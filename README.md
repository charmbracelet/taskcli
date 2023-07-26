# Creating a Task Manager Tutorial

This project is inspired by the incredible work on Task Warrior, an open source
CLI task manager. I use this project quite a bit for managing my projects
without leaving the safety and comfort of my terminal. (⌐■_■)

We built a kanban board TUI in a previous [tutorial][kanban-video], so the
idea here is that we're going to build a task management CLI with [Cobra][cobra] that has Lip Gloss
styles *and* can be viewed using our kanban board.

Here's the plan:

## Checklist

If you're following along with our tutorials for this project, or even if you
want to try and tackle it yourself first, then look at our solutions, here's
what you need to do:

### Data Storage
- [ ] set up a SQLite database
    - [ ] open SQLite DB
    - [ ] add task
    - [ ] delete task
    - [ ] edit task
    - [ ] get tasks

### Making a CLI with Cobra
- [ ] add CLI
    - [ ] add task
    - [ ] delete task
    - [ ] edit task
    - [ ] get tasks

### Add a little... *Je ne sais quoi*
- [ ] print to table layout with [Lip Gloss][lipgloss]
- [ ] print to Kanban layout with [Lip Gloss][lipgloss]

## Project Layout

`db.go` - here we create our custom `task` struct and our data layer.

`main.go` - our main file handles our initial setup including opening a
database and setting the data path for our application.

`cmds.go` - this is where we do all of our Cobra commands and setup for our
CLI.

[lipgloss]: https://github.com/charmbracelet/lipgloss
[charm]: https://github.com/charmbracelet/charm
[cobra]: https://github.com/spf13/cobra
[kanban-video]: https://www.youtube.com/watch?v=ZA93qgdLUzM&list=PLLLtqOZfy0pcFoSIeGXO-SOaP9qLqd_H6
