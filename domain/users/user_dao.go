package users

import (
	"fmt"
	"github.com/taqiabdulaziz/bookstore_users-api/datasources/mysql/users_db"
	"github.com/taqiabdulaziz/bookstore_users-api/utils"
	"github.com/taqiabdulaziz/bookstore_users-api/utils/date_utils"
	"github.com/taqiabdulaziz/bookstore_users-api/utils/mysql_utils"
)

const (
	queryInsertUser = "INSERT INTO users(first_name, last_name, email, date_created) VALUES(?, ?, ?, ?);"
	queryGetUser    = "SELECT id, first_name, last_name, email, date_created FROM users WHERE id = ?;"
	queryUpdateUser = "UPDATE users SET first_name = ?, last_name = ?, email = ? WHERE id = ?;"
)

func (user *User) Get() *utils.RestErr {
	stmt, err := users_db.Client.Prepare(queryGetUser)
	if err != nil {
		return utils.NewInternalServerError(err.Error())
	}

	result := stmt.QueryRow(user.Id)
	if err := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated); err != nil {
		return mysql_utils.ParseError(err)
	}
	return nil
}

func (user *User) Save() *utils.RestErr {
	stmt, err := users_db.Client.Prepare(queryInsertUser)
	if err != nil {
		return utils.NewInternalServerError(err.Error())
	}

	defer stmt.Close()
	user.DateCreated = date_utils.GetNowString()
	insertResult, saveErr := stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated)
	if saveErr != nil {
		return mysql_utils.ParseError(saveErr)
	}

	userId, err := insertResult.LastInsertId()
	if err != nil {
		return utils.NewInternalServerError(fmt.Sprintf("error when trying to save user: %s", err.Error()))
	}

	user.Id = userId
	return nil
}

func (user *User) Update() *utils.RestErr {
	stmt, err := users_db.Client.Prepare(queryUpdateUser)
	if err != nil {
		return utils.NewInternalServerError(err.Error())
	}

	defer stmt.Close()

	_, err = stmt.Exec(user.FirstName, user.LastName, user.Email, user.Id)
	if err != nil {
		return mysql_utils.ParseError(err)
	}

	return nil
}
