package router

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/yourusername/project-management/config"
	"github.com/yourusername/project-management/middleware"
	"github.com/yourusername/project-management/provider"
)

func Setup(app *fiber.App, p *provider.Provider, cfg *config.Config) {
	allowOrigins := "*"
	if len(p.Config.CORS.Origins) > 0 {
		allowOrigins = strings.Join(p.Config.CORS.Origins, ",")
	}

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: allowOrigins,
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, PATCH, DELETE, OPTIONS",
	}))

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// API routes
	api := app.Group("/api")

	// Auth routes (public - no auth middleware)
	authHandler := NewAuthHandler(p)
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)

	// Protected auth routes
	authProtected := auth.Group("", middleware.Auth(cfg))
	authProtected.Get("/me", authHandler.GetMe)
	authProtected.Post("/change-password", authHandler.ChangePassword)

	// Apply auth middleware to all other routes
	api.Use(middleware.Auth(cfg))

	// Initialize handlers
	projectHandler := NewProjectHandler(p)
	issueHandler := NewIssueHandler(p)
	sprintHandler := NewSprintHandler(p)
	commentHandler := NewCommentHandler(p)
	activityHandler := NewActivityHandler(p)
	searchHandler := NewSearchHandler(p)
	wsHandler := NewWebSocketHandler(p)
	workflowHandler := NewWorkflowHandler(p)
	notificationHandler := NewNotificationHandler(p)
	userHandler := NewUserHandler(p)

	// Project routes
	projects := api.Group("/projects")
	projects.Post("/", projectHandler.Create)
	projects.Get("/", projectHandler.List)
	projects.Get("/:id", projectHandler.GetByID)
	projects.Patch("/:id", projectHandler.Update)
	projects.Delete("/:id", projectHandler.Delete)

	// Project member routes
	projects.Get("/:id/members", projectHandler.GetMembers)
	projects.Post("/:id/members", projectHandler.AddMember)
	projects.Delete("/:id/members/:userId", projectHandler.RemoveMember)

	// User routes
	users := api.Group("/users")
	users.Get("/", userHandler.List)
	users.Post("/", userHandler.Create)

	// Workflow routes
	projects.Get("/:id/workflow", workflowHandler.GetByProject)
	projects.Post("/:id/workflow", workflowHandler.CreateDefault)

	// Issue routes
	projects.Post("/:id/issues", issueHandler.Create)
	projects.Get("/:id/board", issueHandler.GetBoard)
	api.Get("/issues/:id", issueHandler.GetByID)
	api.Patch("/issues/:id", issueHandler.Update)
	api.Delete("/issues/:id", issueHandler.Delete)
	api.Post("/issues/:id/transitions", issueHandler.Transition)
	api.Post("/issues/:id/watch", issueHandler.AddWatcher)
	api.Delete("/issues/:id/watch", issueHandler.RemoveWatcher)

	// Sprint routes
	projects.Get("/:id/sprints", sprintHandler.List)
	api.Post("/sprints", sprintHandler.Create)
	api.Get("/sprints/:id", sprintHandler.GetByID)
	api.Patch("/sprints/:id", sprintHandler.Update)
	api.Post("/sprints/:id/start", sprintHandler.Start)
	api.Post("/sprints/:id/complete", sprintHandler.Complete)

	// Comment routes
	api.Get("/issues/:id/comments", commentHandler.List)
	api.Post("/issues/:id/comments", commentHandler.Create)
	api.Patch("/comments/:id", commentHandler.Update)
	api.Delete("/comments/:id", commentHandler.Delete)

	// Activity routes
	projects.Get("/:id/activity", activityHandler.GetProjectActivity)
	api.Get("/issues/:id/activity", activityHandler.GetIssueActivity)

	// Search routes
	api.Get("/search", searchHandler.Search)

	// Notification routes
	api.Get("/notifications", notificationHandler.List)
	api.Post("/notifications/:id/read", notificationHandler.MarkAsRead)
	api.Post("/notifications/read-all", notificationHandler.MarkAllAsRead)

	// WebSocket route
	api.Get("/ws/:projectId", wsHandler.HandleWebSocket)
}
