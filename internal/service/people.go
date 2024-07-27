package service

import (
	"time_tracker/internal/models"
	"time_tracker/internal/repository"
)

type PeopleService struct {
	repo repository.People
}

func NewPeopleService(repo repository.People) *PeopleService {
	return &PeopleService{repo: repo}
}

func (p *PeopleService) GetByPassport(passportSerie int, passportNumber int) (models.People, error) {
	return p.repo.GetByPassport(passportSerie, passportNumber)
}
