package controllers

import (
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/gustcorrea/simple_auth_api/database"
	"github.com/gustcorrea/simple_auth_api/models"
	"github.com/gustcorrea/simple_auth_api/settings"
)

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int    `json:"expires_at"`
}
type Credentials struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type Claims struct {
	ID uint
	jwt.StandardClaims
}

type Refresh struct {
	Token string `json:"refresh_token"`
}

func (token *Token) tokenForUser(user *models.User) {
	expiration := time.Now().Add(5 * time.Minute)
	accessClaims := &Claims{
		ID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiration.Unix(),
		},
	}
	refreshClaims := &Claims{
		ID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiration.Add(24 * time.Hour).Unix(),
		},
	}
	tk := jwt.NewWithClaims(jwt.SigningMethodHS512, accessClaims)
	rfTk := jwt.NewWithClaims(jwt.SigningMethodHS512, refreshClaims)
	strTk, _ := tk.SignedString(settings.AccessKey)
	strRf, _ := rfTk.SignedString(settings.RefreshKey)
	token.AccessToken = strTk
	token.ExpiresAt = 300
	token.RefreshToken = strRf
}

func clientMatch(clientID, clientSecret string) bool {

	if clientID != settings.ClientId || clientSecret != settings.ClientSecret {
		return false
	}
	return true
}

func authenticate(username, password, clientID, clientSecret string, token *Token) *database.ErrorResponse {
	if !clientMatch(clientID, clientSecret) {
		return &database.ErrorResponse{Error: "client id/ secret did not match"}
	} else {
		user := new(models.User)
		result := models.GetUserByUsername(user, username)
		if result != nil {
			return result
		}
		if !user.VerifyPassword(password) {
			return &database.ErrorResponse{Error: "password did not match"}
		}
		token.tokenForUser(user)
		return nil
	}
}

func Authenticate(context *fiber.Ctx) error {
	token := new(Token)
	credentials := new(Credentials)
	context.BodyParser(credentials)
	errorResponse := authenticate(credentials.Username, credentials.Password, credentials.ClientId, credentials.ClientSecret, token)
	if errorResponse != nil {
		return context.Status(401).JSON(errorResponse)
	}
	return context.JSON(token)
}

func parserFunc(key []byte) func(*jwt.Token) (interface{}, error) {
	return func(token *jwt.Token) (interface{}, error) {
		return key, nil
	}
}

func validateToken(token string, claims *Claims, user *models.User, key []byte) (*database.ErrorResponse, int) {
	errorResponse := new(database.ErrorResponse)
	tokenStruct, err := jwt.ParseWithClaims(token, claims, parserFunc(key))
	if tokenStruct.Valid {
		if claims.ExpiresAt < time.Now().Unix() {
			errorResponse.Error = "Token expired"
			return errorResponse, fiber.StatusUnauthorized
		}
		models.GetUserByID(claims.ID, user)
		if user == nil {
			errorResponse.Error = "User not found"
			return errorResponse, fiber.StatusUnauthorized
		}
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			errorResponse.Error = "Invalid Token"
			return errorResponse, fiber.StatusForbidden
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			errorResponse.Error = "Invalid Token"
			return errorResponse, fiber.StatusUnauthorized
		} else {
			errorResponse.Error = "Invalid Token"
			return errorResponse, fiber.StatusForbidden
		}
	}
	return nil, 200
}

func SecureAuth(c *fiber.Ctx) error {

	authorization := strings.Split(c.Get("Authorization"), " ")
	errorResponse := new(database.ErrorResponse)
	var status int

	if len(authorization) != 2 {
		errorResponse.Error = "Invalid token format"
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse)
	} else {
		var flag bool = false
		for _, str := range settings.AllowedAuthorizationPrefix {
			if str == authorization[0] {
				flag = true
				break
			}
		}
		if !flag {
			errorResponse.Error = "Invalid token Prefix"
			return c.Status(fiber.StatusBadRequest).JSON(errorResponse)
		}
	}
	accessToken := authorization[1]
	claims := new(Claims)

	user := new(models.User)
	errorResponse, status = validateToken(accessToken, claims, user, settings.AccessKey)
	if errorResponse != nil {
		return c.Status(status).JSON(errorResponse)
	}
	// token, err := jwt.ParseWithClaims(accessToken, claims, parserFunc(settings.AccessKey))
	// if token.Valid {
	// 	if claims.ExpiresAt < time.Now().Unix() {
	// 		errorResponse.Error = "Token expired"
	// 		return c.Status(fiber.StatusUnauthorized).JSON(errorResponse)
	// 	}
	// 	models.GetUserByID(claims.ID, user)
	// 	if user == nil {
	// 		c.ClearCookie("access_token", "refresh_token")
	// 		errorResponse.Error = "User not found"
	// 		return c.Status(fiber.StatusUnauthorized).JSON(errorResponse)
	// 	}
	// } else if ve, ok := err.(*jwt.ValidationError); ok {
	// 	if ve.Errors&jwt.ValidationErrorMalformed != 0 {
	// 		c.ClearCookie("access_token", "refresh_token")
	// 		return c.SendStatus(fiber.StatusForbidden)
	// 	} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
	// 		c.ClearCookie("access_token", "refresh_token")
	// 		return c.SendStatus(fiber.StatusUnauthorized)
	// 	} else {
	// 		c.ClearCookie("access_token", "refresh_token")
	// 		return c.SendStatus(fiber.StatusForbidden)
	// 	}
	// }
	c.Locals("userID", user.ID)
	return c.Next()
}

func RefreshToken(context *fiber.Ctx) error {

	refresh := new(Refresh)
	claims := new(Claims)
	user := new(models.User)
	context.BodyParser(refresh)

	errorResponse, status := validateToken(refresh.Token, claims, user, settings.RefreshKey)
	if errorResponse != nil {
		return context.Status(status).JSON(errorResponse)
	}

	token := new(Token)
	token.tokenForUser(user)

	return context.JSON(token)
}
