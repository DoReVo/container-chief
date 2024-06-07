package discord

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v2"
)

func NewDiscordService() *DiscordService {
	appPublicKey := os.Getenv("DISCORD_APP_PUBLIC_KEY")
	appId := os.Getenv("DISCORD_APP_ID")
	botToken := os.Getenv("DISCORD_BOT_TOKEN")

	if appPublicKey == "" {
		panic("Discord app public key not defined in environment variable")
	}

	if appId == "" {
		panic("Discord app ID not defined in environment variable")
	}

	if botToken == "" {
		panic("Discord bot token not defined in environment variable")
	}

	service := DiscordService{
		AppPublicKey: appPublicKey,
		AppID:        appId,
		BotToken:     botToken,
	}

	err := service.SetGlobalCommand()
	if err != nil {
		panic(err)
	}

	return &service
}

func (ds *DiscordService) VerifyWebhookSignature(signature string, timestamp string, requestBody string) bool {
	msg := timestamp + (requestBody)

	decodedSignature, err := hex.DecodeString(signature)
	if err != nil {
		return false
	}

	decodedPublicKey, err := hex.DecodeString(ds.AppPublicKey)
	if err != nil {
		return false
	}

	return ed25519.Verify(decodedPublicKey, []byte(msg), decodedSignature)
}

func (ds *DiscordService) SetGlobalCommand() error {
	globalCommands := []Command{
		{
			Name:        "list-containers",
			Description: "List all running containers",
			Type:        1,
			// Options: []CommandOption{
			// 	{
			// 		Type: StringType,
			// 		Name: "lol",
			// 	},
			// },
		},
	}

	url := fmt.Sprintf("%s/applications/%s/commands", DISCORD_BASE_API_URL, ds.AppID)

	if err := ds.SendRequest(fiber.MethodPut, url, globalCommands); err != nil {
		slog.Warn("Error when trying to set global command", "error", err)
		return err
	}

	return nil
}

func (ds *DiscordService) RespondInteraction(interactionId string, interactionToken string, response InteractionResponse) error {
	url := fmt.Sprintf("%v/interactions/%v/%v/callback", DISCORD_BASE_API_URL, interactionId, interactionToken)

	if err := ds.SendRequest(fiber.MethodPost, url, response); err != nil {
		slog.Warn("Error when trying to respond to interaction", "error", err)
		return err
	}

	return nil
}

func (ds *DiscordService) SendRequest(method string, url string, body interface{}) error {
	agent := fiber.AcquireAgent()
	agent.Set("Authorization", "Bot "+ds.BotToken)
	agent.Request().Header.SetMethod(method)
	agent.Request().SetRequestURI(url)
	agent.JSON(body)

	if err := agent.Parse(); err != nil {
		slog.Error("Cannot parse agent", "error", "err")
		return err
	}

	code, body, errs := agent.String()

	if len(errs) > 0 {
		slog.Warn("Error in sending request to Discord", "error", errs)
		return errs[len(errs)-1]
	}

	if code < 200 || code >= 300 {
		slog.Warn("Error response from discord", "response", body)
		return ErrFailedSettingUpGlobalCommand
	}

	slog.Info("Response from discord", "code", code, "response", body)

	return nil
}
