package repository

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"time_tracker/internal/models"
)

type User interface {
	Get(request models.UsersGetRequest) ([]models.User, error)
	GetByID(id string) (user models.User, notFound bool, err error)
	Create(user models.User) (id uint, err error)
	Update(id string, user models.UserUpdateRequest) (notFound bool, err error)
	Delete(id string, deleteRecord bool) (notFound bool, err error)
}

type Task interface {
	CreateAndStart(userId uint, taskTitle string) (id uint, err error)
	Complete(userId uint, taskId uint) error
	Get(request models.TasksGetRequest) ([]models.Task2, error)
	GetByID(id string, completed bool, durationDesc bool) ([]models.Task2, bool, error)
	//Update(id string, user models.UserUpdateRequest) (notFound bool, err error)
	//Delete(id string, deleteRecord bool) (notFound bool, err error)
}

type People interface {
	GetByPassport(passportSerie int, passportNumber int) (models.People, error)
}

type Repository struct {
	User
	Task
	People
}

func NewRepository(db *gorm.DB, httpClient *fiber.Client, externalApiUrl string) *Repository {
	return &Repository{
		User:   NewUserPostgres(db),
		Task:   NewTaskPostgres(db),
		People: NewPeopleExternal(httpClient, externalApiUrl),
	}
}
