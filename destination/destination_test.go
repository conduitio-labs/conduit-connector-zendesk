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

package destination

import (
	"context"
	"errors"
	"testing"

	"github.com/conduitio-labs/conduit-connector-zendesk/config"
	"github.com/conduitio-labs/conduit-connector-zendesk/destination/mocks"
	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestConfigure(t *testing.T) {
	invalidCfg := map[string]string{
		config.KeyDomain:   "test.lab",
		config.KeyUserName: "",
		config.KeyAPIToken: "ajgrmrop&90002p$@7",
	}

	validConfig := map[string]string{
		config.KeyDomain:   "testlab",
		config.KeyUserName: "test",
		config.KeyAPIToken: "ajgrmrop&90002p$@7",
	}

	type field struct {
		cfg map[string]string
	}
	tests := []struct {
		name    string
		field   field
		want    config.Config
		isError bool
	}{
		{
			name: "valid config",
			field: field{
				cfg: validConfig,
			},
			isError: false,
		},
		{
			name: "invalid config",
			field: field{
				cfg: invalidCfg,
			},
			isError: true,
		},
	}
	var destination Destination
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := destination.Configure(context.Background(), tt.field.cfg)
			if tt.isError {
				assert.NotNil(t, err)
			}
		})
	}
}

func TestNewDestination(t *testing.T) {
	dest := NewDestination()
	assert.NotNil(t, dest)
}

func TestOpen(t *testing.T) {
	d := NewDestination()
	err := d.Open(context.Background())
	assert.Nil(t, err)
}

func TestWrite(t *testing.T) {
	tests := []struct {
		name   string
		record sdk.Record
		err    error
		dest   Destination
	}{
		{
			name: "write empty record",
			dest: Destination{
				writer: func() Writer {
					w := &mocks.Writer{}
					w.On("Write", mock.Anything, mock.Anything).Return(nil)

					return w
				}(),
			},
			record: sdk.Record{
				Key:     sdk.RawData(`dummy_key`),
				Payload: sdk.Change{After: sdk.RawData(``)},
			},
			err: nil,
		},
		{
			name: "valid case",
			record: sdk.Record{
				Payload: sdk.Change{After: sdk.RawData(`"dummy_data":"12345"`)},
			},
			dest: Destination{
				writer: func() Writer {
					w := &mocks.Writer{}
					w.On("Write", mock.Anything, mock.Anything).Return(nil)

					return w
				}(),
			},
		},
		{
			name: "write invalid case with flush error",
			record: sdk.Record{
				Payload: sdk.Change{After: sdk.RawData(`"dummy_data":"12345"`)},
			},
			dest: Destination{
				writer: func() Writer {
					w := &mocks.Writer{}
					w.On("Write", mock.Anything, mock.Anything).Return(errors.New("testing error"))

					return w
				}(),
			},
			err: errors.New("testing error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n, err := tt.dest.Write(context.Background(), []sdk.Record{tt.record})
			if tt.err != nil {
				assert.NotNil(t, err)
				assert.Equal(t, err, tt.err)
				assert.Equal(t, n, 0)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, n, 1)
			}
		})
	}
}

func TestTearDown(t *testing.T) {
	tests := []struct {
		name string
		dest Destination
		want error
	}{
		{
			name: "writer valid case for teardown",
			dest: Destination{
				writer: func() Writer {
					w := &mocks.Writer{}
					w.On("Close")

					return w
				}(),
			},
		},
		{
			name: "nil writer case",
			dest: Destination{},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.dest.Teardown(context.Background())
			if tt.want != nil {
				assert.NotNil(t, err)
				return
			}
			assert.Nil(t, err)
		})
	}
}
