package trace

import (
	"log/slog"
	"strings"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

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

func (r CompositeSampler) ShouldSample(parameters sdktrace.SamplingParameters) sdktrace.SamplingResult {
	if len(r.samplers) == 0 {
		return sdktrace.SamplingResult{
			Decision:   sdktrace.Drop,
			Attributes: nil,
			Tracestate: trace.TraceState{},
		}
	}

	var res sdktrace.SamplingResult
	for _, sampler := range r.samplers {
		res = sampler.ShouldSample(parameters)
		if res.Decision == sdktrace.Drop {
			return res
		}
	}

	return res
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
