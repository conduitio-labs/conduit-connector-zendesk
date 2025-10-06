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
	"fmt"
	"strconv"
	"strings"
)

const (
	// Destination type keys
	KeyDestinationType = "destination.type"
	KeyZendeskEnabled  = "destination.zendesk.enabled"
	KeySQLiteEnabled   = "destination.sqlite.enabled"

	// SQLite specific configuration
	KeySQLiteDBPath    = "destination.sqlite.dbPath"
	KeySQLiteTableName = "destination.sqlite.tableName"

	// Default values
	defaultDestinationType = "zendesk"
	defaultSQLiteTableName = "records"
)

// Type represents the type of destination
type Type string

const (
	DestinationTypeZendesk Type = "zendesk"
	DestinationTypeSQLite  Type = "sqlite"
	DestinationTypeMulti   Type = "multi"
)

type MultiConfig struct {
	// Type of destination(s) to use
	DestinationType Type

	// Zendesk configuration
	ZendeskEnabled bool
	ZendeskConfig  Config

	// SQLite configuration
	SQLiteEnabled   bool
	SQLiteDBPath    string
	SQLiteTableName string
}

// ParseMulti parses and validates multi-destination configuration
func ParseMulti(cfg map[string]string) (MultiConfig, error) {
	multiConfig := MultiConfig{
		SQLiteTableName: defaultSQLiteTableName,
	}

	// Parse destination type
	destType := cfg[KeyDestinationType]
	if destType == "" {
		destType = defaultDestinationType
	}

	// Parse enabled flags
	zendeskEnabledStr := cfg[KeyZendeskEnabled]
	sqliteEnabledStr := cfg[KeySQLiteEnabled]

	// If destination type is specified, use it to set enabled flags
	switch strings.ToLower(destType) {
	case string(DestinationTypeZendesk):
		multiConfig.DestinationType = DestinationTypeZendesk
		multiConfig.ZendeskEnabled = true
		multiConfig.SQLiteEnabled = false
	case string(DestinationTypeSQLite):
		multiConfig.DestinationType = DestinationTypeSQLite
		multiConfig.ZendeskEnabled = false
		multiConfig.SQLiteEnabled = true
	case string(DestinationTypeMulti):
		multiConfig.DestinationType = DestinationTypeMulti
		// Parse individual enabled flags
		if zendeskEnabledStr != "" {
			enabled, err := strconv.ParseBool(zendeskEnabledStr)
			if err != nil {
				return MultiConfig{}, fmt.Errorf("invalid value for %s: %w", KeyZendeskEnabled, err)
			}
			multiConfig.ZendeskEnabled = enabled
		} else {
			multiConfig.ZendeskEnabled = true // Default to enabled for backward compatibility
		}

		if sqliteEnabledStr != "" {
			enabled, err := strconv.ParseBool(sqliteEnabledStr)
			if err != nil {
				return MultiConfig{}, fmt.Errorf("invalid value for %s: %w", KeySQLiteEnabled, err)
			}
			multiConfig.SQLiteEnabled = enabled
		}
	default:
		return MultiConfig{}, fmt.Errorf("invalid destination type: %s", destType)
	}

	// If neither destination is enabled, return error
	if !multiConfig.ZendeskEnabled && !multiConfig.SQLiteEnabled {
		return MultiConfig{}, fmt.Errorf("at least one destination must be enabled")
	}

	// Parse Zendesk config if enabled
	if multiConfig.ZendeskEnabled {
		zendeskConfig, err := Parse(cfg)
		if err != nil {
			return MultiConfig{}, fmt.Errorf("failed to parse Zendesk config: %w", err)
		}
		multiConfig.ZendeskConfig = zendeskConfig
	}

	// Parse SQLite config if enabled
	if multiConfig.SQLiteEnabled {
		dbPath := cfg[KeySQLiteDBPath]
		if dbPath == "" {
			return MultiConfig{}, fmt.Errorf("%s is required when SQLite is enabled", KeySQLiteDBPath)
		}
		multiConfig.SQLiteDBPath = dbPath

		if tableName := cfg[KeySQLiteTableName]; tableName != "" {
			multiConfig.SQLiteTableName = tableName
		}
	}

	return multiConfig, nil
}
