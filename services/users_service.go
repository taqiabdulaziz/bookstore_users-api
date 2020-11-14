package services

import (
	"github.com/taqiabdulaziz/bookstore_users-api/domain/users"
	"github.com/taqiabdulaziz/bookstore_users-api/utils"
	"github.com/taqiabdulaziz/bookstore_users-api/utils/crypto_utils"
	"github.com/taqiabdulaziz/bookstore_users-api/utils/date_utils"
)

var (
	UsersService usersServiceInterface = &usersService{}
)

type usersService struct{}

type usersServiceInterface interface {
	CreateUser(users.User) (*users.User, *utils.RestErr)
	GetUser(int64) (*users.User, *utils.RestErr)
	UpdateUser(bool, users.User) (*users.User, *utils.RestErr)
	DeleteUser(int64) *utils.RestErr
	Search(string) (users.Users, *utils.RestErr)
}

func (_ *usersService) CreateUser(user users.User) (*users.User, *utils.RestErr) {
	if err := user.Validate(); err != nil {
		return nil, err
	}
	user.Status = users.StatusActive
	user.DateCreated = date_utils.GetNowDBFormat()
	user.Password = crypto_utils.GetMd5(user.Password)
	if err := user.Save(); err != nil {
		return nil, err
	}

	return &user, nil
}

func (_ *usersService) GetUser(userId int64) (*users.User, *utils.RestErr) {
	result := &users.User{
		Id: userId,
	}

	if err := result.Get(); err != nil {
		return nil, err
	}

	return result, nil
}

func (_ *usersService) UpdateUser(isPartial bool, user users.User) (*users.User, *utils.RestErr) {
	currentUser := &users.User{
		Id: user.Id,
	}

	if err := currentUser.Get(); err != nil {
		return nil, err
	}

	if err := user.Validate(); err != nil {
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

func (_ *usersService) DeleteUser(userId int64) *utils.RestErr {
	user := &users.User{Id: userId}
	return user.Delete()
}

func (_ *usersService) Search(status string) (users.Users, *utils.RestErr) {
	user := &users.User{}
	return user.FindByStatus(status)
}
