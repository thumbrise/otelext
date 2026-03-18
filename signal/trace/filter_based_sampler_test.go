package trace_test

import (
	"context"
	"strings"
	"testing"

	"github.com/thumbrise/otelext/internal/mock"
	"github.com/thumbrise/otelext/signal"
	"github.com/thumbrise/otelext/signal/trace"
	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// TestFilterBasedTraceSampler tests the basic functionality
func TestFilterBasedTraceSampler(t *testing.T) {
	// Clear filters for clean mocks
	signal.ClearFilters()

	sampler := trace.NewFilterBasedSampler()

	params := sdktrace.SamplingParameters{
		ParentContext: context.Background(),
		Attributes: []attribute.KeyValue{
			attribute.String("mocks", "value"),
		},
	}

	result := sampler.ShouldSample(params)

	// With no filters, should record and sample
	if result.Decision != sdktrace.RecordAndSample {
		t.Errorf("Expected RecordAndSample, got %v", result.Decision)
	}

	if len(result.Attributes) != 1 {
		t.Errorf("Expected 1 attribute, got %d", len(result.Attributes))
	}
}

// TestFilterBasedTraceSamplerWithFilters tests with registered filters
func TestFilterBasedTraceSamplerWithFilters(t *testing.T) {
	// Clear filters for clean mocks
	signal.ClearFilters()

	// Create mocks filters
	filter1 := mock.NewFilter("filter1", "Test filter 1", true)  // Should drop
	filter2 := mock.NewFilter("filter2", "Test filter 2", false) // Should not drop

	// Register filters
	signal.RegisterFilter(filter1)
	signal.RegisterFilter(filter2)

	sampler := trace.NewFilterBasedSampler()

	params := sdktrace.SamplingParameters{
		ParentContext: context.Background(),
		Attributes: []attribute.KeyValue{
			attribute.String("mocks", "value"),
		},
	}

	result := sampler.ShouldSample(params)

	// Since filter2 returns false (should not drop), the sampler should drop the trace
	if result.Decision != sdktrace.Drop {
		t.Errorf("Expected Drop, got %v", result.Decision)
	}

	// When dropping, attributes should be nil
	if result.Attributes != nil {
		t.Errorf("Expected nil attributes when dropping, got %v", result.Attributes)
	}
}

// TestFilterBasedTraceSamplerDescription tests the FDescription method
func TestFilterBasedTraceSamplerDescription(t *testing.T) {
	// Clear filters for clean mocks
	signal.ClearFilters()

	sampler := trace.NewFilterBasedSampler()

	// Test with no filters
	desc := sampler.Description()
	if desc == "" {
		t.Error("Expected non-empty description")
	}

	if !strings.Contains(desc, "Iterates on registered attributes based filters") {
		t.Error("Expected description to contain base text")
	}

	// Add filters and mocks description
	filter1 := mock.NewFilter("filter1", "First filter", true)
	filter2 := mock.NewFilter("filter2", "Second filter", false)

	signal.RegisterFilter(filter1)
	signal.RegisterFilter(filter2)

	desc = sampler.Description()
	if !strings.Contains(desc, "First filter") {
		t.Error("Expected description to contain filter descriptions")
	}

	if !strings.Contains(desc, "Second filter") {
		t.Error("Expected description to contain filter descriptions")
	}
}

// TestFilterBasedTraceSamplerWithEmptyAttributes tests with empty attributes
func TestFilterBasedTraceSamplerWithEmptyAttributes(t *testing.T) {
	// Clear filters for clean mocks
	signal.ClearFilters()

	// Create a filter that drops when no attributes
	filter := mock.NewFilter("empty-filter", "Empty attributes filter", true) // Filter says should drop

	signal.RegisterFilter(filter)

	sampler := trace.NewFilterBasedSampler()

	params := sdktrace.SamplingParameters{
		ParentContext: context.Background(),
		Attributes:    []attribute.KeyValue{}, // Empty attributes
	}

	result := sampler.ShouldSample(params)

	// Filter returns true (should drop), but sampler logic is: if ANY filter returns false -> drop
	// Since our filter returns true, sampler should NOT drop (RecordAndSample)
	if result.Decision != sdktrace.RecordAndSample {
		t.Errorf("Expected RecordAndSample (filter says drop, so sampler doesn't drop), got %v", result.Decision)
	}
}

// TestFilterBasedTraceSamplerWithMultipleFilters tests behavior with multiple filters
func TestFilterBasedTraceSamplerWithMultipleFilters(t *testing.T) {
	// Clear filters for clean mocks
	signal.ClearFilters()

	// Create multiple filters with different behaviors
	filters := []*mock.Filter{
		mock.NewFilter("f1", "Filter 1", true),
		mock.NewFilter("f2", "Filter 2", false), // This will cause drop
		mock.NewFilter("f3", "Filter 3", true),
	}

	for _, f := range filters {
		signal.RegisterFilter(f)
	}

	sampler := trace.NewFilterBasedSampler()

	params := sdktrace.SamplingParameters{
		ParentContext: context.Background(),
		Attributes: []attribute.KeyValue{
			attribute.String("mocks", "value"),
		},
	}

	result := sampler.ShouldSample(params)

	// Since filter2 returns false, the trace should be dropped
	if result.Decision != sdktrace.Drop {
		t.Errorf("Expected Drop due to filter2, got %v", result.Decision)
	}
}
