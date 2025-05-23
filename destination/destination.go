// Copyright © 2022 Meroxa, Inc. & Gophers Lab Technologies Pvt. Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:generate mockery --name=Writer

package destination

import (
	"context"

	"github.com/conduitio-labs/conduit-connector-zendesk/config"
	"github.com/conduitio-labs/conduit-connector-zendesk/zendesk"
	cconfig "github.com/conduitio/conduit-commons/config"
	"github.com/conduitio/conduit-commons/opencdc"
	sdk "github.com/conduitio/conduit-connector-sdk"
)

type Writer interface {
	Write(ctx context.Context, records []opencdc.Record) error
	Close()
}

type Destination struct {
	sdk.UnimplementedDestination
	cfg    Config // destination specific config for zendesk
	writer Writer // interface that implements to write tickets to zendesk
}

// NewDestination initialises a new Destination.
func NewDestination() sdk.Destination {
	return sdk.DestinationWithMiddleware(&Destination{}, sdk.DefaultDestinationMiddleware()...)
}

// Parameters returns a map of named Parameters that describe how to configure the Source.
func (d *Destination) Parameters() cconfig.Parameters {
	return map[string]cconfig.Parameter{
		config.KeyDomain: {
			Default:     "",
			Description: "A domain is referred as the organization name to which zendesk is registered",
			Validations: []cconfig.Validation{cconfig.ValidationRequired{}},
		},
		config.KeyUserName: {
			Default:     "",
			Description: "Login to zendesk performed using username",
			Validations: []cconfig.Validation{cconfig.ValidationRequired{}},
		},
		config.KeyAPIToken: {
			Default:     "",
			Description: "password to login",
			Validations: []cconfig.Validation{cconfig.ValidationRequired{}},
		},
		KeyMaxRetries: {
			Default:     "3",
			Description: "max API retries, before returning an error",
		},
	}
}

// Configure parses and initializes the config.
func (d *Destination) Configure(_ context.Context, cfg cconfig.Config) error {
	configuration, err := Parse(cfg)
	if err != nil {
		return err
	}

	d.cfg = configuration

	return nil
}

// Open initializes a http client.
func (d *Destination) Open(_ context.Context) error {
	d.writer = zendesk.NewBulkImporter(d.cfg.UserName, d.cfg.APIToken, d.cfg.Domain, d.cfg.MaxRetries)

	return nil
}

// Write writes records into a Destination.
func (d *Destination) Write(ctx context.Context, records []opencdc.Record) (int, error) {
	err := d.writer.Write(ctx, records)
	if err != nil {
		return 0, err
	}

	return len(records), nil
}

// Teardown closes any connections which were previously connected from previous requests.
func (d *Destination) Teardown(_ context.Context) error {
	if d.writer != nil {
		d.writer.Close()
	}

	return nil
}
