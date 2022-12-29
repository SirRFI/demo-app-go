package task

import "time"

type ID uint

type Task struct {
	id          ID
	title       string
	description string
	createdAt   time.Time
	updatedAt   *time.Time
}

func NewTask(id ID, title string, description string, createdAt time.Time, updatedAt *time.Time) Task {
	return Task{id: id, title: title, description: description, createdAt: createdAt, updatedAt: updatedAt}
}

func (t *Task) Update(title string, description string) {
	t.title = title
	t.description = description
	now := time.Now()
	t.updatedAt = &now
}

func (t Task) Id() ID {
	return t.id
}

func (t Task) Title() string {
	return t.title
}

func (t Task) Description() string {
	return t.description
}

func (t Task) CreatedAt() time.Time {
	return t.createdAt
}

func (t Task) UpdatedAt() *time.Time {
	return t.updatedAt
}

type AddTaskCommand struct {
	title       string
	description string
	createdAt   time.Time
}

func NewAddTaskCommand(title string, description string) AddTaskCommand {
	return AddTaskCommand{title: title, description: description, createdAt: time.Now()}
}

func (a AddTaskCommand) Title() string {
	return a.title
}

func (a AddTaskCommand) Description() string {
	return a.description
}

func (a AddTaskCommand) CreatedAt() time.Time {
	return a.createdAt
}
