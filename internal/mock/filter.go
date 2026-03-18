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
