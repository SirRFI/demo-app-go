package handlers

import (
	"demo-app-go/storage"
	"demo-app-go/task"
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type TaskHandler struct {
	repository *storage.TaskRepository
}

func NewTaskHandler(repository *storage.TaskRepository) *TaskHandler {
	return &TaskHandler{repository: repository}
}

type taskResponse struct {
	Id          task.ID    `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   *time.Time `json:"updatedAt"`
}

func (h *TaskHandler) List(c echo.Context) error {
	entities, err := h.repository.List()
	if err != nil {
		return err
	}

	result := make([]taskResponse, len(entities))
	for i, entity := range entities {
		result[i] = createTaskResponse(entity)
	}

	return c.JSON(http.StatusOK, result)
}

func (h *TaskHandler) Get(c echo.Context) error {
	id, err := getTaskId(c)
	if err != nil {
		return err
	}

	entity, err := h.repository.GetByID(id)
	if errors.Is(err, storage.ErrResourceNotFound) {
		return c.NoContent(http.StatusNotFound)
	}
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, createTaskResponse(entity))
}

type taskRequest struct {
	Title       string `json:"title" validate:"required,min=1"`
	Description string `json:"description"`
}

func (h *TaskHandler) Add(c echo.Context) error {
	data := &taskRequest{}
	err := c.Bind(data)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	data.Title = strings.TrimSpace(data.Title)
	data.Description = strings.TrimSpace(data.Description)

	err = c.Validate(data)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	command := task.NewAddTaskCommand(data.Title, data.Description)
	entity, err := h.repository.Add(command)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, createTaskResponse(entity))
}

func (h *TaskHandler) Update(c echo.Context) error {
	id, err := getTaskId(c)
	if err != nil {
		return err
	}

	data := &taskRequest{}
	err = c.Bind(data)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	data.Title = strings.TrimSpace(data.Title)
	data.Description = strings.TrimSpace(data.Description)

	err = c.Validate(data)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	entity, err := h.repository.GetByID(id)
	if errors.Is(err, storage.ErrResourceNotFound) {
		return c.NoContent(http.StatusNotFound)
	}
	if err != nil {
		return err
	}
	entity.Update(data.Title, data.Description)
	err = h.repository.Save(entity)
	if errors.Is(err, storage.ErrResourceNotFound) {
		return c.NoContent(http.StatusNotFound)
	}
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, createTaskResponse(entity))
}

func (h *TaskHandler) Delete(c echo.Context) error {
	id, err := getTaskId(c)
	if err != nil {
		return err
	}

	err = h.repository.Delete(id)
	if errors.Is(err, storage.ErrResourceNotFound) {
		return c.NoContent(http.StatusNotFound)
	}
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func getTaskId(c echo.Context) (task.ID, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return 0, err
	}
	return task.ID(id), nil
}

func createTaskResponse(entity task.Task) taskResponse {
	return taskResponse{
		Id:          entity.Id(),
		Title:       entity.Title(),
		Description: entity.Description(),
		CreatedAt:   entity.CreatedAt(),
		UpdatedAt:   entity.UpdatedAt(),
	}
}
