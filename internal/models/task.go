package models

import (
	"database/sql"
	"time"
)

const TasksTable = "tasks"

type Task struct {
	Base
	Title       string       `gorm:"not null" json:"title"`
	StartedAt   time.Time    `gorm:"not null" json:"started_at"`
	CompletedAt sql.NullTime `gorm:"" json:"completed_at"`
	UserId      uint         `gorm:"not null" json:"user_id"`
}

type TasksGetFilter struct {
	Title           string `json:"title" validate:"max=255"`
	StartedAtFrom   string `json:"started_at_from" validate:"str_min_max=19 20"`
	StartedAtTo     string `json:"started_at_to" validate:"str_min_max=19 20"`
	CompletedAtFrom string `json:"completed_at_from" validate:"str_min_max=19 20"`
	CompletedAtTo   string `json:"completed_at_to" validate:"str_min_max=19 20"`
}

type TasksGetRequest struct {
	Filter TasksGetFilter `json:"filter"`
	Page   Page           `json:"page"`
}

type Task2 struct {
	Id              uint   `json:"id"`
	UserId          uint   `json:"user_id"`
	Title           string `json:"title"`
	StartedAt       string `json:"started_at"`
	CompletedAt     string `json:"completed_at"`
	DurationSeconds uint   `json:"duration_seconds"`
	DurationString  string `json:"duration_string"`
}

type TasksSearchResponse struct {
	Response[[]Task2]
	Page Page `json:"page"`
}

type TaskCreateRequest struct {
	Id     uint   `json:"id"`
	UserId uint   `gorm:"not null" json:"user_id" validate:"required,numeric,min=1"`
	Title  string `gorm:"not null" json:"task_title" validate:"required,min=1,max=255"`
}

type TaskCompleteRequest struct {
	Id     uint `json:"id"`
	UserId uint `gorm:"not null" json:"user_id" validate:"required,numeric,min=1"`
}

func (Task) TableName() string {
	return TasksTable
}
