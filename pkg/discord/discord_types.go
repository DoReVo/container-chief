package discord

import "fmt"

const DISCORD_BASE_API_URL = "https://discord.com/api"

var ErrFailedSettingUpGlobalCommand = fmt.Errorf("failed setting up command")

const (
	StringType  CommandOptionType = 3
	IntegerType CommandOptionType = 4
)

const (
	// Pong ACK a Ping
	Pong InteractionCallbackType = 1
	// ChannelMessageWithSource respond to an interaction with a message
	ChannelMessageWithSource InteractionCallbackType = 4
	// DeferredChannelMessageWithSource ACK an interaction and edit a response later, the user sees a loading state
	DeferredChannelMessageWithSource InteractionCallbackType = 5
	// DeferredUpdateMessage for components, ACK an interaction and edit the original message later; the user does not see a loading state
	DeferredUpdateMessage InteractionCallbackType = 6
	// UpdateMessage for components, edit the message the component was attached to
	UpdateMessage InteractionCallbackType = 7
	// ApplicationCommandAutocompleteResult respond to an autocomplete interaction with suggested choices
	ApplicationCommandAutocompleteResult InteractionCallbackType = 8
	// Modal respond to an interaction with a popup modal
	Modal InteractionCallbackType = 9
	// PremiumRequired respond to an interaction with an upgrade button, only available for apps with monetization enabled
	PremiumRequired InteractionCallbackType = 10
)

type (
	InteractionCallbackType int
	InteractionWebhook      struct {
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
)

type (
	CommandOptionType int
	DiscordUser       struct {
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
)

type InteractionCallbackData struct {
	Content string `json:"content"`
}

type InteractionResponse struct {
	Data InteractionCallbackData `json:"data"`
	Type InteractionCallbackType `json:"type"`
}

type DiscordService struct {
	AppPublicKey string
	AppID        string
	BotToken     string
}

type Command struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Options     []CommandOption `json:"options"`
	Type        int             `json:"type"`
}

type CommandOption struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Type        CommandOptionType `json:"type"`
	Required    bool              `json:"required"`
}
