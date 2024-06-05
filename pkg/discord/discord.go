package discord

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
