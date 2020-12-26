package database

import (
	"strconv"
	"time"

	"github.com/jackc/pgconn"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database exports access to database for models
var Database *gorm.DB

var models []interface{}

func connectToDB() error {
	dsn := "host=localhost user=postgres password=Chv5taffvs dbname=auth port=5432"

	var err error

	Database, err = gorm.Open(postgres.New(postgres.Config{DSN: dsn, PreferSimpleProtocol: true}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	return err
}

func addToPool() error {
	sqlDB, err := Database.DB()

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	sqlDB.SetConnMaxLifetime(time.Hour)

	return err
}

// RegisterModel appends models for automigration tool
func RegisterModel(model interface{}) bool {
	models = append(models, model)
	return true
}

func migrate() error {
	for _, m := range models {
		Database.AutoMigrate(m)
	}
	return nil
}

// InitDatabase runs connection, pooling and migration
func InitDatabase() {
	errors := []error{connectToDB(), addToPool(), migrate()}

	for _, err := range errors {
		if err != nil {
			panic(err)
		}
	}
}

// ErrorResponse exposes struct for response
type ErrorResponse struct {
	Error string `json:"error"`
}

// ErrorCodes return error responses by error code
func ErrorCodes(errorCode string) *ErrorResponse {
	code, _ := strconv.Atoi(errorCode)
	switch code {
	case 23505:
		return &ErrorResponse{Error: "Duplicate data"}
	default:
		return nil
	}
}

// ErrorString return error responses by error string
func ErrorString(errorString string) *ErrorResponse {
	switch errorString {
	case "record not found":
		return &ErrorResponse{Error: "Object not found"}
	default:
		return nil
	}
}

func ParseResponse(result *gorm.DB) *ErrorResponse {
	if err, ok := result.Error.(*pgconn.PgError); ok {
		return ErrorCodes(err.Code)
	}
	if result.Error != nil {
		return ErrorString(result.Error.Error())
	}
	return nil
}
