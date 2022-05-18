/*
Copyright © 2022 Meroxa, Inc.

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
package config

import (
	"testing"

	"github.com/conduitio/conduit-connector-zendesk/config"
	"github.com/stretchr/testify/assert"
)

func TestParse_Destination(t *testing.T) {
	tests := []struct {
		name    string
		config  map[string]string
		want    Config
		isError bool
	}{
		{
			name: "Login with configured buffer size",
			config: map[string]string{
				KeyBufferSize:      "10",
				config.KeyDomain:   "testlab",
				config.KeyUserName: "test@testlab.com",
				config.KeyAPIToken: "gkdsaj)({jgo43646435#$!ga",
			},
			want: Config{
				BufferSize: 10,
				Config: config.Config{
					Domain:   "testlab",
					UserName: "test@testlab.com",
					APIToken: "gkdsaj)({jgo43646435#$!ga",
				},
			},
			isError: false,
		},
		{
			name: "Login without buffer size value",
			config: map[string]string{
				KeyBufferSize:      "",
				config.KeyDomain:   "testlab",
				config.KeyUserName: "test@testlab.com",
				config.KeyAPIToken: "gkdsaj)({jgo43646435#$!ga",
			},
			want: Config{
				BufferSize: 100,
				Config: config.Config{
					Domain:   "testlab",
					UserName: "test@testlab.com",
					APIToken: "gkdsaj)({jgo43646435#$!ga",
				},
			},
			isError: false,
		},
		{
			name: "Login without buffer size",
			config: map[string]string{
				config.KeyDomain:   "testlab",
				config.KeyUserName: "test@testlab.com",
				config.KeyAPIToken: "gkdsaj)({jgo43646435#$!ga",
			},
			want: Config{
				BufferSize: 100,
				Config: config.Config{
					Domain:   "testlab",
					UserName: "test@testlab.com",
					APIToken: "gkdsaj)({jgo43646435#$!ga",
				},
			},
			isError: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := Parse(tt.config)
			if tt.isError {
				assert.NotNil(t, err)
			} else {
				assert.NotNil(t, res)
				assert.Equal(t, res, tt.want)
			}
		})
	}
}
