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

package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/conduitio/conduit-commons/opencdc"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// Writer implements the destination.Writer interface for SQLite
type Writer struct {
	db        *sql.DB
	tableName string
}

// NewWriter creates a new SQLite writer
func NewWriter(dbPath string, tableName string) (*Writer, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open SQLite database: %w", err)
	}

	// Create table if it doesn't exist
	createTableSQL := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id TEXT PRIMARY KEY,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			position TEXT,
			operation TEXT,
			metadata TEXT,
			key TEXT,
			payload_before TEXT,
			payload_after TEXT
		)
	`, tableName)

	if _, err := db.ExecContext(context.Background(), createTableSQL); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return &Writer{
		db:        db,
		tableName: tableName,
	}, nil
}

// Write implements the Writer interface
func (w *Writer) Write(ctx context.Context, records []opencdc.Record) error {
	tx, err := w.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	// #nosec G201 - tableName is controlled internally and not user input
	insertSQL := fmt.Sprintf(`
		INSERT OR REPLACE INTO %s (
			id, updated_at, position, operation, metadata, key, payload_before, payload_after
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, w.tableName)

	stmt, err := tx.PrepareContext(ctx, insertSQL)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, record := range records {
		// Generate a unique ID if not present
		id := string(record.Key.Bytes())
		if id == "" {
			id = fmt.Sprintf("%d_%s", time.Now().UnixNano(), record.Position)
		}

		// Convert metadata to JSON
		metadataJSON, err := json.Marshal(record.Metadata)
		if err != nil {
			return fmt.Errorf("failed to marshal metadata: %w", err)
		}

		// Convert position to string
		positionStr := string(record.Position)

		// Convert payloads to JSON strings
		var payloadBeforeStr, payloadAfterStr string
		if record.Payload.Before != nil {
			// If it's already RawData, use it directly
			if rawData, ok := record.Payload.Before.(opencdc.RawData); ok {
				payloadBeforeStr = string(rawData)
			} else {
				payloadBeforeJSON, err := json.Marshal(record.Payload.Before)
				if err != nil {
					return fmt.Errorf("failed to marshal payload.before: %w", err)
				}
				payloadBeforeStr = string(payloadBeforeJSON)
			}
		}

		if record.Payload.After != nil {
			// If it's already RawData, use it directly
			if rawData, ok := record.Payload.After.(opencdc.RawData); ok {
				payloadAfterStr = string(rawData)
			} else {
				payloadAfterJSON, err := json.Marshal(record.Payload.After)
				if err != nil {
					return fmt.Errorf("failed to marshal payload.after: %w", err)
				}
				payloadAfterStr = string(payloadAfterJSON)
			}
		}

		_, err = stmt.ExecContext(ctx,
			id,
			time.Now(),
			positionStr,
			record.Operation.String(),
			string(metadataJSON),
			string(record.Key.Bytes()),
			payloadBeforeStr,
			payloadAfterStr,
		)
		if err != nil {
			return fmt.Errorf("failed to insert record: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Close closes the database connection
func (w *Writer) Close() {
	if w.db != nil {
		w.db.Close()
	}
}
