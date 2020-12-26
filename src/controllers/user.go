package controllers

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gustcorrea/simple_auth_api/database"
	"github.com/gustcorrea/simple_auth_api/models"
)

func CreateUser(context *fiber.Ctx) error {
	user := new(models.User)
	context.BodyParser(user)
	errorResponse := models.RegisterUser(user)
	if errorResponse != nil {
		return context.Status(409).JSON(errorResponse)
	}
	return context.JSON(user)
}

func getUser(userId string) (user *models.User, errorResponse *database.ErrorResponse) {
	user = new(models.User)
	id, _ := strconv.Atoi(userId)
	errorResponse = models.GetUserByID(uint(id), user)
	return user, errorResponse
}

func GetUser(context *fiber.Ctx) error {
	userId := fmt.Sprintf("%v", context.Locals("userID"))
	user, errorResponse := getUser(userId)
	if errorResponse != nil {
		return context.Status(404).JSON(errorResponse)
	}
	return context.JSON(user)
}

func EditUser(context *fiber.Ctx) error {
	userId := fmt.Sprintf("%v", context.Locals("userID"))
	user, errorResponse := getUser(userId)
	newUser := new(models.EditableUser)
	context.BodyParser(newUser)
	if errorResponse != nil {
		return context.Status(404).JSON(errorResponse)
	}
	models.EditUser(user, newUser)
	return context.JSON(user)
}

func ChangePassword(context *fiber.Ctx) error {
	userId := fmt.Sprintf("%v", context.Locals("userID"))
	user, errorResponse := getUser(userId)
	if errorResponse != nil {
		return context.Status(404).JSON(errorResponse)
	}
	password := new(models.PasswordBody)
	context.BodyParser(password)
	models.ChangePassword(user, password.Password)
	return context.JSON(fiber.Map{
		"result": true,
	})
}
