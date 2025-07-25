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

package destination

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseMulti(t *testing.T) {
	t.Run("zendesk only destination", func(t *testing.T) {
		cfg := map[string]string{
			"destination.type": "zendesk",
			"zendesk.domain":   "test",
			"zendesk.userName": "user@example.com",
			"zendesk.apiToken": "token123",
		}

		config, err := ParseMulti(cfg)
		require.NoError(t, err)

		assert.Equal(t, DestinationTypeZendesk, config.DestinationType)
		assert.True(t, config.ZendeskEnabled)
		assert.False(t, config.SQLiteEnabled)
		assert.Equal(t, "test", config.ZendeskConfig.Domain)
	})

	t.Run("sqlite only destination", func(t *testing.T) {
		cfg := map[string]string{
			"destination.type":             "sqlite",
			"destination.sqlite.dbPath":    "/tmp/test.db",
			"destination.sqlite.tableName": "custom_table",
		}

		config, err := ParseMulti(cfg)
		require.NoError(t, err)

		assert.Equal(t, DestinationTypeSQLite, config.DestinationType)
		assert.False(t, config.ZendeskEnabled)
		assert.True(t, config.SQLiteEnabled)
		assert.Equal(t, "/tmp/test.db", config.SQLiteDBPath)
		assert.Equal(t, "custom_table", config.SQLiteTableName)
	})

	t.Run("multi destination with both enabled", func(t *testing.T) {
		cfg := map[string]string{
			"destination.type":            "multi",
			"destination.zendesk.enabled": "true",
			"destination.sqlite.enabled":  "true",
			"zendesk.domain":              "test",
			"zendesk.userName":            "user@example.com",
			"zendesk.apiToken":            "token123",
			"destination.sqlite.dbPath":   "/tmp/test.db",
		}

		config, err := ParseMulti(cfg)
		require.NoError(t, err)

		assert.Equal(t, DestinationTypeMulti, config.DestinationType)
		assert.True(t, config.ZendeskEnabled)
		assert.True(t, config.SQLiteEnabled)
		assert.Equal(t, "/tmp/test.db", config.SQLiteDBPath)
		assert.Equal(t, "records", config.SQLiteTableName) // default
	})

	t.Run("multi destination with zendesk disabled", func(t *testing.T) {
		cfg := map[string]string{
			"destination.type":            "multi",
			"destination.zendesk.enabled": "false",
			"destination.sqlite.enabled":  "true",
			"destination.sqlite.dbPath":   "/tmp/test.db",
		}

		config, err := ParseMulti(cfg)
		require.NoError(t, err)

		assert.Equal(t, DestinationTypeMulti, config.DestinationType)
		assert.False(t, config.ZendeskEnabled)
		assert.True(t, config.SQLiteEnabled)
	})

	t.Run("default to zendesk when no type specified", func(t *testing.T) {
		cfg := map[string]string{
			"zendesk.domain":   "test",
			"zendesk.userName": "user@example.com",
			"zendesk.apiToken": "token123",
		}

		config, err := ParseMulti(cfg)
		require.NoError(t, err)

		assert.Equal(t, DestinationTypeZendesk, config.DestinationType)
		assert.True(t, config.ZendeskEnabled)
		assert.False(t, config.SQLiteEnabled)
	})

	t.Run("error when no destinations enabled", func(t *testing.T) {
		cfg := map[string]string{
			"destination.type":            "multi",
			"destination.zendesk.enabled": "false",
			"destination.sqlite.enabled":  "false",
		}

		_, err := ParseMulti(cfg)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "at least one destination must be enabled")
	})

	t.Run("error when sqlite enabled without db path", func(t *testing.T) {
		cfg := map[string]string{
			"destination.type": "sqlite",
		}

		_, err := ParseMulti(cfg)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "destination.sqlite.dbPath is required")
	})

	t.Run("error with invalid destination type", func(t *testing.T) {
		cfg := map[string]string{
			"destination.type": "invalid",
		}

		_, err := ParseMulti(cfg)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid destination type")
	})

	t.Run("error with invalid boolean for enabled flags", func(t *testing.T) {
		cfg := map[string]string{
			"destination.type":            "multi",
			"destination.zendesk.enabled": "not-a-bool",
		}

		_, err := ParseMulti(cfg)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid value for destination.zendesk.enabled")
	})
}
