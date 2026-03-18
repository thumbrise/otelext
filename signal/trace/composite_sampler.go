package trace

import (
	"log/slog"
	"strings"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

type CompositeSampler struct {
	Samplers []sdktrace.Sampler
}

func NewCompositeSampler(samplers ...sdktrace.Sampler) *CompositeSampler {
	if len(samplers) == 0 {
		slog.Warn("no samplers passed in composite sampler, so always drop")
	}

	return &CompositeSampler{Samplers: samplers}
}

func (r CompositeSampler) ShouldSample(parameters sdktrace.SamplingParameters) sdktrace.SamplingResult {
	if len(r.Samplers) == 0 {
		return sdktrace.SamplingResult{
			Decision:   sdktrace.Drop,
			Attributes: nil,
			Tracestate: trace.TraceState{},
		}
	}

	var res sdktrace.SamplingResult
	for _, sampler := range r.Samplers {
		res = sampler.ShouldSample(parameters)
		if res.Decision == sdktrace.Drop {
			return res
		}
	}

	return res
}

func (r CompositeSampler) Description() string {
	if len(r.Samplers) == 0 {
		return "no samplers passed in composite sampler"
	}

	descriptions := make([]string, 0, len(r.Samplers))
	for _, sampler := range r.Samplers {
		descriptions = append(descriptions, sampler.Description())
	}

	descriptionsStr := strings.Join(descriptions, "\n")

	return "Decorates chain of samplers: " + descriptionsStr
}
