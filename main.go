package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	gap "github.com/muesli/go-app-paths"
)

// setupPath uses XDG to create the necessary data dirs for the program.
func setupPath() string {
	// get XDG paths
	scope := gap.NewScope(gap.User, "tasks")
	dirs, err := scope.DataDirs()
	if err != nil {
		log.Fatal(err)
	}
	// create the app base dir, if it doesn't exist
	var taskDir string
	if len(dirs) > 0 {
		taskDir = dirs[0]
	} else {
		taskDir, _ = os.UserHomeDir()
	}
	if err := initTaskDir(taskDir); err != nil {
		log.Fatal(err)
	}
	return taskDir
}

func main() {
	path := setupPath()
	db, err := sql.Open("sqlite3", fmt.Sprintf("%s/tasks.db", path))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	t := taskDB{db}

	if !t.tableExists("tasks") {
		err := t.createTable()
		if err != nil {
			log.Fatal(err)
		}
	}
	if err := t.insert("cook currywurst", ""); err != nil {
		log.Fatal(err)
	}
	tasks, _ := t.getTasks()
	fmt.Printf("%#v", tasks)
}
