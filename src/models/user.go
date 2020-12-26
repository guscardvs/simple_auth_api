package models

import (
	"github.com/gustcorrea/simple_auth_api/database"
	"golang.org/x/crypto/bcrypt"
)

// User Model Definition
type User struct {
	ID        uint   `gorm:"primaryKey;" json:"id"`
	Username  string `json:"username" gorm:"unique"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Password  string `json:"-"`
}

// EditableUser defines fields available for edit
type EditableUser struct {
	Username  string `json:"username" gorm:"unique"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type PasswordBody struct {
	Password string `json:"password"`
}

var n = database.RegisterModel(&User{})

func (user *User) parsePassword() error {
	pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), -1)
	user.Password = string(pass)
	return err
}

func (user *User) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}

func RegisterUser(user *User) *database.ErrorResponse {
	user.parsePassword()
	result := database.Database.Create(user)
	return database.ParseResponse(result)
}

func GetUserByID(id uint, user *User) *database.ErrorResponse {
	result := database.Database.First(&user, id)
	return database.ParseResponse(result)
}

func EditUser(oldUser *User, newUser *EditableUser) *database.ErrorResponse {
	EditModel(oldUser, newUser)
	result := database.Database.Save(oldUser)
	return database.ParseResponse(result)

}

func ChangePassword(user *User, password string) *database.ErrorResponse {
	user.Password = password
	user.parsePassword()
	result := database.Database.Save(user)
	return database.ParseResponse(result)
}

func GetUserByUsername(user *User, username string) *database.ErrorResponse {
	result := database.Database.Where("username = ?", username).First(user)
	return database.ParseResponse(result)
}
