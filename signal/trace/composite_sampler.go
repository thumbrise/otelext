package trace

import (
	"log/slog"
	"strings"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// CompositeSampler decorates samplers.
//
// Mixed decision resolving strategy is restrict-only.
//
// For example:
//
// If any sampler returns sdktrace.Drop, the
// decision is immediately Drop. If a sampler returns sdktrace.RecordOnly,
// that decision is recorded, and subsequent samplers cannot upgrade it to
// sdktrace.RecordAndSample. If a sampler returns sdktrace.RecordAndSample,
// it is only considered if no prior sampler decided sdktrace.RecordOnly.
// This ensures that restrictions are respected and the behavior is restrict-only
// where any sampler can restrict the decision, but not upgrade it beyond
// what a previous sampler allowed.
type CompositeSampler struct {
	samplers []sdktrace.Sampler
}

func NewCompositeSampler(samplers ...sdktrace.Sampler) *CompositeSampler {
	if len(samplers) == 0 {
		slog.Warn("no samplers passed in composite sampler, so always drop")
	}

	copied := make([]sdktrace.Sampler, len(samplers))
	copy(copied, samplers)

	return &CompositeSampler{samplers: copied}
}

// ShouldSample determines strictest decision in configured samplers.
//
// See CompositeSampler
func (r CompositeSampler) ShouldSample(parameters sdktrace.SamplingParameters) sdktrace.SamplingResult {
	if len(r.samplers) == 0 {
		return sdktrace.SamplingResult{
			Decision:   sdktrace.Drop,
			Attributes: nil,
			Tracestate: trace.TraceState{},
		}
	}

	// Relies on OTel SDK ordering: Drop(0) < RecordOnly(1) < RecordAndSample(2).
	// See TestSamplingDecisionOrdering.
	const strictest = sdktrace.Drop

	result := r.samplers[0].ShouldSample(parameters)

	for _, sampler := range r.samplers[1:] {
		if result.Decision == strictest {
			break
		}

		current := sampler.ShouldSample(parameters)
		if current.Decision < result.Decision {
			result = current
		}
	}

	return result
}

func (r CompositeSampler) Description() string {
	if len(r.samplers) == 0 {
		return "no samplers passed in composite sampler"
	}

	descriptions := make([]string, 0, len(r.samplers))
	for _, sampler := range r.samplers {
		descriptions = append(descriptions, sampler.Description())
	}

	descriptionsStr := strings.Join(descriptions, "\n")

	return "Decorates chain of samplers: " + descriptionsStr
}
