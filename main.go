package main

import (
	"container-chief/pkg/control"
	"container-chief/pkg/discord"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

const DEFAULT_SERVER_PORT = "3000"

func main() {
	slog.Info("Starting container-chief")

	var server_port string

	server_port_str, server_port_set := os.LookupEnv("SERVER_PORT")

	if !server_port_set {
		slog.Info("SERVER_PORT env not set, setting default to", "port", DEFAULT_SERVER_PORT)
		server_port = DEFAULT_SERVER_PORT
	} else {
		slog.Info("SERVER_PORT env found, setting it to", "port", server_port_str)
		server_port = server_port_str
	}

	server := fiber.New()
	chiefService := control.NewChiefService()

	defer func() {
		slog.Info("Stopping container-chief")

		if err := server.Shutdown(); err == nil {
			slog.Error("Error shutting down server", "error", err)
		}

		if err := chiefService.Cli.Close(); err != nil {
			slog.Error("Error shutting down docker CLI", "error", err)
		}

		slog.Info("Exiting container-chief")
	}()

	server.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(struct {
			Timestamp string `json:"timestamp"`
			Message   string `json:"message"`
		}{Message: "Hello from container-chief", Timestamp: time.Now().Format(time.RFC3339)})
	})

	server.Post("/discord-webhook", func(c *fiber.Ctx) error {
		bodyContent := discord.InteractionWebhook{}
		err := c.BodyParser(&bodyContent)

		if err != nil {
			slog.Warn("Cannot parse body", "error", err)
		} else {
			fmt.Println(bodyContent)
		}

		fmt.Println("DOne print")

		return c.JSON(fiber.Map{
			"message": "ok",
		})
	})

	server.Listen(":" + server_port)
}
