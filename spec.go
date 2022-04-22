package zendesk

import (
	"github.com/conduitio/conduit-connector-zendesk/config"

	sdk "github.com/conduitio/conduit-connector-sdk"
)

func Specification() sdk.Specification {
	return sdk.Specification{
		Name:    "zendesk",
		Summary: "An zendesk source and destination plugin for Conduit, written in Go",
		Version: "v0.1.0",
		Author:  "Meroxa,Inc.",
		SourceParams: map[string]sdk.Parameter{
			config.ConfigKeyDomain: {
				Default:     "",
				Required:    true,
				Description: "A domain is referred as the organization name to which zendesk is registered",
			},
			config.ConfigKeyUserName: {
				Default:     "",
				Required:    true,
				Description: "Login to zendesk performed using username",
			},
			config.ConfigKeyAPIToken: {
				Default:     "",
				Required:    true,
				Description: "password to login",
			},
			config.ConfigKeyIterationInterval: {
				Default:     "2m",
				Required:    false,
				Description: "Fetch interval for consecutive iterations",
			},
		},
		DestinationParams: map[string]sdk.Parameter{},
	}
}
