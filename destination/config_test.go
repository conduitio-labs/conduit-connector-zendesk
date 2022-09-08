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
	"testing"

	"github.com/conduitio-labs/conduit-connector-zendesk/config"
	"github.com/stretchr/testify/assert"
)

func TestParse_Destination(t *testing.T) {
	tests := []struct {
		name    string
		config  map[string]string
		want    Config
		isError bool
		err     error
	}{
		{
			name: "valid configuration",
			config: map[string]string{
				config.KeyDomain:   "testlab",
				config.KeyUserName: "test@testlab.com",
				config.KeyAPIToken: "gkdsaj)({jgo43646435#$!ga",
				KeyMaxRetries:      "5",
			},
			want: Config{
				MaxRetries: 5,
				Config: config.Config{
					Domain:   "testlab",
					UserName: "test@testlab.com",
					APIToken: "gkdsaj)({jgo43646435#$!ga",
				},
			},
			isError: false,
			err:     nil,
		},
		{
			name: "Login without max retries",
			config: map[string]string{
				config.KeyDomain:   "testlab",
				config.KeyUserName: "test@testlab.com",
				config.KeyAPIToken: "gkdsaj)({jgo43646435#$!ga",
			},
			want: Config{
				MaxRetries: 3,
				Config: config.Config{
					Domain:   "testlab",
					UserName: "test@testlab.com",
					APIToken: "gkdsaj)({jgo43646435#$!ga",
				},
			},
			isError: false,
			err:     nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := Parse(tt.config)
			if tt.isError {
				assert.NotNil(t, err)
				assert.EqualError(t, err, tt.err.Error())
			} else {
				assert.NotNil(t, res)
				assert.Equal(t, res, tt.want)
			}
		})
	}
}
