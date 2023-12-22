package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/s6352410016/go-fiber-gorm-rest-api-crud-multiple-upload-files-postgresql/config"
	"github.com/s6352410016/go-fiber-gorm-rest-api-crud-multiple-upload-files-postgresql/database"
	"github.com/s6352410016/go-fiber-gorm-rest-api-crud-multiple-upload-files-postgresql/routes"
)

func main() {
	config.LoadEnv()
	database.ConnectDB()

	app := fiber.New()
	routes.SetUpRoutes(app)

	app.Listen(":8080")
}
