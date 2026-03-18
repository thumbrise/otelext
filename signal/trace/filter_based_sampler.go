package trace

import (
	"context"
	"strings"

	"github.com/thumbrise/otelext/signal"
	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

type FilterBasedSampler struct{}

func NewFilterBasedSampler() *FilterBasedSampler {
	return &FilterBasedSampler{}
}

func (s FilterBasedSampler) ShouldSample(parameters sdktrace.SamplingParameters) sdktrace.SamplingResult {
	psc := trace.SpanContextFromContext(parameters.ParentContext)
	result := sdktrace.SamplingResult{
		Decision:   sdktrace.RecordAndSample,
		Attributes: parameters.Attributes,
		Tracestate: psc.TraceState(),
	}

	for _, f := range signal.RegisteredFilters() {
		if !f.ShouldDrop(parameters.ParentContext, attribute.NewSet(parameters.Attributes...)) {
			result.Decision = sdktrace.Drop
			result.Attributes = nil
		}
	}

	return result
}

func (s FilterBasedSampler) Description() string {
	descriptions := make([]string, 0, len(signal.RegisteredFilters()))
	for _, f := range signal.RegisteredFilters() {
		descriptions = append(descriptions, f.Description(context.Background()))
	}

	descriptionsStr := strings.Join(descriptions, ", ")

	return "Iterates on registered attributes based filters. Filters: " + descriptionsStr
}
