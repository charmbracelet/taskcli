package main

import (
	"database/sql"
	"fmt"
	"os"
)

type task struct {
	id      uint
	name    string
	project string
}

type taskDB struct {
	db *sql.DB
}

func initTaskDir(path string) (err error) {
	if _, e := os.Stat(path); e != nil {
		err = os.Mkdir(path, 0o755)
	}
	return err
}

func (t *taskDB) tableExists(name string) bool {
	if _, err := t.db.Query("select * from tasks;"); err == nil {
		return true
	}
	return false
}

func (t *taskDB) createTable() (err error) {
	cmd := `CREATE TABLE "tasks" ( "id" INTEGER, "name" TEXT NOT NULL, "project" TEXT, PRIMARY KEY("id" AUTOINCREMENT));`
	_, err = t.db.Exec(cmd)
	return err
}

func (t *taskDB) insert(name, project string) (err error) {
	tx, err := t.db.Begin()
	if err != nil {
		return fmt.Errorf("unable to begin db: %w", err)
	}

	stmt, err := tx.Prepare("insert into tasks(name, project) VALUES( ?, ?)")
	if err != nil {
		return fmt.Errorf("unable to prepare txn: %w", err)
	}
	defer stmt.Close()

	if _, err := stmt.Exec(name, project); err != nil {
		return fmt.Errorf("unable to insert: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("unable to commit txn: %w", err)
	}
	return err
}

func (t *taskDB) delete(id uint) (err error) {
	tx, err := t.db.Begin()
	if err != nil {
		return fmt.Errorf("unable to begin db: %w", err)
	}

	stmt, err := tx.Prepare("delete from tasks where id = ?")
	if err != nil {
		return fmt.Errorf("unable to prepare txn: %w", err)
	}
	defer stmt.Close()

	if _, err := stmt.Exec(id); err != nil {
		return fmt.Errorf("unable to delete: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("unable to commit txn: %w", err)
	}
	return err
}

func (t *taskDB) edit(id uint, name string) (err error) {
	tx, err := t.db.Begin()
	if err != nil {
		return fmt.Errorf("unable to begin db: %w", err)
	}

	stmt, err := tx.Prepare("update tasks set name = ? where id = ?")
	if err != nil {
		return fmt.Errorf("unable to prepare txn: %w", err)
	}
	defer stmt.Close()

	if _, err := stmt.Exec(name, id); err != nil {
		return fmt.Errorf("unable to edit: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("unable to commit txn: %w", err)
	}
	return err
}

func (t *taskDB) editProject(id uint, project string) (err error) {
	tx, err := t.db.Begin()
	if err != nil {
		return fmt.Errorf("unable to begin db: %w", err)
	}

	stmt, err := tx.Prepare("update tasks set project = ? where id = ?")
	if err != nil {
		return fmt.Errorf("unable to prepare txn: %w", err)
	}
	defer stmt.Close()

	if _, err := stmt.Exec(project, id); err != nil {
		return fmt.Errorf("unable to edit: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("unable to commit txn: %w", err)
	}
	return err
}

func (t *taskDB) getTasks() (tasks []task, err error) {
	tx, err := t.db.Begin()
	if err != nil {
		return tasks, fmt.Errorf("unable to begin db: %w", err)
	}

	stmt, err := tx.Prepare("select * from tasks")
	if err != nil {
		return tasks, fmt.Errorf("unable to prepare txn: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return tasks, fmt.Errorf("unable to get values: %w", err)
	}

	for rows.Next() {
		var task task
		err = rows.Scan(&task.id, &task.name, &task.project)
		if err != nil {
			return tasks, err
		}

		tasks = append(tasks, task)
	}

	if err := tx.Commit(); err != nil {
		return tasks, fmt.Errorf("unable to commit txn: %w", err)
	}
	return tasks, err
}
