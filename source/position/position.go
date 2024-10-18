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

package position

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/conduitio/conduit-commons/opencdc"
)

// Mode defines an iterator mode.
type Mode string

const (
	ModeSnapshot = "snapshot"
	ModeCDC      = "cdc"
)

type TicketPosition struct {
	// Mode represents current iterator mode.
	Mode         Mode      `json:"mode"`
	LastModified time.Time `json:"last_modified_time"`
	ID           float64   `json:"id"` // two tickets can have the same update time, id is to keep the position unique across tickets
}

// ToRecordPosition will extract the after_url from the ticket result json
func (pos *TicketPosition) ToRecordPosition() (opencdc.Position, error) {
	res, err := json.Marshal(pos)
	if err != nil {
		return opencdc.Position{}, fmt.Errorf("error in parsing the position %w", err)
	}

	return res, nil
}

// ParsePosition will unmarshal the TicketPosition used to record the next position
func ParsePosition(p opencdc.Position) (*TicketPosition, error) {
	var err error

	if p == nil {
		return nil, nil
	}

	var tp TicketPosition
	// parse the next position to opencdc.Record
	err = json.Unmarshal(p, &tp)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse the after_cursor position: %w", err)
	}

	return &tp, err
}
