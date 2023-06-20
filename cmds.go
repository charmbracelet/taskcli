package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tasks",
	Short: "A CLI task management tool for ~slaying~ your to do list.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var addCmd = &cobra.Command{
	Use:   "add NAME",
	Short: "Add a new task with an optional project name",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := openDB(setupPath())
		if err != nil {
			return err
		}
		defer t.db.Close()
		project, err := cmd.Flags().GetString("project")
		if err != nil {
			return err
		}
		if err := t.insert(args[0], project); err != nil {
			return err
		}
		return nil
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a task by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := openDB(setupPath())
		if err != nil {
			return err
		}
		defer t.db.Close()
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		return t.delete(uint(id))
	},
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a task by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := openDB(setupPath())
		if err != nil {
			return err
		}
		defer t.db.Close()
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			return err
		}
		project, err := cmd.Flags().GetString("project")
		if err != nil {
			return err
		}
		status, err := cmd.Flags().GetString("status")
		if err != nil {
			return err
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		newTask := task{uint(id), name, project, status, time.Time{}}
		return t.update(newTask)
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all your tasks",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := openDB(setupPath())
		if err != nil {
			return err
		}
		defer t.db.Close()
		tasks, err := t.getTasks()
		if err != nil {
			return err
		}
		table := setupTable(tasks)
		fmt.Print(table.View())
		return nil
	},
}

func setupTable(tasks []task) table.Model {
	columns := []table.Column{
		{Title: "ID", Width: 3},
		{Title: "Name", Width: 30},
		{Title: "Project", Width: 20},
		{Title: "Status", Width: 10},
		{Title: "Created At", Width: 10},
	}
	var rows []table.Row
	for _, task := range tasks {
		rows = append(rows, table.Row{fmt.Sprintf("%d", task.ID), task.Name, task.Project, task.Status, task.Created.Format("2006-01-02")})
	}
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(false),
		table.WithHeight(len(tasks)),
	)
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	t.SetStyles(s)
	return t
}

func init() {
	addCmd.Flags().StringP(
		"project",
		"p",
		"",
		"specify a project for your task",
	)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(listCmd)
	updateCmd.Flags().StringP(
		"name",
		"n",
		"",
		"specify a name for your task",
	)
	updateCmd.Flags().StringP(
		"project",
		"p",
		"",
		"specify a project for your task",
	)
	updateCmd.Flags().StringP(
		"status",
		"s",
		"",
		"specify a status for your task",
	)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(deleteCmd)
}
