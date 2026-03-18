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
	// Check if any filter wants to drop the trace
	for _, f := range signal.RegisteredFilters() {
		if f.ShouldDrop(parameters.ParentContext, attribute.NewSet(parameters.Attributes...)) {
			return sdktrace.SamplingResult{
				Decision:   sdktrace.Drop,
				Attributes: nil,
				Tracestate: trace.SpanContextFromContext(parameters.ParentContext).TraceState(),
			}
		}
	}

	// No filters dropped, record and sample
	psc := trace.SpanContextFromContext(parameters.ParentContext)
	return sdktrace.SamplingResult{
		Decision:   sdktrace.RecordAndSample,
		Attributes: parameters.Attributes,
		Tracestate: psc.TraceState(),
	}
}

func (s FilterBasedSampler) Description() string {
	descriptions := make([]string, 0, len(signal.RegisteredFilters()))
	for _, f := range signal.RegisteredFilters() {
		descriptions = append(descriptions, f.Description(context.Background()))
	}

	descriptionsStr := strings.Join(descriptions, ", ")

	return "Iterates on registered attributes based filters. Filters: " + descriptionsStr
}
