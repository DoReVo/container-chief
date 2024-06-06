package discord

import (
	"crypto/ed25519"
	"encoding/hex"
	"os"
)

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
}

func NewDiscordService() *DiscordService {
	appPublicKey := os.Getenv("DISCORD_APP_PUBLIC_KEY")

	if appPublicKey == "" {
		panic("Discord app public key not defined in environment variable")
	}

	return &DiscordService{
		AppPublicKey: appPublicKey,
	}
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
