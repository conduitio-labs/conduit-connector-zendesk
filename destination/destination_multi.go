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
	"github.com/conduitio-labs/conduit-connector-zendesk/sqlite"
	"github.com/conduitio-labs/conduit-connector-zendesk/zendesk"
	cconfig "github.com/conduitio/conduit-commons/config"
	"github.com/conduitio/conduit-commons/opencdc"
	sdk "github.com/conduitio/conduit-connector-sdk"
)

type MultiDestination struct {
	sdk.UnimplementedDestination
	cfg    MultiConfig
	writer Writer
}

// NewMultiDestination initialises a new multi-destination
func NewMultiDestination() sdk.Destination {
	return sdk.DestinationWithMiddleware(&MultiDestination{}, sdk.DefaultDestinationMiddleware()...)
}

// Parameters returns a map of named Parameters that describe how to configure the destination
func (d *MultiDestination) Parameters() cconfig.Parameters {
	params := map[string]cconfig.Parameter{
		// Destination type selection
		KeyDestinationType: {
			Default:     string(DestinationTypeZendesk),
			Description: "Type of destination: 'zendesk', 'sqlite', or 'multi' for both",
		},

		// Zendesk parameters
		config.KeyDomain: {
			Default:     "",
			Description: "Zendesk domain (organization name)",
			Validations: []cconfig.Validation{},
		},
		config.KeyUserName: {
			Default:     "",
			Description: "Zendesk username (email)",
			Validations: []cconfig.Validation{},
		},
		config.KeyAPIToken: {
			Default:     "",
			Description: "Zendesk API token",
			Validations: []cconfig.Validation{},
		},
		KeyMaxRetries: {
			Default:     "3",
			Description: "Max API retries for Zendesk rate limiting",
		},
		KeyZendeskEnabled: {
			Default:     "",
			Description: "Enable Zendesk destination (only used when destination.type=multi)",
		},

		// SQLite parameters
		KeySQLiteEnabled: {
			Default:     "",
			Description: "Enable SQLite destination (only used when destination.type=multi)",
		},
		KeySQLiteDBPath: {
			Default:     "",
			Description: "Path to SQLite database file",
			Validations: []cconfig.Validation{},
		},
		KeySQLiteTableName: {
			Default:     defaultSQLiteTableName,
			Description: "Name of the table to store records in SQLite",
		},
	}

	return params
}

// Configure parses and initializes the config
func (d *MultiDestination) Configure(_ context.Context, cfg cconfig.Config) error {
	configuration, err := ParseMulti(cfg)
	if err != nil {
		return err
	}

	d.cfg = configuration
	return nil
}

// Open initializes the appropriate writer(s) based on configuration
func (d *MultiDestination) Open(_ context.Context) error {
	var writers []Writer

	// Initialize Zendesk writer if enabled
	if d.cfg.ZendeskEnabled {
		zendeskWriter := zendesk.NewBulkImporter(
			d.cfg.ZendeskConfig.UserName,
			d.cfg.ZendeskConfig.APIToken,
			d.cfg.ZendeskConfig.Domain,
			d.cfg.ZendeskConfig.MaxRetries,
		)
		writers = append(writers, zendeskWriter)
	}

	// Initialize SQLite writer if enabled
	if d.cfg.SQLiteEnabled {
		sqliteWriter, err := sqlite.NewWriter(d.cfg.SQLiteDBPath, d.cfg.SQLiteTableName)
		if err != nil {
			// Clean up any initialized writers
			for _, w := range writers {
				w.Close()
			}
			return err
		}
		writers = append(writers, sqliteWriter)
	}

	// Use single writer or multi-writer based on configuration
	if len(writers) == 1 {
		d.writer = writers[0]
	} else {
		d.writer = NewMultiWriter(writers...)
	}

	return nil
}

// Write writes records into the configured destination(s)
func (d *MultiDestination) Write(ctx context.Context, records []opencdc.Record) (int, error) {
	err := d.writer.Write(ctx, records)
	if err != nil {
		return 0, err
	}

	return len(records), nil
}

// Teardown closes any connections
func (d *MultiDestination) Teardown(_ context.Context) error {
	if d.writer != nil {
		d.writer.Close()
	}

	return nil
}
