package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	gap "github.com/muesli/go-app-paths"
)

/*
Here's the plan, we're going to store our data in a dedicated data directory at
`XDG_DATA_HOME/.tasks`. That's where we will store a copy of our SQLite DB.
*/

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

// openDB opens a SQLite database and stores that database in our special spot.
func openDB(path string) (*taskDB, error) {
	db, err := sql.Open("sqlite3", filepath.Join(path, "tasks.db"))
	if err != nil {
		return nil, err
	}
	t := taskDB{db, path}
	if !t.tableExists("tasks") {
		err := t.createTable()
		if err != nil {
			return nil, err
		}
	}
	return &t, nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
