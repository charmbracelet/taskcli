package main

import (
	"database/sql"
	"fmt"
	"os"
	"reflect"
	"time"
)

// TODO: make a single Update/Set function

type status int

const (
	todo status = iota
	completed
)

// SQL helpers
const (
	updateCmd = "UPDATE tasks SET name = ?, project = ?, status = ? WHERE id = ?"
	insertCmd = "INSERT INTO tasks(name, project, status, created) VALUES( ?, ?, ?, ?)"
	deleteCmd = "DELETE FROM tasks WHERE id = ?"

	commitTxnErr  = "unable to commit txn: %w"
	prepareTxnErr = "unable to prepare txn: %w"
	editTxnErr    = "unable to edit: %w"
	beginDBErr    = "unable to begin db: %w"
)

func (s status) String() string {
	return [...]string{"todo", "done"}[s]
}

// note for reflect: only exported fields of a struct are settable.
type task struct {
	ID      uint
	Name    string
	Project string
	Status  string
	Created time.Time
}

type taskDB struct {
	db *sql.DB
}

func initTaskDir(path string) error {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return os.Mkdir(path, 0o755)
		}

		return err
	}
	return nil
}

func (t *taskDB) tableExists(name string) bool {
	if _, err := t.db.Query("SELECT * FROM tasks"); err == nil {
		return true
	}
	return false
}

func (t *taskDB) createTable() error {
	cmd := `CREATE TABLE "tasks" ( "id" INTEGER, "name" TEXT NOT NULL, "project" TEXT, "status" TEXT, "created" DATETIME, PRIMARY KEY("id" AUTOINCREMENT))`
	_, err := t.db.Exec(cmd)
	return err
}

func (t *taskDB) insert(name, project string) error {
	tx, err := t.db.Begin()
	if err != nil {
		return fmt.Errorf(beginDBErr, err)
	}
	stmt, err := tx.Prepare(insertCmd)
	if err != nil {
		return fmt.Errorf(prepareTxnErr, err)
	}
	defer stmt.Close()
	if _, err := stmt.Exec(name, project, todo.String(), time.Now().Format(time.RFC822)); err != nil {
		return fmt.Errorf("unable to insert: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf(commitTxnErr, err)
	}
	return err
}

func (t *taskDB) delete(id uint) error {
	tx, err := t.db.Begin()
	if err != nil {
		return fmt.Errorf(beginDBErr, err)
	}
	stmt, err := tx.Prepare(deleteCmd)
	if err != nil {
		return fmt.Errorf(prepareTxnErr, err)
	}
	defer stmt.Close()
	if _, err := stmt.Exec(id); err != nil {
		return fmt.Errorf("unable to delete: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf(commitTxnErr, err)
	}
	return err
}

// Update the task in the db. Provide new values for the fields you want to
// change, keep them empty if unchanged.
func (t *taskDB) update(task task) error {
	tx, err := t.db.Begin()
	if err != nil {
		return fmt.Errorf(beginDBErr, err)
	}

	orig, err := t.getTask(task.ID)
	orig.merge(task)

	stmt, err := tx.Prepare(updateCmd)
	if err != nil {
		return fmt.Errorf(prepareTxnErr, err)
	}
	defer stmt.Close()
	if _, err := stmt.Exec(orig.Name, orig.Project, orig.Status, orig.ID); err != nil {
		return fmt.Errorf(editTxnErr, err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf(commitTxnErr, err)
	}
	return err
}

func (orig *task) merge(t task) {
	uValues := reflect.ValueOf(&t).Elem()
	oValues := reflect.ValueOf(orig).Elem()
	for i := 0; i < uValues.NumField(); i++ {
		uField := uValues.Field(i).Interface()
		if oValues.CanSet() {
			if v, ok := uField.(int64); ok && uField != 0 {
				oValues.Field(i).SetInt(v)
			}
			if v, ok := uField.(string); ok && uField != "" {
				oValues.Field(i).SetString(v)
			}
		}
	}
}

func (t *taskDB) getTasks() ([]task, error) {
	var tasks []task
	tx, err := t.db.Begin()
	if err != nil {
		return tasks, fmt.Errorf(beginDBErr, err)
	}
	stmt, err := tx.Prepare("SELECT * FROM tasks")
	if err != nil {
		return tasks, fmt.Errorf(prepareTxnErr, err)
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		return tasks, fmt.Errorf("unable to get values: %w", err)
	}
	for rows.Next() {
		var task task
		err = rows.Scan(&task.ID, &task.Name, &task.Project, &task.Status, &task.Created)
		if err != nil {
			return tasks, err
		}

		tasks = append(tasks, task)
	}
	if err := tx.Commit(); err != nil {
		return tasks, fmt.Errorf(commitTxnErr, err)
	}
	return tasks, err
}

func (t *taskDB) getTask(id uint) (task, error) {
	var task task
	tx, err := t.db.Begin()
	if err != nil {
		return task, fmt.Errorf(beginDBErr, err)
	}
	stmt, err := tx.Prepare("SELECT * FROM tasks WHERE id = ?")
	if err != nil {
		return task, fmt.Errorf(prepareTxnErr, err)
	}
	defer stmt.Close()
	err = tx.QueryRow(
		`SELECT * FROM tasks WHERE id = ?`,
		id).Scan(
		&task.ID,
		&task.Name,
		&task.Project,
		&task.Status,
		&task.Created)
	if err != nil {
		return task, err
	}
	if err := tx.Commit(); err != nil {
		return task, fmt.Errorf(commitTxnErr, err)
	}
	return task, nil
}
