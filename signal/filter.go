package signal

import (
	"context"
	"log"

	"go.opentelemetry.io/otel/attribute"
)

type Filter interface {
	ShouldDrop(ctx context.Context, attrs attribute.Set) bool
	Key(ctx context.Context) string
	Description(ctx context.Context) string
}

var (
	filters    []Filter
	filterKeys = make(map[string]bool)
)

func RegisterFilter(filter Filter) {
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
	return filters
}

func ClearFilters() {
	filters = make([]Filter, 0)
}
