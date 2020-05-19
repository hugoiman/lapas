package controllers

import (
	"fmt"
	models "lapas/pkg/models"
)

// User is class
type User models.User

// GetUser is function
func GetUser() {
	data := models.GetUser("10")
	fmt.Printf("%+v\n User: ", data)
}

// GetUsers is function
func GetUsers() {
	data := models.GetUsers()
	for _, v := range data {
		fmt.Printf("%+v\n User: ", v)
	}
}
