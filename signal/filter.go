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
	mu.Lock()
	defer mu.Unlock()

	// Check for key collision with existing filters
	newKey := filter.Key(context.Background())
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

	return filters
}

func ClearFilters() {
	mu.Lock()
	defer mu.Unlock()

	filters = make([]Filter, 0)
	filterKeys = make(map[string]bool)
}
