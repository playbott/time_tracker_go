package handler

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"strconv"
	"time_tracker/internal/models"
)

// GetUsers user search with filter and pagination.
//
//	@Description	User search with filter and pagination.
//	@Summary		user search with filter and pagination
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			params	body		models.UsersGetRequest	true	"search parameters"
//	@Success		200		{object}	models.UsersResponse
//	@Failure		400		{object}	models.ResponseError
//	@Failure		500		{object}	models.ResponseError
//	@Router			/api/v1/user [get]
func (h *Handler) GetUsers(c *fiber.Ctx) error {
	var err error
	v := c.Locals("validator").(*validator.Validate)
	userGetRequest := models.UsersGetRequest{}
	if err = c.BodyParser(&userGetRequest); err != nil {
		logrus.Error("Failed to parse request body: ", err)
		return c.Status(fiber.StatusBadRequest).JSON(models.ResponseError{Message: "Cannot parse JSON input data"})
	}

	err = v.Struct(userGetRequest)
	if err != nil {
		logrus.Error(err)
		return c.Status(fiber.StatusBadRequest).JSON(models.ResponseError{Message: err.Error()})
	}

	users, err := h.services.User.Get(userGetRequest)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ResponseError{Message: "Unknown error"})
	}

	return c.JSON(models.UsersResponse{Response: models.Response[[]models.User]{Data: users}, Page: userGetRequest.Page})
}

// GetUser func gets a user by id.
//
//	@Description	Get user by id.
//	@Summary		get user by id
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"user id"
//	@Success		200	{object}	models.Response[models.User]
//	@Failure		404	{object}	string
//	@Router			/api/v1/user/{id} [get]
func (h *Handler) GetUser(c *fiber.Ctx) error {
	v := c.Locals("validator").(*validator.Validate)
	id := c.Params("id")
	err := v.Var(id, "required,numeric")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ResponseError{Message: err.Error()})
	}

	user, notfound, err := h.services.User.GetByID(id)
	if err != nil {
		if notfound {
			return c.Status(fiber.StatusNotFound).JSON(models.ResponseError{Message: "User not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ResponseError{Message: "Unknown error"})
	}

	return c.JSON(models.Response[models.User]{Data: user})
}

// CreateUser creates a new user.
//
//	@Description	Create a new user.
//	@Summary		create a new user
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			passport	body		models.UserCreateRequest	true	"series and number of the user's passport"
//	@Success		201			{object}	models.User
//	@Router			/api/v1/user [post]
func (h *Handler) CreateUser(c *fiber.Ctx) error {
	v := c.Locals("validator").(*validator.Validate)
	var err error
	var userCreateRequest models.UserCreateRequest

	if err = c.BodyParser(&userCreateRequest); err != nil {
		logrus.Error("Failed to parse request body: ", err)
		return c.Status(fiber.StatusBadRequest).JSON(models.ResponseError{Message: "Cannot parse JSON input data"})
	}

	err = v.Struct(userCreateRequest)
	if err != nil {
		logrus.Error(err)
		return c.Status(fiber.StatusBadRequest).JSON(models.ResponseError{Message: err.Error()})
	}

	passportSerie, err := strconv.Atoi(userCreateRequest.PassportNumber[:4])
	passportNumber, err2 := strconv.Atoi(userCreateRequest.PassportNumber[5:])
	if (err != nil) || (err2 != nil) {
		logrus.Error("Failed to parse passport number. Value: ", userCreateRequest.PassportNumber)
		return c.Status(fiber.StatusBadRequest).JSON(models.ResponseError{Message: "Invalid passport number"})
	}

	people, err := h.services.People.GetByPassport(passportSerie, passportNumber)
	if err != nil {
		logrus.Error("Failed to get people info: ", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(models.ResponseError{Message: "Failed to create user"})
	}

	var user models.User

	user.People = people
	user.PassportNumber = fmt.Sprintf("%d %d", passportSerie, passportNumber)

	id, err := h.services.User.Create(user)
	if err != nil {
		logrus.Info("Failed to create user: ", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(models.ResponseError{Message: "Failed to create user"})
	}
	user.ID = id
	logrus.Info("User created. ID: ", user.ID)
	return c.Status(fiber.StatusCreated).JSON(models.Response[models.User]{Data: user})
}

// UpdateUser func for updates user data.
//
//	@Description	Update a user data.
//	@Summary		update a user data
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string						true	"user id"
//	@Param			user	body		models.UserUpdateRequest	true	"User"
//	@Success		201		{object}	models.User
//	@Router			/api/v1/user/{id} [patch]
func (h *Handler) UpdateUser(c *fiber.Ctx) error {
	var err error
	v := c.Locals("validator").(*validator.Validate)
	id := c.Params("id")
	err = v.Var(id, "required,numeric")
	if err != nil {
		logrus.Error("User ID is required: ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(models.ResponseError{Message: err.Error()})
	}
	userUpdateRequest := models.UserUpdateRequest{}
	if err = c.BodyParser(&userUpdateRequest); err != nil {
		logrus.Error("Failed to parse request body: ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(models.ResponseError{Message: "Failed to parse JSON input."})
	}

	err = v.Struct(userUpdateRequest)
	if err != nil {
		logrus.Error(err)
		return c.Status(fiber.StatusBadRequest).JSON(models.ResponseError{Message: err.Error()})
	}

	notfound, err := h.services.User.Update(id, userUpdateRequest)
	if err != nil {
		if notfound {
			return c.Status(fiber.StatusNotFound).JSON(models.ResponseError{Message: "User not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ResponseError{Message: "Unknown error"})
	}

	return c.JSON(models.Response[models.UserUpdateRequest]{Data: userUpdateRequest, Message: "User updated"})
}

// DeleteUser delete user by id.
//
//	@Description	Delete a user by id.
//	@Summary		delete a user by id
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"user id"
//	@Success		201	{object}	models.User
//	@Router			/api/v1/user/{id} [delete]
func (h *Handler) DeleteUser(c *fiber.Ctx) error {
	var err error
	v := c.Locals("validator").(*validator.Validate)
	id := c.Params("id")
	err = v.Var(id, "required,numeric")
	if err != nil {
		logrus.Error("User ID is required: ", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(models.ResponseError{Message: err.Error()})
	}

	notfound, err := h.services.User.Delete(id)
	if err != nil {
		if notfound {
			return c.Status(fiber.StatusNotFound).JSON(models.ResponseError{Message: "User not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ResponseError{Message: "Unknown error"})
	}

	return c.JSON(models.Response[string]{Message: "User deleted"})
}
