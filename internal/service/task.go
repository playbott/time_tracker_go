package service

import (
	"fmt"
	"strings"
	"time_tracker/internal/models"
	"time_tracker/internal/repository"
	"time_tracker/pkg"
)

type TaskService struct {
	repo repository.Task
}

func NewTaskService(repo repository.Task) *TaskService {
	return &TaskService{repo: repo}
}

func (us *TaskService) CreateAndStart(userId uint, taskTitle string) (uint, error) {
	return us.repo.CreateAndStart(userId, taskTitle)
}

func (us *TaskService) Complete(userId uint, taskId uint) error {
	return us.repo.Complete(userId, taskId)
}

func (us *TaskService) Get(request models.TasksGetRequest) ([]models.Task2, error) {
	return us.repo.Get(request)
}

func (us *TaskService) GetByID(id string, completed bool, durationDesc bool) ([]models.Task2, bool, error) {
	tasks, notFound, err := us.repo.GetByID(id, completed, durationDesc)
	if err != nil {
		return tasks, notFound, err
	}

	hms := func(v int, m string) string {
		return fmt.Sprintf("%02d %s", v, m)
	}

	newTasks := make([]models.Task2, 0)

	for _, task := range tasks {
		if task.CompletedAt != "" {
			duration := pkg.SecondsToTime(int64(task.DurationSeconds))
			task.DurationString = strings.Join([]string{
				hms(duration.Hour(), "час"),
				hms(duration.Minute(), "мин"),
				hms(duration.Second(), "сек"),
			}, ", ")
		}
		newTasks = append(newTasks, task)
	}

	return newTasks, notFound, nil
}
