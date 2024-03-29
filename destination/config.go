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

	"github.com/conduitio-labs/conduit-connector-zendesk/config"
)

const (
	// KeyMaxRetries is the number of times the writer needs to retry, in case of 429 error, before returning an error
	KeyMaxRetries = "maxRetries"

	defaultMaxRetries = "3"
)

type Config struct {
	config.Config
	MaxRetries uint64
}

// Parse parses and validates configuration.
func Parse(cfg map[string]string) (Config, error) {
	defaultConfig, err := config.Parse(cfg)
	if err != nil {
		return Config{}, err
	}

	maxRetriesString := cfg[KeyMaxRetries]
	if maxRetriesString == "" {
		maxRetriesString = defaultMaxRetries
	}

	maxRetries, err := strconv.ParseUint(maxRetriesString, 10, 64)
	if err != nil {
		return Config{}, fmt.Errorf(
			"%q config value should be a positive integer",
			KeyMaxRetries,
		)
	}

	destinationConfig := Config{
		Config:     defaultConfig,
		MaxRetries: maxRetries,
	}

	return destinationConfig, nil
}
