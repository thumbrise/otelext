package signal_test

import (
	"bytes"
	"context"
	"log"
	"strings"
	"testing"

	"github.com/thumbrise/otelext/internal/mock"
	"github.com/thumbrise/otelext/signal"
	"go.opentelemetry.io/otel/attribute"
)

// TestRegisterAndGet tests basic registration and retrieval
func TestRegisterAndGet(t *testing.T) {
	signal.ClearFilters()

	filter1 := mock.NewFilter("filter1", "First filter", false)

	filter2 := mock.NewFilter("filter2", "Second filter", true)

	if len(signal.RegisteredFilters()) != 0 {
		t.Errorf("Expected 0 registered filters, got %d", len(signal.RegisteredFilters()))
	}

	signal.RegisterFilter(filter1)
	signal.RegisterFilter(filter2)

	registered := signal.RegisteredFilters()
	if len(registered) != 2 {
		t.Errorf("Expected 2 registered filters, got %d", len(registered))
	}

	if registered[0] != signal.Filter(filter1) || registered[1] != signal.Filter(filter2) {
		t.Errorf("Expected [filter1, filter2], got [%v, %v]", registered[0], registered[1])
	}
}

// TestFilterKeyCollision tests warning when filters have the same key
func TestFilterKeyCollision(t *testing.T) {
	signal.ClearFilters()

	// Capture log output
	var logOutput bytes.Buffer

	// Save original output
	originalOutput := log.Writer()

	log.SetOutput(&logOutput)
	defer log.SetOutput(originalOutput)

	filter1 := mock.NewFilter("same-key", "First filter with same key", false)

	filter2 := mock.NewFilter("same-key", "Second filter with same key", true)

	signal.RegisterFilter(filter1)

	// Register second filter (should trigger warning)
	signal.RegisterFilter(filter2)

	// Check that warning was logged
	logContent := logOutput.String()
	if len(logContent) == 0 {
		t.Error("Expected warning about key collision, but no log output")
	} else if !strings.Contains(logContent, "same-key") {
		t.Errorf("Expected warning to mention 'same-key', got: %s", logContent)
	}

	// Test that both filters are still registered (no filtering in registration)
	registered := signal.RegisteredFilters()
	if len(registered) != 2 {
		t.Errorf("Expected 2 registered filters despite collision, got %d", len(registered))
	}
}

// TestFilterMethods tests the Filter interface methods
func TestFilterMethods(t *testing.T) {
	signal.ClearFilters()

	ctx := context.Background()
	attrs := attribute.NewSet(attribute.String("test", "value"))

	filter := mock.NewFilter("test-key", "test description", true)

	if key := filter.Key(ctx); key != "test-key" {
		t.Errorf("Key() = %q, want %q", key, "test-key")
	}

	if desc := filter.Description(ctx); desc != "test description" {
		t.Errorf("Description() = %q, want %q", desc, "test description")
	}

	if drop := filter.ShouldDrop(ctx, attrs); drop != true {
		t.Errorf("ShouldDrop() = %v, want %v", drop, true)
	}
}

// TestEmptyFilters tests behavior with no registered filters
func TestEmptyFilters(t *testing.T) {
	signal.ClearFilters()

	registered := signal.RegisteredFilters()
	if len(registered) != 0 {
		t.Errorf("Expected 0 registered filters when empty, got %d", len(registered))
	}
}
