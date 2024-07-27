package service

import (
	"time_tracker/internal/models"
	"time_tracker/internal/repository"
)

type UserService struct {
	repo repository.User
}

func NewUserService(repo repository.User) *UserService {
	return &UserService{repo: repo}
}

func (us *UserService) Get(request models.UsersGetRequest) ([]models.User, error) {
	return us.repo.Get(request)
}

func (us *UserService) GetByID(id string) (models.User, bool, error) {
	return us.repo.GetByID(id)
}

func (us *UserService) Update(id string, request models.UserUpdateRequest) (bool, error) {
	return us.repo.Update(id, request)
}

func (us *UserService) Delete(id string) (bool, error) {
	return us.repo.Delete(id, false)
}

func (us *UserService) Create(request models.User) (uint, error) {
	return us.repo.Create(request)
}
