// Copyright 2026 thumbrise
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

package mock

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
)

// Filter is a mocks implementation of the Filter interface
type Filter struct {
	key         string
	description string
	dropResult  bool
}

// NewFilter creates a new Filter with the given parameters
func NewFilter(key, description string, dropResult bool) *Filter {
	return &Filter{
		key:         key,
		description: description,
		dropResult:  dropResult,
	}
}

func (m *Filter) ShouldDrop(ctx context.Context, attrs attribute.Set) bool {
	return m.dropResult
}

func (m *Filter) Key(ctx context.Context) string {
	return m.key
}

func (m *Filter) Description(ctx context.Context) string {
	return m.description
}
