package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/yourusername/project-management/config"
	"github.com/yourusername/project-management/provider"
	"github.com/yourusername/project-management/router"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	// Load configuration
	cfg := config.Load()

	// Initialize provider (dependency injection)
	p, err := provider.New(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize provider: %v", err)
	}
	defer p.Close()

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    "ERROR",
					"message": err.Error(),
				},
			})
		},
	})

	// Setup routes (auth middleware applied selectively in router)
	router.Setup(app, p, cfg)

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan

		log.Println("Shutting down server...")
		app.Shutdown()
	}()

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Server starting on %s (env: %s)", addr, cfg.Server.Env)

	if err := app.Listen(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
