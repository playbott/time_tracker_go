package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"time_tracker/internal/models"
)

type PeopleExternal struct {
	httpClient *fiber.Client
	url        string
}

func NewPeopleExternal(httpClient *fiber.Client, url string) *PeopleExternal {
	return &PeopleExternal{httpClient: httpClient, url: url}
}

func (r *PeopleExternal) GetByPassport(passportSerie int, passportNumber int) (models.People, error) {
	defer fiber.ReleaseClient(r.httpClient)
	var people models.People
	url := fmt.Sprintf("%s/info?passportSerie=%d&passportNumber=%d", r.url, passportSerie, passportNumber)
	agent := r.httpClient.Get(url)
	body := make([]byte, 0)
	code, body, errs := agent.Bytes()

	if code != fiber.StatusOK {
		return people, errors.New(fmt.Sprintf("people info fetch failed. status code: %v", code))
	}

	for _, e := range errs {
		if e != nil {
			logrus.Errorf("People info fetch failed: %v", e)
		}
	}

	if len(body) == 0 {
		logrus.Errorf("People info fetch failed: body is empty")
		return people, errors.New("people info fetch failed: body is empty")
	}

	err := json.Unmarshal(body, &people)
	if err != nil {
		logrus.Errorf("People info fetch failed: %v", err)
		return people, err
	}

	return people, nil
}
