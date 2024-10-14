package main

import (
	"time"

	"gofiber-cleanarch-test/internal/infrastructure/database"
	"gofiber-cleanarch-test/internal/infrastructure/repository"
	"gofiber-cleanarch-test/internal/interfaces/http/controllers"
	"gofiber-cleanarch-test/internal/interfaces/http/middleware"
	"gofiber-cleanarch-test/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	database.ConnectDB()

	// repo init
	userRepo := repository.NewUserRepository()

	// service init
	userService := service.NewUserService(userRepo, database.DB)
	authService := service.NewAuthService(userRepo, database.DB)

	// controller init
	userController := controllers.NewUserController(userService)
	authController := controllers.NewAuthController(authService)

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			return c.Status(code).JSON(fiber.Map{
				"error":   true,
				"message": err.Error(),
			})
		},
	})

	app.Use(cors.New())
	app.Use(helmet.New())
	app.Use(logger.New())
	app.Use(limiter.New(limiter.Config{
		Max:               10,
		Expiration:        1 * time.Minute,
		LimiterMiddleware: limiter.SlidingWindow{}, // sliding window rate limiter,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"errors":  true,
				"message": "Too many requests, please try again later.",
			})
		},
	}))

	// app.Get("/monitor", monitor.New()) // still beta on fiber

	// routing
	api := app.Group("/api")
	v1 := api.Group("/v1")
	// below will be the endpoint with prefix /api/v1
	v1.Get("/users", middleware.IsAuth, middleware.IsSuperAdmin, userController.GetAllUsers)
	v1.Get("/users/:id", middleware.IsAuth, middleware.IsSuperAdminOrIsSelf, userController.GetUserById)
	v1.Post("/users", middleware.IsAuth, middleware.IsSuperAdmin, userController.CreateUser)
	v1.Patch("/users/:id", middleware.IsAuth, middleware.IsSuperAdminOrIsSelf, userController.EditUser)
	v1.Patch("/users/:id/password", middleware.IsAuth, middleware.IsSelf, userController.EditUserPassword)
	v1.Delete("/users/:id", middleware.IsAuth, middleware.IsSuperAdmin, userController.DeleteUser)

	v1.Post("/login", authController.Login)

	app.Listen(":3000")
}
