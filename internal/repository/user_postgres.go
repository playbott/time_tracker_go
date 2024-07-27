package repository

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strings"
	"time_tracker/internal/models"
)

type UserPostgres struct {
	db *gorm.DB
}

func NewUserPostgres(db *gorm.DB) *UserPostgres {
	return &UserPostgres{db: db}
}

func (r *UserPostgres) Get(request models.UsersGetRequest) ([]models.User, error) {
	pageNumber := request.Page.Number
	pageSize := request.Page.Size

	ilike := func(s string) string {
		return "'%" + s + "%'"
	}

	query := fmt.Sprintf(`select * from %s
         where deleted_at is null
			 and name ilike %s
             and surname ilike %s
             and patronymic ilike %s
             and passport_number ilike %s
             and address ilike %s
				limit %d
				offset %d
             `, models.UsersTable,
		ilike(request.Filter.Name),
		ilike(request.Filter.Surname),
		ilike(request.Filter.Patronymic),
		ilike(request.Filter.PassportNumber),
		ilike(request.Filter.Address),
		pageSize, (pageNumber-1)*pageSize)

	users := make([]models.User, 0)
	if err := r.db.Raw(query).Scan(&users).Error; err != nil {
		logrus.Errorf("Failed to get users: %v", err)
	}
	return users, nil
}

func (r *UserPostgres) GetByID(id string) (models.User, bool, error) {
	var user models.User
	query := fmt.Sprintf(`select * from %s where deleted_at is null
			 and id = %s`, models.UsersTable, id)
	if err := r.db.Raw(query).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			logrus.Error("User not found ID: ", id)
			return user, true, err
		}
		logrus.Error("Failed to get user by id: ", err.Error())
		return user, false, err
	}

	return user, false, nil
}

func (r *UserPostgres) Update(id string, user models.UserUpdateRequest) (notFound bool, err error) {
	params := make([]string, 0)
	if user.Name != "" {
		params = append(params, fmt.Sprintf("name = '%s'", user.Name))
	}
	if user.Surname != "" {
		params = append(params, fmt.Sprintf("surname = '%s'", user.Surname))
	}
	if user.Patronymic != "" {
		params = append(params, fmt.Sprintf("patronymic = '%s'", user.Patronymic))
	}
	if user.Address != "" {
		params = append(params, fmt.Sprintf("address = '%s'", user.Address))
	}
	if user.PassportNumber != "" {
		params = append(params, fmt.Sprintf("passport_number = '%s'", user.PassportNumber))
	}

	if len(params) > 0 {
		tx := r.db.Begin()
		var count int64 = 0
		_ = tx.Table(models.UsersTable).Where("id = ?", id).Count(&count).Error
		if err == gorm.ErrRecordNotFound || count == 0 {
			logrus.Error("User not found ID: ", id)
			return true, errors.New("user not found")
		}

		query := fmt.Sprintf(`update %s set %s where id = %s;`, models.UsersTable, strings.Join(params, ", "), id)
		if err = tx.Exec(query).Error; err != nil {
			tx.Rollback()
			logrus.Errorf("Failed to update user data ID: %s. Error: %s", id, err.Error())
			return false, err
		}
		tx.Commit()
	}

	return false, nil
}

func (r *UserPostgres) Delete(id string, deleteRecord bool) (notFound bool, err error) {
	var query string
	tx := r.db.Begin()
	var count int64 = 0
	err = tx.Table(models.UsersTable).Where("id = ? and deleted_at is null", id).Count(&count).Error
	if err == gorm.ErrRecordNotFound || count == 0 {
		logrus.Error("User not found ID: ", id)
		return true, errors.New("user not found")
	}

	if deleteRecord {
		query = fmt.Sprintf(`delete from %s where id = %s;`, models.UsersTable, id)
	} else {
		query = fmt.Sprintf(`update %s set deleted_at = now() where id = %s;`, models.UsersTable, id)
	}

	if err = tx.Exec(query).Error; err != nil {
		tx.Rollback()
		logrus.Errorf("Failed to update user data ID: %s. Error: %s", id, err.Error())
		return false, err
	}
	tx.Commit()
	return false, nil
}

func (r *UserPostgres) Create(user models.User) (uint, error) {
	var err error
	var newUserId uint
	query := fmt.Sprintf(`insert into %s (name, surname, patronymic, address, passport_number)
values ('%s', '%s', '%s', '%s', '%s')  RETURNING id;`,
		models.UsersTable, user.Name, user.Surname, user.Patronymic, user.Address, user.PassportNumber)
	if err = r.db.Raw(query).Scan(&newUserId).Error; err != nil {
		logrus.Errorf("Failed to create user data: %s", err.Error())
		return 0, err
	}
	return newUserId, nil
}
