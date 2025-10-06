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
	"context"
	"errors"
	"testing"

	"github.com/conduitio-labs/conduit-connector-zendesk/destination/mocks"
	"github.com/conduitio/conduit-commons/opencdc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestMultiWriter_Write(t *testing.T) {
	ctx := context.Background()
	records := []opencdc.Record{
		{Key: opencdc.RawData("key1")},
		{Key: opencdc.RawData("key2")},
	}

	t.Run("successful write to all writers", func(t *testing.T) {
		writer1 := &mocks.Writer{}
		writer2 := &mocks.Writer{}

		writer1.On("Write", ctx, records).Return(nil)
		writer2.On("Write", ctx, records).Return(nil)

		multiWriter := NewMultiWriter(writer1, writer2)
		err := multiWriter.Write(ctx, records)

		require.NoError(t, err)
		writer1.AssertExpectations(t)
		writer2.AssertExpectations(t)
	})

	t.Run("error from one writer", func(t *testing.T) {
		writer1 := &mocks.Writer{}
		writer2 := &mocks.Writer{}

		writer1.On("Write", ctx, records).Return(nil)
		writer2.On("Write", ctx, records).Return(errors.New("write failed"))

		multiWriter := NewMultiWriter(writer1, writer2)
		err := multiWriter.Write(ctx, records)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "writer 1 failed")
		assert.Contains(t, err.Error(), "write failed")
		writer1.AssertExpectations(t)
		writer2.AssertExpectations(t)
	})

	t.Run("errors from multiple writers", func(t *testing.T) {
		writer1 := &mocks.Writer{}
		writer2 := &mocks.Writer{}
		writer3 := &mocks.Writer{}

		writer1.On("Write", ctx, records).Return(errors.New("error1"))
		writer2.On("Write", ctx, records).Return(nil)
		writer3.On("Write", ctx, records).Return(errors.New("error3"))

		multiWriter := NewMultiWriter(writer1, writer2, writer3)
		err := multiWriter.Write(ctx, records)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "writer 0 failed")
		assert.Contains(t, err.Error(), "writer 2 failed")
		assert.Contains(t, err.Error(), "error1")
		assert.Contains(t, err.Error(), "error3")
		writer1.AssertExpectations(t)
		writer2.AssertExpectations(t)
		writer3.AssertExpectations(t)
	})

	t.Run("single writer", func(t *testing.T) {
		writer := &mocks.Writer{}
		writer.On("Write", ctx, records).Return(nil)

		multiWriter := NewMultiWriter(writer)
		err := multiWriter.Write(ctx, records)

		require.NoError(t, err)
		writer.AssertExpectations(t)
	})

	t.Run("no writers", func(t *testing.T) {
		multiWriter := NewMultiWriter()
		err := multiWriter.Write(ctx, records)

		require.NoError(t, err)
	})
}

func TestMultiWriter_Close(t *testing.T) {
	t.Run("close all writers", func(t *testing.T) {
		writer1 := &mocks.Writer{}
		writer2 := &mocks.Writer{}
		writer3 := &mocks.Writer{}

		writer1.On("Close").Once()
		writer2.On("Close").Once()
		writer3.On("Close").Once()

		multiWriter := NewMultiWriter(writer1, writer2, writer3)
		multiWriter.Close()

		writer1.AssertExpectations(t)
		writer2.AssertExpectations(t)
		writer3.AssertExpectations(t)
	})

	t.Run("close with no writers", func(_ *testing.T) {
		multiWriter := NewMultiWriter()
		// Should not panic
		multiWriter.Close()
	})

	t.Run("close with mock that panics", func(t *testing.T) {
		writer1 := &mocks.Writer{}
		writer2 := &mocks.Writer{}

		writer1.On("Close").Once()
		writer2.On("Close").Run(func(_ mock.Arguments) {
			// Writer 2 close doesn't panic but does nothing
		}).Once()

		multiWriter := NewMultiWriter(writer1, writer2)
		multiWriter.Close()

		writer1.AssertExpectations(t)
		writer2.AssertExpectations(t)
	})
}
