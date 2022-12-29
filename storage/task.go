package storage

import (
	"database/sql"
	"demo-app-go/task"
	"errors"
	"github.com/jmoiron/sqlx"
	"time"
)

var ErrResourceNotFound = errors.New("resource not found")

type TaskRepository struct {
	db *sqlx.DB
}

func NewTaskRepository(db *sqlx.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

type taskRecord struct {
	Id          task.ID    `db:"id"`
	Title       string     `db:"title"`
	Description string     `db:"description"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`
}

func (r *TaskRepository) List() ([]task.Task, error) {
	var records []taskRecord
	err := r.db.Select(&records, "SELECT * FROM task;")
	if errors.Is(err, sql.ErrNoRows) {
		return []task.Task{}, nil
	}
	if err != nil {
		return nil, err
	}

	result := make([]task.Task, len(records))
	for i, record := range records {
		result[i] = createTask(record)
	}

	return result, nil
}

func (r *TaskRepository) GetByID(id task.ID) (task.Task, error) {
	var record taskRecord
	err := r.db.Get(&record, "SELECT * FROM task WHERE id=?;", id)
	if errors.Is(err, sql.ErrNoRows) {
		return task.Task{}, ErrResourceNotFound
	}
	if err != nil {
		return task.Task{}, err
	}

	return createTask(record), nil
}

func (r *TaskRepository) Add(addTask task.AddTaskCommand) (task.Task, error) {
	rows, err := r.db.NamedQuery(
		"INSERT INTO task (title, description, created_at) VALUES (:title, :description, :createdAt) RETURNING *;",
		map[string]any{
			"title":       addTask.Title(),
			"description": addTask.Description(),
			"createdAt":   addTask.CreatedAt(),
		},
	)
	if err != nil {
		return task.Task{}, err
	}

	var record taskRecord
	if rows.Next() == false {
		return task.Task{}, errors.New("sql: Next() failed")
	}
	err = rows.StructScan(&record)
	if err != nil {
		return task.Task{}, err
	}

	return createTask(record), nil
}

func (r *TaskRepository) Save(task task.Task) error {
	_, err := r.db.NamedExec(
		"UPDATE task SET title=:title, description=:description, updated_at=:updatedAt WHERE id=:id;",
		map[string]any{
			"id":          task.Id(),
			"title":       task.Title(),
			"description": task.Description(),
			"updatedAt":   task.UpdatedAt(),
		},
	)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrResourceNotFound
	}
	if err != nil {
		return err
	}

	return nil
}

func (r *TaskRepository) Delete(id task.ID) error {
	_, err := r.db.Exec("DELETE FROM task WHERE id=?;", id)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrResourceNotFound
	}
	if err != nil {
		return err
	}

	return nil
}

func createTask(record taskRecord) task.Task {
	return task.NewTask(record.Id, record.Title, record.Description, record.CreatedAt, record.UpdatedAt)
}
