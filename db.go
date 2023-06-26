package main

import (
	"database/sql"
	"fmt"
	"os"
	"reflect"
	"time"
)

type status int

const (
	todo status = iota
	inProgress
	done
)

func (s status) String() string {
	return [...]string{"todo", "in progress", "done"}[s]
}

/*
A note on SQL statements:
Make sure you're using parameterized SQL statements to avoid
SQL injections. This format creates prepared statements at run time.
learn more: https://go.dev/doc/database/sql-injection
*/

// note for reflect: only exported fields of a struct are settable.
type task struct {
	ID      uint
	Name    string
	Project string
	Status  string
	Created time.Time
}

type taskDB struct {
	db      *sql.DB
	dataDir string
}

func initTaskDir(path string) error {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return os.Mkdir(path, 0o770)
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
	_, err := t.db.Exec(`CREATE TABLE "tasks" ( "id" INTEGER, "name" TEXT NOT NULL, "project" TEXT, "status" TEXT, "created" DATETIME, PRIMARY KEY("id" AUTOINCREMENT))`)
	return err
}

func (t *taskDB) insert(name, project string) error {
	// We don't care about the returned values, so we're using Exec. If we
	// wanted to reuse these statements, it would be more efficient to use
	// prepared statements. Learn more:
	// https://go.dev/doc/database/prepared-statements
	_, err := t.db.Exec(
		"INSERT INTO tasks(name, project, status, created) VALUES( ?, ?, ?, ?)",
		name,
		project,
		todo.String(),
		time.Now())
	return err
}

func (t *taskDB) delete(id uint) error {
	_, err := t.db.Exec("DELETE FROM tasks WHERE id = ?", id)
	return err
}

// Update the task in the db. Provide new values for the fields you want to
// change, keep them empty if unchanged.
func (t *taskDB) update(task task) error {
	// Get the existing state of the task we want to update.
	orig, err := t.getTask(task.ID)
	if err != nil {
		return err
	}
	orig.merge(task)
	_, err = t.db.Exec(
		"UPDATE tasks SET name = ?, project = ?, status = ? WHERE id = ?",
		orig.Name,
		orig.Project,
		orig.Status,
		orig.ID)
	return err
}

// merge the changed fields to the original task
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
	rows, err := t.db.Query("SELECT * FROM tasks")
	if err != nil {
		return tasks, fmt.Errorf("unable to get values: %w", err)
	}
	for rows.Next() {
		var task task
		err = rows.Scan(
			&task.ID,
			&task.Name,
			&task.Project,
			&task.Status,
			&task.Created,
		)
		if err != nil {
			return tasks, err
		}
		tasks = append(tasks, task)
	}
	return tasks, err
}

func (t *taskDB) getTask(id uint) (task, error) {
	var task task
	err := t.db.QueryRow("SELECT * FROM tasks WHERE id = ?", id).
		Scan(
			&task.ID,
			&task.Name,
			&task.Project,
			&task.Status,
			&task.Created,
		)
	return task, err
}
