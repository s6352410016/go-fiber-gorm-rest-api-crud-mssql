package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/s6352410016/go-fiber-gorm-rest-api-crud-mssql/handlers"
)

func SetUpRoutes(app *fiber.App) {
	product := app.Group("/api")
	product.Get("/products", handlers.GetAll)
	product.Get("/product/:id", handlers.GetById)
	product.Post("/create", handlers.Create)
	product.Put("/update/:id", handlers.Update)
	product.Delete("/delete/:id", handlers.Delete)
	product.Get("/image/:filename", handlers.GetImage)
}
