package handler

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"time_tracker/internal/models"
)

// GetTasks task search with filter and pagination.
//
//	@Description	Task search with filter and pagination.
//	@Summary		task search with filter and pagination
//	@Tags			Task
//	@Accept			json
//	@Produce		json
//	@Param			params	body		models.TasksGetRequest	true	"search parameters"
//	@Success		200		{object}	models.TasksSearchResponse
//	@Failure		400		{object}	models.ResponseError
//	@Failure		500		{object}	models.ResponseError
//	@Router			/api/v1/task [get]
func (h *Handler) GetTasks(c *fiber.Ctx) error {
	var err error
	v := c.Locals("validator").(*validator.Validate)
	var taskGetRequest models.TasksGetRequest
	if err = c.BodyParser(&taskGetRequest); err != nil {
		logrus.Error("Failed to parse request body: ", err)
		return c.Status(fiber.StatusBadRequest).JSON(models.ResponseError{Message: "Cannot parse JSON input data"})
	}

	err = v.Struct(taskGetRequest)
	if err != nil {
		logrus.Error(err)
		return c.Status(fiber.StatusBadRequest).JSON(models.ResponseError{Message: err.Error()})
	}

	fmt.Printf("%+v\n", taskGetRequest)

	tasks, err := h.services.Task.Get(taskGetRequest)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ResponseError{Message: "Unknown error"})
	}

	return c.JSON(models.TasksSearchResponse{Response: models.Response[[]models.Task2]{Data: tasks}, Page: taskGetRequest.Page})
}

// GetUserTasks func gets a task by user id.
//
//	@Description	Get task by user id.
//	@Summary		get task by user id
//	@Tags			Task
//	@Accept			json
//	@Produce		json
//	@Param			id				query		number	true	"user id"
//	@Param			completed		query		boolean	true	"completed tasks"
//	@Param			durationDesc	query		boolean	true	"duration descending"
//	@Success		200				{object}	models.Response[[]models.Task2]
//	@Failure		404				{object}	string
//	@Router			/api/v1/task/user [get]
func (h *Handler) GetUserTasks(c *fiber.Ctx) error {
	v := c.Locals("validator").(*validator.Validate)
	q := c.Queries()
	userId := q["id"]
	completed := q["completed"]
	durationDesc := q["durationDesc"]
	err := v.Var(userId, "required,numeric,min=1")
	err2 := v.Var(completed, "required,boolean")
	err3 := v.Var(durationDesc, "required,boolean")
	if err != nil || err2 != nil || err3 != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ResponseError{Message: "Invalid parameters"})
	}

	tasks, notfound, err := h.services.Task.GetByID(userId, completed == "true", durationDesc == "true")
	if err != nil {
		if notfound {
			return c.Status(fiber.StatusNotFound).JSON(models.ResponseError{Message: "Task not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ResponseError{Message: "Unknown error"})
	}

	return c.JSON(models.Response[[]models.Task2]{Data: tasks})
}

// CreateAndStartTask creates and starts a new task.
//
//	@Description	Create and start a new task.
//	@Summary		create and start a new task
//	@Tags			Task
//	@Accept			json
//	@Produce		json
//	@Param			task	body		models.TaskCreateRequest	true	"user id and task title"
//	@Success		201		{object}	models.Response[models.TaskCreateRequest]
//	@Router			/api/v1/task/create-start [post]
func (h *Handler) CreateAndStartTask(c *fiber.Ctx) error {
	v := c.Locals("validator").(*validator.Validate)
	var err error
	var newTask models.TaskCreateRequest

	if err = c.BodyParser(&newTask); err != nil {
		logrus.Error("Failed to parse request body: ", err)
		return c.Status(fiber.StatusBadRequest).JSON(models.ResponseError{Message: "Cannot parse JSON input data"})
	}

	err = v.Struct(newTask)
	if err != nil {
		logrus.Error(err)
		return c.Status(fiber.StatusBadRequest).JSON(models.ResponseError{Message: err.Error()})
	}

	taskId, err := h.services.Task.CreateAndStart(newTask.UserId, newTask.Title)
	if err != nil {
		logrus.Info("Failed to create task: ", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(models.ResponseError{Message: "Failed to create user"})
	}
	newTask.Id = taskId
	logrus.Info("Task created and started. ID: ", newTask.Id)
	return c.Status(fiber.StatusCreated).JSON(models.Response[models.TaskCreateRequest]{Data: newTask})
}

// CompleteTask completes a task.
//
//	@Description	Complete a task.
//	@Summary		complete a task
//	@Tags			Task
//	@Accept			json
//	@Produce		json
//	@Param			task	body		models.TaskCompleteRequest	true	"task id and user id"
//	@Success		200		{object}	models.Response[string]
//	@Router			/api/v1/task/complete [post]
func (h *Handler) CompleteTask(c *fiber.Ctx) error {
	v := c.Locals("validator").(*validator.Validate)
	var err error
	var task models.TaskCompleteRequest

	if err = c.BodyParser(&task); err != nil {
		logrus.Error("Failed to parse request body: ", err)
		return c.Status(fiber.StatusBadRequest).JSON(models.ResponseError{Message: "Cannot parse JSON input data"})
	}

	err = v.Struct(task)
	if err != nil {
		logrus.Error(err)
		return c.Status(fiber.StatusBadRequest).JSON(models.ResponseError{Message: err.Error()})
	}

	err = h.services.Task.Complete(task.UserId, task.Id)
	if err != nil {
		logrus.Info("Failed to create task: ", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(models.ResponseError{Message: "Failed to create user"})
	}

	logrus.Info("Task completed. ID: ", task.Id)

	return c.Status(fiber.StatusOK).JSON(models.Response[string]{Message: "Task completed"})
}
