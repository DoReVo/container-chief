package handler

import (
	"container-chief/pkg/discord"
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

func DiscordWebhookHandler(c *fiber.Ctx) error {
	bodyContent := discord.InteractionWebhook{}

	err := c.BodyParser(&bodyContent)
	if err != nil {
		slog.Warn("Cannot parse body", "error", err)
		return c.Status(400).JSON(fiber.Map{"ok": false})
	}

	// Pretty print the body
	slog.Info("Body content", "body", bodyContent)

	if bodyContent.Type == 1 {
		slog.Info("Ping received")

		signature := c.Get("X-Signature-Ed25519")
		timestamp := c.Get("X-Signature-Timestamp")

		isValidSignature := discordService.VerifyWebhookSignature(signature, timestamp, string(c.BodyRaw()))

		if !isValidSignature {
			return c.Status(401).JSON(fiber.Map{"ok": false})
		}

		return c.JSON(fiber.Map{"type": 1})
	} else if bodyContent.Type == 2 {
		// Send an interaction response
		token := bodyContent.Token
		id := bodyContent.ID

		response := discord.InteractionResponse{
			Type: discord.ChannelMessageWithSource,
			Data: discord.InteractionCallbackData{
				Content: "Received",
			},
		}

		if err := discordService.RespondInteraction(id, token, response); err != nil {
			slog.Warn("Error trying to respond to interaction", "error", err)
		}

	}

	return c.JSON(fiber.Map{
		"message": "ok",
	})
}
