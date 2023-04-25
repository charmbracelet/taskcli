package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	home, _ := os.UserHomeDir()
	taskDir := fmt.Sprintf("%s/.tasks", home)
	if err := initTaskDir(taskDir); err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("sqlite3", fmt.Sprintf("%s/tasks.db", taskDir))
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
	/*
		if err := t.insert("get milk", ""); err != nil {
			log.Fatal(err)
		}

		if err := t.insert("get cereal", ""); err != nil {
			log.Fatal(err)
		}
	*/
	// if err := t.delete(1); err != nil {
	// 	log.Fatal(err)
	// }
	tasks, _ := t.getTasks()
	fmt.Printf("%#v", tasks)

	if err := t.editProject(3, "groceries"); err != nil {
		log.Fatal(err)
	}

	tasks, _ = t.getTasks()
	fmt.Printf("%#v", tasks)
}
