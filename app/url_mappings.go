package app

import (
	"github.com/taqiabdulaziz/bookstore_users-api/controllers/ping"
	"github.com/taqiabdulaziz/bookstore_users-api/controllers/users"
)

func MapUrls() {
	router.GET("/ping", ping.Ping)
	router.GET("/users/:user_id", users.GetUser)
	router.POST("/users", users.CreateUser)
}
