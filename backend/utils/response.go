package utils

import "github.com/gofiber/fiber/v2"

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
}

type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func SuccessResponse(c *fiber.Ctx, data interface{}) error {
	return c.JSON(Response{
		Success: true,
		Data:    data,
	})
}

func ErrorResponse(c *fiber.Ctx, status int, code, message string) error {
	return c.Status(status).JSON(Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
		},
	})
}

func ValidationError(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusBadRequest, "VALIDATION_ERROR", message)
}

func NotFoundError(c *fiber.Ctx, resource string) error {
	return ErrorResponse(c, fiber.StatusNotFound, "NOT_FOUND", resource+" not found")
}

func ConflictError(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusConflict, "CONFLICT", message)
}

func InternalError(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusInternalServerError, "INTERNAL_ERROR", message)
}
