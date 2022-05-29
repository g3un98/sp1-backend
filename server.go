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

	netflixApi := app.Group("/netflix")
	netflixApi.Post("/info", netflixInfo)

	wavveApi := app.Group("/wavve")
	wavveApi.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, Wavve!")
	})

	app.Listen(":8000")
}
