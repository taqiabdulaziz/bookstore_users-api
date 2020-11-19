package users

import (
	"fmt"
	"github.com/taqiabdulaziz/bookstore_users-api/datasources/mysql/users_db"
	"github.com/taqiabdulaziz/bookstore_users-api/logger"
	"github.com/taqiabdulaziz/bookstore_users-api/utils"
)

const (
	queryInsertUser       = "INSERT INTO users(first_name, last_name, email, date_created, status, password) VALUES(?, ?, ?, ?, ?, ?);"
	queryGetUser          = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE id = ?;"
	queryUpdateUser       = "UPDATE users SET first_name = ?, last_name = ?, email = ? WHERE id = ?;"
	queryDeleteUser       = "DELETE FROM users WHERE id = ?"
	queryFindUserByStatus = "SELECT id, first_name, last_name, email, date_created, status from users WHERE status = ?;"
)

func (user *User) Get() *utils.RestErr {
	stmt, err := users_db.Client.Prepare(queryGetUser)
	if err != nil {
		logger.Error("error when trying to prepare get user statement", err)
		return utils.NewInternalServerError("database error")
	}

	result := stmt.QueryRow(user.Id)
	if err := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); err != nil {
		logger.Error("error when trying to prepare get user by id", err)
		return utils.NewInternalServerError("database error")
	}
	return nil
}

func (user *User) Save() *utils.RestErr {
	stmt, err := users_db.Client.Prepare(queryInsertUser)
	if err != nil {
		logger.Error("error when trying save user", err)
		return utils.NewInternalServerError("database error")
	}

	defer stmt.Close()

	insertResult, saveErr := stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated, user.Status, user.Password)
	if saveErr != nil {
		logger.Error("error when trying to save user", err)
		return utils.NewInternalServerError("database error")
	}

	userId, err := insertResult.LastInsertId()
	if err != nil {
		logger.Error("error when trying to get last insert id after creating new user", err)
		return utils.NewInternalServerError("database error")
	}

	user.Id = userId
	return nil
}

func (user *User) Update() *utils.RestErr {
	stmt, err := users_db.Client.Prepare(queryUpdateUser)
	if err != nil {
		logger.Error("error when trying to prepare statement update", err)
		return utils.NewInternalServerError("database error")
	}

	defer stmt.Close()

	_, err = stmt.Exec(user.FirstName, user.LastName, user.Email, user.Id)
	if err != nil {
		logger.Error("error when trying to exec update user", err)
		return utils.NewInternalServerError("database error")
	}

	return nil
}

func (user *User) Delete() *utils.RestErr {
	stmt, err := users_db.Client.Prepare(queryDeleteUser)
	if err != nil {
		logger.Error("error when trying to prepare delete user", err)
		return utils.NewInternalServerError("database error")
	}
	defer stmt.Close()
	if _, err = stmt.Exec(user.Id); err != nil {
		logger.Error("error when trying to exec delete user", err)
		return utils.NewInternalServerError("database error")
	}
	return nil
}

func (user *User) FindByStatus(status string) (Users, *utils.RestErr) {
	stmt, err := users_db.Client.Prepare(queryFindUserByStatus)
	if err != nil {
		logger.Error("error when trying to prepare find by status user", err)
		return nil, utils.NewInternalServerError("database error")
	}

	defer stmt.Close()

	rows, err := stmt.Query(status)
	if err != nil {
		logger.Error("error when trying to query find by status user", err)
		return nil, utils.NewInternalServerError("database error")
	}

	defer rows.Close()

	results := make([]User, 0)
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); err != nil {
			logger.Error("error when trying to scan rows find by status user", err)
			return nil, utils.NewInternalServerError("database error")
		}

		results = append(results, user)
	}

	if len(results) == 0 {
		return nil, utils.NewNotFoundError(fmt.Sprintf("no users matching status: %s", status))
	}
	return results, nil
}
