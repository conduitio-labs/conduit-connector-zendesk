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

package source

import (
	"context"
	"fmt"

	"github.com/conduitio-labs/conduit-connector-zendesk/config"
	"github.com/conduitio-labs/conduit-connector-zendesk/source/iterator"
	"github.com/conduitio-labs/conduit-connector-zendesk/source/position"

	sdk "github.com/conduitio/conduit-connector-sdk"
)

type Source struct {
	sdk.UnimplementedSource
	config   Config
	iterator Iterator
}

type Iterator interface {
	HasNext(ctx context.Context) bool
	Next(ctx context.Context) (sdk.Record, error)
	Stop()
}

// NewSource initialises a new source.
func NewSource() sdk.Source {
	return sdk.SourceWithMiddleware(&Source{}, sdk.DefaultSourceMiddleware()...)
}

// Parameters returns a map of named Parameters that describe how to configure the Source.
func (s *Source) Parameters() map[string]sdk.Parameter {
	return map[string]sdk.Parameter{
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
		KeyPollingPeriod: {
			Default:     "6s",
			Required:    false,
			Description: "Fetch interval for consecutive iterations",
		},
	}
}

// Configure parses zendesk config
func (s *Source) Configure(ctx context.Context, cfg map[string]string) error {
	zendeskConfig, err := Parse(cfg)
	if err != nil {
		return err
	}
	s.config = zendeskConfig
	return nil
}

// Open prepare the plugin to start sending records from the given position
func (s *Source) Open(ctx context.Context, rp sdk.Position) error {
	ticketPos, err := position.ParsePosition(rp)
	if err != nil {
		return err
	}

	s.iterator, err = iterator.NewCDCIterator(
		ctx,
		s.config.UserName,
		s.config.APIToken,
		s.config.Domain,
		s.config.PollingPeriod,
		ticketPos,
	)
	if err != nil {
		return err
	}
	return nil
}

// Read gets the next object from the zendesk api
func (s *Source) Read(ctx context.Context) (sdk.Record, error) {
	if !s.iterator.HasNext(ctx) {
		return sdk.Record{}, sdk.ErrBackoffRetry
	}

	r, err := s.iterator.Next(ctx)
	if err != nil {
		return sdk.Record{}, err
	}
	return r, nil
}

func (s *Source) Teardown(ctx context.Context) error {
	sdk.Logger(ctx).Trace().Msg("shutting down zendesk client")
	if s.iterator != nil {
		s.iterator.Stop()
		s.iterator = nil
	}
	return nil
}

func (s *Source) Ack(ctx context.Context, pos sdk.Position) error {
	ticketPos, err := position.ParsePosition(pos)
	if err != nil {
		return fmt.Errorf("invalid position: %w", err)
	}
	sdk.Logger(ctx).Trace().
		Float64("id", ticketPos.ID).
		Time("update_time", ticketPos.LastModified).
		Msg("ack received")
	return nil
}
