package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/s6352410016/go-fiber-gorm-rest-api-crud-mssql/config"
	"github.com/s6352410016/go-fiber-gorm-rest-api-crud-mssql/database"
	"github.com/s6352410016/go-fiber-gorm-rest-api-crud-mssql/routes"
)

func main() {
	config.LoadEnv()
	database.ConnectDB()

	app := fiber.New()
	routes.SetUpRoutes(app)

	app.Listen(":8080")
}
