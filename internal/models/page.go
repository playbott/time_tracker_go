package models

type Page struct {
	Number int `json:"number" validate:"required,numeric,min=1"`
	Size   int `json:"size" validate:"required,numeric,min=1"`
}
