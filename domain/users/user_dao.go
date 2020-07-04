package users

import (
	"fmt"
	"github.com/taqiabdulaziz/bookstore_users-api/utils"
)

var (
	usersDB = make(map[int64]*User)
)

func (user *User) Get() *utils.RestErr {
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
	if usersDB[user.Id] != nil {
		if usersDB[user.Id].Email == user.Email {
			return utils.NewBadRequestError(fmt.Sprintf("user %s already registerred", user.Email))
		}
		return utils.NewBadRequestError(fmt.Sprintf("user %d already exists", user.Id))
	}

	usersDB[user.Id] = user
	return nil
}
