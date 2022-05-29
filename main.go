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

	app.Post("/user", addUser)
	app.Delete("/user", delUser)
	app.Put("/user", setUser)

	app.Post("/login", login)

	app.Get("/group/:groupId", getGroup)
	app.Post("/group", addGroup)
    app.Delete("/group/:groupId", delGroup)
	app.Put("/group/:groupId", setGroup)

	netflixApi := app.Group("/netflix")
	netflixApi.Post("/info", netflixInfo)
	netflixApi.Delete("/membership", netflixUnsubscribe)

	wavveApi := app.Group("/wavve")
	wavveApi.Post("/info", wavveInfo)

	app.Listen(":8000")
}
