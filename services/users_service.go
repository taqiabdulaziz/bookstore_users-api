package services

import (
	"github.com/taqiabdulaziz/bookstore_users-api/domain/users"
	"github.com/taqiabdulaziz/bookstore_users-api/utils"
	"github.com/taqiabdulaziz/bookstore_users-api/utils/date_utils"
)

func CreateUser(user users.User) (*users.User, *utils.RestErr) {
	if err := user.Validate(); err != nil {
		return nil, err
	}
	user.Status = users.StatusActive
	user.DateCreated = date_utils.GetNowDBFormat()
	if err := user.Save(); err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUser(userId int64) (*users.User, *utils.RestErr) {
	result := &users.User{
		Id: userId,
	}

	if err := result.Get(); err != nil {
		return nil, err
	}

	return result, nil
}

func UpdateUser(isPartial bool, user users.User) (*users.User, *utils.RestErr) {
	currentUser, err := GetUser(user.Id)
	if err != nil {
		return nil, err
	}

	if err = user.Validate(); err != nil {
		return nil, err
	}

	if !isPartial {
		currentUser.FirstName = user.FirstName
		currentUser.LastName = user.LastName
		currentUser.Email = user.Email
	} else {
		if user.FirstName != "" {
			currentUser.FirstName = user.FirstName
		}
		if user.LastName != "" {
			currentUser.LastName = user.LastName
		}
		if user.Email != "" {
			currentUser.Email = user.Email
		}
	}

	if err := currentUser.Update(); err != nil {
		return nil, err
	}

	return currentUser, nil
}

func DeleteUser(userId int64) *utils.RestErr {
	user := &users.User{Id: userId}
	return user.Delete()
}

func Search(status string) ([]users.User, *utils.RestErr) {
	user := &users.User{}
	return user.FindByStatus(status)
}
