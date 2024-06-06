package discord

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v2"
)

const DISCORD_BASE_API_URL = "https://discord.com/api"

type InteractionWebhook struct {
	AuthorizingIntegrationOwners map[string]interface{} `json:"authorizing_integration_owners"`
	AppPermissions               string                 `json:"app_permissions"`
	ApplicationID                string                 `json:"application_id"`
	ID                           string                 `json:"id"`
	Token                        string                 `json:"token"`
	Entitlements                 []interface{}          `json:"entitlements"`
	User                         DiscordUser            `json:"user"`
	Type                         int                    `json:"type"`
	Version                      int                    `json:"version"`
}

type DiscordUser struct {
	AvatarDecorationData interface{} `json:"avatar_decoration_data"`
	Clan                 interface{} `json:"clan"`
	Avatar               string      `json:"avatar"`
	Discriminator        string      `json:"discriminator"`
	GlobalName           string      `json:"global_name"`
	ID                   string      `json:"id"`
	Username             string      `json:"username"`
	PublicFlags          int         `json:"public_flags"`
	Bot                  bool        `json:"bot"`
	System               bool        `json:"system"`
}

type DiscordService struct {
	AppPublicKey string
	AppID        string
	BotToken     string
}

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

type Command struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Type        int             `json:"type"`
	Options     []CommandOption `json:"options"`
}

type CommandOptionType int

const (
	StringType  CommandOptionType = 3
	IntegerType CommandOptionType = 4
)

type CommandOption struct {
	Type        CommandOptionType `json:"type"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Required    bool              `json:"required"`
}

var ErrFailedSettingUpGlobalCommand = fmt.Errorf("Failed setting up command")

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
	agent := fiber.Put(DISCORD_BASE_API_URL+"/applications/"+ds.AppID+"/commands").JSON(globalCommands).Set("Authorization", "Bot "+ds.BotToken)

	code, body, errs := agent.String()

	slog.Info("Error object", "error", errs)

	if errs != nil && len(errs) > 0 {
		slog.Warn("Error while trying to set global command", "error", errs)
		return errs[len(errs)-1]
	}

	if code != 200 {
		slog.Warn("Failed setting up global command", "response", body)
		return ErrFailedSettingUpGlobalCommand
	}

	slog.Info("Response from setting global command", "code", code, "response", body)

	return nil
}
