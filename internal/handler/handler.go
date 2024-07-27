package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"time_tracker/internal/middleware"
	"time_tracker/internal/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) Routes(app *fiber.App) {
	var _validator = validator.New()
	app.Use(middleware.ValidatorMiddleware(_validator))

	app.Get("/swagger/*", swagger.HandlerDefault)

	api := app.Group("/api")
	apiV1 := api.Group("/v1")

	user := apiV1.Group("/user")
	user.Get("/", h.GetUsers)
	user.Get("/:id", h.GetUser)
	user.Post("/", h.CreateUser)
	user.Patch("/:id", h.UpdateUser)
	user.Delete("/:id", h.DeleteUser)

	task := apiV1.Group("/task")
	task.Get("/", h.GetTasks)
	task.Get("/user", h.GetUserTasks)
	task.Post("/create-start", h.CreateAndStartTask)
	task.Post("/complete", h.CompleteTask)
}
