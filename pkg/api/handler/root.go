package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func RootHandler(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(struct {
		Timestamp string `json:"timestamp"`
		Message   string `json:"message"`
	}{Message: "Hello from container-chief", Timestamp: time.Now().Format(time.RFC3339)})
}
