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
	"fmt"

	"github.com/conduitio/conduit-commons/opencdc"
)

// MultiWriter writes to multiple destinations
type MultiWriter struct {
	writers []Writer
}

// NewMultiWriter creates a new MultiWriter
func NewMultiWriter(writers ...Writer) *MultiWriter {
	return &MultiWriter{
		writers: writers,
	}
}

// Write writes records to all enabled destinations
func (m *MultiWriter) Write(ctx context.Context, records []opencdc.Record) error {
	var errs []error

	for i, writer := range m.writers {
		if err := writer.Write(ctx, records); err != nil {
			errs = append(errs, fmt.Errorf("writer %d failed: %w", i, err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("multi-writer errors: %v", errs)
	}

	return nil
}

// Close closes all writers
func (m *MultiWriter) Close() {
	for _, writer := range m.writers {
		writer.Close()
	}
}
