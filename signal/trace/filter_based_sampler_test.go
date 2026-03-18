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

	// SamplingResult.Attributes are additional attributes per OTel spec, sampler should not echo back input
	if result.Attributes != nil {
		t.Errorf("Expected nil additional attributes, got %v", result.Attributes)
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

	// Since filter1 returns true (should drop), the sampler should drop the trace
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

	// Filter returns true (should drop), so sampler should drop the trace
	if result.Decision != sdktrace.Drop {
		t.Errorf("Expected Drop (filter says drop), got %v", result.Decision)
	}
}

// TestFilterBasedTraceSamplerWithMultipleFilters tests behavior with multiple filters
func TestFilterBasedTraceSamplerWithMultipleFilters(t *testing.T) {
	// Clear filters for clean mocks
	signal.ClearFilters()

	// Create multiple filters with different behaviors
	filters := []*mock.Filter{
		mock.NewFilter("f1", "Filter 1", true),  // Should drop
		mock.NewFilter("f2", "Filter 2", false), // Should not drop
		mock.NewFilter("f3", "Filter 3", true),  // Should drop
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

	// Since filter1 and filter3 return true (should drop), the trace should be dropped
	if result.Decision != sdktrace.Drop {
		t.Errorf("Expected Drop due to filter1 and filter3, got %v", result.Decision)
	}
}

// TestFilterBasedTraceSamplerAllFiltersPass tests when all filters pass (don't drop)
func TestFilterBasedTraceSamplerAllFiltersPass(t *testing.T) {
	// Clear filters for clean mocks
	signal.ClearFilters()

	// Create filters that all pass (don't drop)
	filters := []*mock.Filter{
		mock.NewFilter("f1", "Filter 1", false), // Should not drop
		mock.NewFilter("f2", "Filter 2", false), // Should not drop
		mock.NewFilter("f3", "Filter 3", false), // Should not drop
	}

	for _, f := range filters {
		signal.RegisterFilter(f)
	}

	sampler := trace.NewFilterBasedSampler()

	params := sdktrace.SamplingParameters{
		ParentContext: context.Background(),
		Attributes: []attribute.KeyValue{
			attribute.String("test", "value"),
		},
	}

	result := sampler.ShouldSample(params)

	// All filters pass, so trace should be recorded and sampled
	if result.Decision != sdktrace.RecordAndSample {
		t.Errorf("Expected RecordAndSample when all filters pass, got %v", result.Decision)
	}

	// SamplingResult.Attributes are additional attributes per OTel spec, sampler should not echo back input
	if result.Attributes != nil {
		t.Errorf("Expected nil additional attributes, got %v", result.Attributes)
	}
}
