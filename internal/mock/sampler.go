package mock

import (
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// Sampler is a mocks implementation of the sdktrace.Sampler interface
type Sampler struct {
	decision    sdktrace.SamplingDecision
	description string
}

// NewSampler creates a new Sampler with the given parameters
func NewSampler(decision sdktrace.SamplingDecision, description string) *Sampler {
	return &Sampler{
		decision:    decision,
		description: description,
	}
}

func (m *Sampler) ShouldSample(parameters sdktrace.SamplingParameters) sdktrace.SamplingResult {
	return sdktrace.SamplingResult{
		Decision:   m.decision,
		Attributes: nil,
		Tracestate: oteltrace.TraceState{},
	}
}

func (m *Sampler) Description() string {
	return m.description
}
