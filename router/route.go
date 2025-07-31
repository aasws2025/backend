package router

import (
	"api/handler/event"
	"api/handler/user"
	"api/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	api.Get("/docs/*", swagger.HandlerDefault)

	api.Post("/register", user.CreateUser)
	api.Post("/login", user.Authorize)

	api.Get("/event", event.GetAllEvent)
	api.Get("/event/:id", event.GetEventID)

	protected := api.Group("/protected")
	protected.Use(middleware.JWTAuthMiddleware)

	protected.Post("/event", event.TambahEvent)
	protected.Put("/event/:id", event.EditEvent)
	protected.Delete("/event/:id", event.DeleteEvent)
}
