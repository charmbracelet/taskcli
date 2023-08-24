package main

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestDelete(t *testing.T) {
	tests := []struct {
		want task
	}{
		{
			want: task{
				ID:      1,
				Name:    "get milk",
				Project: "groceries",
				Status:  "todo",
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.want.Name, func(t *testing.T) {
			tDB := setup()
			defer teardown(tDB)
			if err := tDB.insert(tc.want.Name, tc.want.Project); err != nil {
				t.Fatalf("unable to insert tasks: %v", err)
			}
			tasks, err := tDB.getTasks()
			if err != nil {
				t.Fatalf("unable to get tasks: %v", err)
			}
			tc.want.Created = tasks[0].Created
			if !reflect.DeepEqual(tc.want, tasks[0]) {
				t.Fatalf("got %v, want %v", tc.want, tasks)
			}
			if err := tDB.delete(1); err != nil {
				t.Fatalf("unable to delete tasks: %v", err)
			}
			tasks, err = tDB.getTasks()
			if err != nil {
				t.Fatalf("unable to get tasks: %v", err)
			}
			if len(tasks) != 0 {
				t.Fatalf("expected tasks to be empty, got: %v", tasks)
			}
		})
	}
}

func TestGetTask(t *testing.T) {
	tests := []struct {
		want task
	}{
		{
			want: task{
				ID:      1,
				Name:    "get milk",
				Project: "groceries",
				Status:  todo.String(),
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.want.Name, func(t *testing.T) {
			tDB := setup()
			defer teardown(tDB)
			if err := tDB.insert(tc.want.Name, tc.want.Project); err != nil {
				t.Fatalf("we ran into an unexpected error: %v", err)
			}
			task, err := tDB.getTask(tc.want.ID)
			if err != nil {
				t.Fatalf("we ran into an unexpected error: %v", err)
			}
			tc.want.Created = task.Created
			if !reflect.DeepEqual(task, tc.want) {
				t.Fatalf("got: %#v, want: %#v", task, tc.want)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	tests := []struct {
		new  *task
		old  *task
		want task
	}{
		{
			new: &task{
				ID:      1,
				Name:    "strawberries",
				Project: "",
				Status:  "",
			},
			old: &task{
				ID:      1,
				Name:    "get milk",
				Project: "groceries",
				Status:  todo.String(),
			},
			want: task{
				ID:      1,
				Name:    "strawberries",
				Project: "groceries",
				Status:  todo.String(),
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.new.Name, func(t *testing.T) {
			tDB := setup()
			defer teardown(tDB)
			if err := tDB.insert(tc.old.Name, tc.old.Project); err != nil {
				t.Fatalf("we ran into an unexpected error: %v", err)
			}
			if err := tDB.update(*tc.new); err != nil {
				t.Fatalf("we ran into an unexpected error: %v", err)
			}
			task, err := tDB.getTask(tc.want.ID)
			if err != nil {
				t.Fatalf("we ran into an unexpected error: %v", err)
			}
			tc.want.Created = task.Created
			if !reflect.DeepEqual(task, tc.want) {
				t.Fatalf("got: %#v, want: %#v", task, tc.want)
			}
		})
	}
}

func TestMerge(t *testing.T) {
	tests := []struct {
		new  task
		old  task
		want task
	}{
		{
			new: task{
				ID:      1,
				Name:    "strawberries",
				Project: "",
				Status:  "",
			},
			old: task{
				ID:      1,
				Name:    "get milk",
				Project: "groceries",
				Status:  todo.String(),
			},
			want: task{
				ID:      1,
				Name:    "strawberries",
				Project: "groceries",
				Status:  todo.String(),
			},
		},
	}
	for _, tc := range tests {
		tc.old.merge(tc.new)
		if !reflect.DeepEqual(tc.old, tc.want) {
			t.Fatalf("got: %#v, want %#v", tc.new, tc.want)
		}
	}
}

func TestGetTasksByStatus(t *testing.T) {
	tests := []struct {
		want task
	}{
		{
			want: task{
				ID:      1,
				Name:    "get milk",
				Project: "groceries",
				Status:  todo.String(),
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.want.Name, func(t *testing.T) {
			tDB := setup()
			defer teardown(tDB)
			if err := tDB.insert(tc.want.Name, tc.want.Project); err != nil {
				t.Fatalf("we ran into an unexpected error: %v", err)
			}
			tasks, err := tDB.getTasksByStatus(tc.want.Status)
			if err != nil {
				t.Fatalf("we ran into an unexpected error: %v", err)
			}
			if len(tasks) < 1 {
				t.Fatalf("expected 1 value, got %#v", tasks)
			}
			tc.want.Created = tasks[0].Created
			if !reflect.DeepEqual(tasks[0], tc.want) {
				t.Fatalf("got: %#v, want: %#v", tasks, tc.want)
			}
		})
	}
}

func setup() *taskDB {
	path := filepath.Join(os.TempDir(), "test.db")
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatal(err)
	}
	t := taskDB{db, path}
	if !t.tableExists("tasks") {
		err := t.createTable()
		if err != nil {
			log.Fatal(err)
		}
	}
	return &t
}

func teardown(tDB *taskDB) {
	tDB.db.Close()
	os.Remove(tDB.dataDir)
}
