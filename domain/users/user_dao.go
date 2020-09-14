package users

import (
	"fmt"
	"github.com/taqiabdulaziz/bookstore_users-api/datasources/mysql/users_db"
	"github.com/taqiabdulaziz/bookstore_users-api/utils"
	"github.com/taqiabdulaziz/bookstore_users-api/utils/date_utils"
	"strings"
)

var (
	usersDB = make(map[int64]*User)
)

const (
	indexUniqueEmail = "users_email_uindex"
	queryInsertUser = "INSERT INTO users(first_name, last_name, email, date_created) VALUES(?, ?, ?, ?);"
)

func (user *User) Get() *utils.RestErr {
	if err := users_db.Client.Ping(); err != nil {
		panic(err)
	}

	result := usersDB[user.Id]
	if result == nil {
		return utils.NewNotFoundError(fmt.Sprintf("user %d not found", user.Id))
	}

	user.Id = result.Id
	user.Email = result.Email
	user.FirstName = result.FirstName
	user.LastName = result.LastName
	user.DateCreated = result.DateCreated
	return nil
}

func (user *User) Save() *utils.RestErr {
	stmt, err :=  users_db.Client.Prepare(queryInsertUser)
	if err != nil {
		return utils.NewInternalServerError(err.Error())
	}

	defer stmt.Close()
	user.DateCreated = date_utils.GetNowString()
	insertResult, err := stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated)
	if err != nil {
		if strings.Contains(err.Error(), indexUniqueEmail) {
			return utils.NewBadRequestError(fmt.Sprintf("email %s already exists", user.Email))
		}
		return utils.NewInternalServerError(fmt.Sprintf("error when trying to save user: %s", err.Error()))
	}

	userId, err := insertResult.LastInsertId()
	if err != nil {
		return utils.NewInternalServerError(fmt.Sprintf("error when trying to save user: %s", err.Error()))
	}

	if usersDB[user.Id] != nil {
		if usersDB[user.Id].Email == user.Email {
			return utils.NewBadRequestError(fmt.Sprintf("user %s already registerred", user.Email))
		}
		return utils.NewBadRequestError(fmt.Sprintf("user %d already exists", user.Id))
	}

	user.Id = userId
	return nil
}
