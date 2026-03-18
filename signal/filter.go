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

package signal

import (
	"context"
	"log"
	"sync"

	"go.opentelemetry.io/otel/attribute"
)

type Filter interface {
	ShouldDrop(ctx context.Context, attrs attribute.Set) bool
	Key(ctx context.Context) string
	Description(ctx context.Context) string
}

var (
	mu         sync.RWMutex
	filters    []Filter
	filterKeys = make(map[string]bool)
)

func RegisterFilter(filter Filter) {
	// Evaluate key before acquiring lock to avoid deadlock if Key() calls back into the registry
	newKey := filter.Key(context.Background())

	mu.Lock()
	defer mu.Unlock()

	// Check for key collision with existing filters
	if filterKeys[newKey] {
		log.Printf("WARNING: Filter key collision detected: '%s'. Multiple filters with the same key may cause unexpected behavior during debugging", newKey)
	} else {
		filterKeys[newKey] = true
	}

	filters = append(filters, filter)
}

func RegisteredFilters() []Filter {
	mu.RLock()
	defer mu.RUnlock()

	result := make([]Filter, len(filters))
	copy(result, filters)

	return result
}

func ClearFilters() {
	mu.Lock()
	defer mu.Unlock()

	filters = make([]Filter, 0)
	filterKeys = make(map[string]bool)
}
