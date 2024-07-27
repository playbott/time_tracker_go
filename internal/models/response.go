package models

type Response[T any] struct {
	Data    T      `json:"data"`
	Message string `json:"message"`
}

type ResponseError struct {
	Message string `json:"message"`
}
