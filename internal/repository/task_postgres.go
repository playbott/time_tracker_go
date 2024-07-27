package repository

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strings"
	"time"
	"time_tracker/internal/models"
)

type TaskPostgres struct {
	db *gorm.DB
}

func NewTaskPostgres(db *gorm.DB) *TaskPostgres {
	return &TaskPostgres{db: db}
}

func (r *TaskPostgres) CreateAndStart(userId uint, taskTitle string) (uint, error) {
	var err error
	var newTaskId uint
	query := fmt.Sprintf(`insert into %s (user_id, title, started_at, created_at, updated_at)
values (%d, '%s', now(), now(), now()) RETURNING id;`,
		models.TasksTable, userId, taskTitle)
	if err = r.db.Raw(query).Scan(&newTaskId).Error; err != nil {
		logrus.Errorf("Failed to create task: %s", err.Error())
		return 0, err
	}
	return newTaskId, nil
}

func (r *TaskPostgres) Complete(userId uint, taskId uint) error {
	query := fmt.Sprintf(`update %s set completed_at = now() where user_id = %d and id = %d;`,
		models.TasksTable, userId, taskId)
	if err := r.db.Raw(query).Error; err != nil {
		logrus.Errorf("Failed to complete task: %s", err.Error())
		return err
	}
	return nil
}

func (r *TaskPostgres) Get(request models.TasksGetRequest) ([]models.Task2, error) {

	parseDateTime := func(dateTime string) (time.Time, error) {
		return time.Parse(time.RFC3339, dateTime)
	}

	pageNumber := request.Page.Number
	pageSize := request.Page.Size

	startDateTimeFrom, startDateTimeFromErr := parseDateTime(request.Filter.StartedAtFrom)
	startDateTimeTo, startDateTimeToErr := parseDateTime(request.Filter.StartedAtTo)
	completedDateTimeFrom, completedDateTimeFromErr := parseDateTime(request.Filter.CompletedAtFrom)
	completedDateTimeTo, completedDateTimeToErr := parseDateTime(request.Filter.CompletedAtTo)
	if startDateTimeFromErr != nil || startDateTimeToErr != nil || completedDateTimeFromErr != nil || completedDateTimeToErr != nil {
		logrus.Errorf("Error parsing datetime in filter fields")
	}

	fmt.Print(
		startDateTimeFrom, "\n",
		startDateTimeTo, "\n",
		completedDateTimeFrom, "\n",
		completedDateTimeTo, "\n",
	)
	if startDateTimeFromErr == nil && startDateTimeToErr == nil {
		if startDateTimeFrom.After(startDateTimeTo) {
			logrus.Errorf("Start date must be before end date")
			return nil, errors.New("start date must be before end date")
		}
	}

	if completedDateTimeFromErr == nil && completedDateTimeToErr == nil {
		if completedDateTimeFrom.After(completedDateTimeTo) {
			logrus.Errorf("Completed date must be before end date")
			return nil, errors.New("completed date must be before end date")
		}
	}

	params := []string{"deleted_at is null"}

	if request.Filter.Title != "" {
		params = append(params, fmt.Sprintf("title ilike '%%%s%%'", request.Filter.Title))
	}

	if startDateTimeFromErr == nil {
		params = append(params, fmt.Sprintf("date(started_at) >= '%s'", startDateTimeFrom.Format(time.RFC3339)))
	}

	if startDateTimeToErr == nil {
		params = append(params, fmt.Sprintf("date(started_at) <= '%s'", startDateTimeTo.Format(time.RFC3339)))
	}

	if completedDateTimeFromErr == nil {
		params = append(params, fmt.Sprintf("date(completed_at) >= '%s'", completedDateTimeFrom.Format(time.RFC3339)))
	}

	if completedDateTimeToErr == nil {
		params = append(params, fmt.Sprintf("date(completed_at) <= '%s'", completedDateTimeTo.Format(time.RFC3339)))
	}

	query := fmt.Sprintf(`
select id,
       title,
       started_at,
       completed_at,
       user_id,
       CASE
           WHEN started_at IS NULL OR completed_at IS NULL THEN NULL
           ELSE EXTRACT(EPOCH FROM completed_at - started_at)::integer
           END AS duration_seconds
from %s
where %s
limit %d offset %d
             `, models.TasksTable,
		strings.Join(params, " and "),
		pageSize, (pageNumber-1)*pageSize)
	fmt.Printf("query: %s\n", query)
	tasks := make([]models.Task2, 0)
	if err := r.db.Raw(query).Scan(&tasks).Error; err != nil {
		logrus.Errorf("Failed to get tasks: %v", err)
	}
	return tasks, nil
}

func (r *TaskPostgres) GetByID(userId string, completed bool, durationDesc bool) ([]models.Task2, bool, error) {
	var params = []string{
		"deleted_at is null",
		fmt.Sprintf("user_id = %s", userId),
	}

	if completed {
		params = append(params, "completed_at is not null")
	}

	orderBy := "order by duration_seconds"
	if durationDesc {
		orderBy += " desc"
	}

	var tasks []models.Task2
	query := fmt.Sprintf(`
select id,
title,
started_at,
completed_at,
user_id,
CASE
   WHEN started_at IS NULL OR completed_at IS NULL THEN NULL
   ELSE EXTRACT(EPOCH FROM completed_at - started_at)::integer
   END AS duration_seconds
from %s where %s
%s
`, models.TasksTable, strings.Join(params, " and "), orderBy)
	if err := r.db.Raw(query).Scan(&tasks).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			logrus.Error("Task not found ID: ", userId)
			return tasks, true, err
		}
		logrus.Error("Failed to get task by id: ", err.Error())
		return tasks, false, err
	}

	return tasks, false, nil
}
