package main

import (
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	app := fiber.New(fiber.Config{
		Prefork:     true,
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
	})
	app.Use(logger.New(), recover.New())

	app.Post("/user", postUser)
	app.Delete("/user", deleteUser)
	app.Put("/user", putUser)

	app.Post("/login", postLogin)

	app.Get("/group/:groupId", getGroup)
	app.Post("/group", postGroup)
	app.Delete("/group/:groupId", deleteGroup)
	app.Put("/group/:groupId", putGroup)

	netflixApi := app.Group("/netflix")
	netflixApi.Post("/account", postNetflixAccount)
	netflixApi.Put("/account", putNetflixAccount)
	netflixApi.Delete("/membership", deleteNetflixMembership)

	wavveApi := app.Group("/wavve")
	wavveApi.Post("/account", postWavveAccount)

	app.Listen(":8000")
}
