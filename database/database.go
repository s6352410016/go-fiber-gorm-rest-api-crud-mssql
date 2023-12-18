package database

import (
	"fmt"
	"log"
	"os"

	"github.com/s6352410016/go-fiber-gorm-rest-api-crud-mssql/models"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := fmt.Sprintf(
		"sqlserver://%s:%s@%s:%s?database=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed To Connect Database\n", err.Error())
		os.Exit(2)
	}

	log.Println("Connected To Database Successfully")
	db.Logger = logger.Default.LogMode(logger.Info)
	db.AutoMigrate(&models.Product{})
	DB = db
}
