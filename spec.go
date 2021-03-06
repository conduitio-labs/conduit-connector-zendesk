/*
Copyright © 2022 Meroxa, Inc. & Gophers Lab Technologies Pvt. Ltd.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package zendesk

import (
	"github.com/conduitio-labs/conduit-connector-zendesk/config"
	"github.com/conduitio-labs/conduit-connector-zendesk/destination"
	"github.com/conduitio-labs/conduit-connector-zendesk/source"
	sdk "github.com/conduitio/conduit-connector-sdk"
)

func Specification() sdk.Specification {
	return sdk.Specification{
		Name:        "zendesk",
		Summary:     "Zendesk conduit plugin",
		Description: "A zendesk source and destination plugin for Conduit",
		Version:     "v0.1.0",
		Author:      "Gophers Lab Technologies Pvt Ltd",
		SourceParams: map[string]sdk.Parameter{
			config.KeyDomain: {
				Default:     "",
				Required:    true,
				Description: "A domain is referred as the organization name to which zendesk is registered",
			},
			config.KeyUserName: {
				Default:     "",
				Required:    true,
				Description: "Login to zendesk performed using username",
			},
			config.KeyAPIToken: {
				Default:     "",
				Required:    true,
				Description: "password to login",
			},
			source.KeyPollingPeriod: {
				Default:     "6s",
				Required:    false,
				Description: "Fetch interval for consecutive iterations",
			},
		},
		DestinationParams: map[string]sdk.Parameter{
			config.KeyDomain: {
				Default:     "",
				Required:    true,
				Description: "A domain is referred as the organization name to which zendesk is registered",
			},
			config.KeyUserName: {
				Default:     "",
				Required:    true,
				Description: "Login to zendesk performed using username",
			},
			config.KeyAPIToken: {
				Default:     "",
				Required:    true,
				Description: "password to login",
			},
			destination.KeyBufferSize: {
				Default:     "100",
				Required:    false,
				Description: "max tickets to be created in one API call",
			},
			destination.KeyMaxRetries: {
				Default:     "3",
				Required:    false,
				Description: "max API retries, before returning an error",
			},
		},
	}
}
