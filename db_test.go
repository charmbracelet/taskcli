package main

import (
	"database/sql"
	"log"
	"os"
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
		tDB := setup()
		defer tDB.db.Close()

		if err := tDB.insert("get milk", "groceries"); err != nil {
			teardown()
			t.Fatalf("unable to insert tasks: %v", err)
		}

		tasks, err := tDB.getTasks()
		if err != nil {
			teardown()
			t.Fatalf("unable to get tasks: %v", err)
		}

		if !reflect.DeepEqual(tc.want, tasks[0]) {
			teardown()
			t.Fatalf("got %v, want %v", tc.want, tasks)
		}

		if err := tDB.delete(1); err != nil {
			teardown()
			t.Fatalf("unable to delete tasks: %v", err)
		}

		tasks, err = tDB.getTasks()
		if err != nil {
			teardown()
			t.Fatalf("unable to get tasks: %v", err)
		}

		if len(tasks) != 0 {
			teardown()
			t.Fatalf("expected tasks to be empty, got: %v", tasks)
		}
		teardown()
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
		tDB := setup()
		defer tDB.db.Close()
		if err := tDB.insert(tc.want.Name, tc.want.Project); err != nil {
			teardown()
			t.Fatalf("we ran into an unexpected error: %v", err)
		}
		task, err := tDB.getTask(tc.want.ID)
		if err != nil {
			teardown()
			t.Fatalf("we ran into an unexpected error: %v", err)
		}
		if !reflect.DeepEqual(task, tc.want) {
			teardown()
			t.Fatalf("got: %#v, want: %#v", task, tc.want)
		}
	}
	teardown()
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
		tDB := setup()
		defer tDB.db.Close()
		if err := tDB.insert(tc.old.Name, tc.old.Project); err != nil {
			teardown()
			t.Fatalf("we ran into an unexpected error: %v", err)
		}
		if err := tDB.update(*tc.new); err != nil {
			teardown()
			t.Fatalf("we ran into an unexpected error: %v", err)
		}
		task, err := tDB.getTask(tc.want.ID)
		if err != nil {
			teardown()
			t.Fatalf("we ran into an unexpected error: %v", err)
		}
		if !reflect.DeepEqual(task, tc.want) {
			teardown()
			t.Fatalf("got: %#v, want: %#v", task, tc.want)
		}
	}
	teardown()
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

func setup() *taskDB {
	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatal(err)
	}
	t := taskDB{db}

	if !t.tableExists("tasks") {
		err := t.createTable()
		if err != nil {
			teardown()
			log.Fatal(err)
		}
	}
	return &t
}

func teardown() {
	os.Remove("test.db")
}
