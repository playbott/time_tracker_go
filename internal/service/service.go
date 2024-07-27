package service

import (
	"time_tracker/internal/models"
	"time_tracker/internal/repository"
)

type User interface {
	Get(request models.UsersGetRequest) ([]models.User, error)
	GetByID(id string) (models.User, bool, error)
	Create(request models.User) (uint, error)
	Update(id string, user models.UserUpdateRequest) (bool, error)
	Delete(id string) (bool, error)
}

type Task interface {
	CreateAndStart(userId uint, taskTitle string) (id uint, err error)
	Complete(userId uint, taskId uint) error
	Get(request models.TasksGetRequest) ([]models.Task2, error)
	GetByID(id string, completed bool, durationDesc bool) ([]models.Task2, bool, error)
	//Update(id string, user models.UserUpdateRequest) (bool, error)
	//Delete(id string) (bool, error)
}

type People interface {
	GetByPassport(passportSerie int, passportNumber int) (models.People, error)
}

type Service struct {
	User
	Task
	People
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		User:   NewUserService(repos.User),
		Task:   NewTaskService(repos.Task),
		People: NewPeopleService(repos.People),
	}
}
