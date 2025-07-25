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
	"os"
	"testing"

	"github.com/conduitio/conduit-commons/opencdc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewWriter(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "test-*.db")
		require.NoError(t, err)
		defer os.Remove(tmpFile.Name())

		writer, err := NewWriter(tmpFile.Name(), "test_table")
		require.NoError(t, err)
		require.NotNil(t, writer)
		defer writer.Close()

		// Verify table was created
		var tableName string
		err = writer.db.QueryRowContext(context.Background(), "SELECT name FROM sqlite_master WHERE type='table' AND name='test_table'").Scan(&tableName)
		require.NoError(t, err)
		assert.Equal(t, "test_table", tableName)
	})

	t.Run("invalid database path", func(t *testing.T) {
		writer, err := NewWriter("/invalid/path/test.db", "test_table")
		assert.Error(t, err)
		assert.Nil(t, writer)
	})
}

func TestWriter_Write(t *testing.T) {
	ctx := context.Background()

	t.Run("successful write", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "test-*.db")
		require.NoError(t, err)
		defer os.Remove(tmpFile.Name())

		writer, err := NewWriter(tmpFile.Name(), "test_table")
		require.NoError(t, err)
		defer writer.Close()

		records := []opencdc.Record{
			{
				Position:  opencdc.Position("pos1"),
				Operation: opencdc.OperationCreate,
				Metadata:  opencdc.Metadata{"key": "value"},
				Key:       opencdc.RawData("record1"),
				Payload: opencdc.Change{
					After: opencdc.RawData(`{"id": 1, "name": "test"}`),
				},
			},
			{
				Position:  opencdc.Position("pos2"),
				Operation: opencdc.OperationUpdate,
				Metadata:  opencdc.Metadata{"key": "value2"},
				Key:       opencdc.RawData("record2"),
				Payload: opencdc.Change{
					Before: opencdc.RawData(`{"id": 2, "name": "old"}`),
					After:  opencdc.RawData(`{"id": 2, "name": "new"}`),
				},
			},
		}

		err = writer.Write(ctx, records)
		require.NoError(t, err)

		// Verify records were written
		var count int
		err = writer.db.QueryRowContext(context.Background(), "SELECT COUNT(*) FROM test_table").Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, 2, count)

		// Verify first record
		var id, key, operation string
		err = writer.db.QueryRowContext(context.Background(), "SELECT id, key, operation FROM test_table WHERE id = ?", "record1").Scan(&id, &key, &operation)
		require.NoError(t, err)
		assert.Equal(t, "record1", id)
		assert.Equal(t, "record1", key)
		assert.Equal(t, "create", operation)
	})

	t.Run("write with empty key", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "test-*.db")
		require.NoError(t, err)
		defer os.Remove(tmpFile.Name())

		writer, err := NewWriter(tmpFile.Name(), "test_table")
		require.NoError(t, err)
		defer writer.Close()

		records := []opencdc.Record{
			{
				Position:  opencdc.Position("pos1"),
				Operation: opencdc.OperationCreate,
				Key:       opencdc.RawData(""),
				Payload: opencdc.Change{
					After: opencdc.RawData(`{"id": 1}`),
				},
			},
		}

		err = writer.Write(ctx, records)
		require.NoError(t, err)

		// Verify record was written with generated ID
		var count int
		err = writer.db.QueryRowContext(context.Background(), "SELECT COUNT(*) FROM test_table").Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, 1, count)
	})

	t.Run("update existing record", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "test-*.db")
		require.NoError(t, err)
		defer os.Remove(tmpFile.Name())

		writer, err := NewWriter(tmpFile.Name(), "test_table")
		require.NoError(t, err)
		defer writer.Close()

		// Write initial record
		records := []opencdc.Record{
			{
				Position:  opencdc.Position("pos1"),
				Operation: opencdc.OperationCreate,
				Key:       opencdc.RawData("record1"),
				Payload: opencdc.Change{
					After: opencdc.RawData(`{"id": 1, "name": "initial"}`),
				},
			},
		}
		err = writer.Write(ctx, records)
		require.NoError(t, err)

		// Update the record
		records = []opencdc.Record{
			{
				Position:  opencdc.Position("pos2"),
				Operation: opencdc.OperationUpdate,
				Key:       opencdc.RawData("record1"),
				Payload: opencdc.Change{
					Before: opencdc.RawData(`{"id": 1, "name": "initial"}`),
					After:  opencdc.RawData(`{"id": 1, "name": "updated"}`),
				},
			},
		}
		err = writer.Write(ctx, records)
		require.NoError(t, err)

		// Verify only one record exists (updated)
		var count int
		err = writer.db.QueryRowContext(context.Background(), "SELECT COUNT(*) FROM test_table").Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, 1, count)

		// Verify the payload was updated
		var payloadAfter string
		err = writer.db.QueryRowContext(context.Background(), "SELECT payload_after FROM test_table WHERE id = ?", "record1").Scan(&payloadAfter)
		require.NoError(t, err)
		// Decode the base64 JSON to check content
		assert.Contains(t, payloadAfter, "updated")
	})
}

func TestWriter_Close(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-*.db")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	writer, err := NewWriter(tmpFile.Name(), "test_table")
	require.NoError(t, err)

	// Close should not error
	writer.Close()

	// Verify database is closed
	err = writer.db.PingContext(context.Background())
	// Check that the error indicates the database is closed
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "closed")

	// Multiple closes should not panic
	writer.Close()
}
