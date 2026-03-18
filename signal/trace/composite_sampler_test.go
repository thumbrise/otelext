package trace_test

import (
	"testing"

	"github.com/thumbrise/otelext/signal/trace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/thumbrise/otelext/internal/mock"
)

func TestCompositeSampler(t *testing.T) {
	tests := []struct {
		name     string
		samplers []sdktrace.Sampler
		params   sdktrace.SamplingParameters
		want     sdktrace.SamplingDecision
	}{
		{
			name:     "no samplers",
			samplers: []sdktrace.Sampler{},
			params:   sdktrace.SamplingParameters{},
			want:     sdktrace.Drop,
		},
		{
			name: "first sampler drops",
			samplers: []sdktrace.Sampler{
				mock.NewSampler(sdktrace.Drop, "drop sampler"),
				mock.NewSampler(sdktrace.RecordAndSample, "record sampler"),
			},
			params: sdktrace.SamplingParameters{},
			want:   sdktrace.Drop,
		},
		{
			name: "all samplers record and sample",
			samplers: []sdktrace.Sampler{
				mock.NewSampler(sdktrace.RecordAndSample, "record sampler 1"),
				mock.NewSampler(sdktrace.RecordAndSample, "record sampler 2"),
			},
			params: sdktrace.SamplingParameters{},
			want:   sdktrace.RecordAndSample,
		},
		{
			name: "mixed decisions",
			samplers: []sdktrace.Sampler{
				mock.NewSampler(sdktrace.RecordOnly, "record only"),
				mock.NewSampler(sdktrace.Drop, "drop"),
				mock.NewSampler(sdktrace.RecordAndSample, "record and sample"),
			},
			params: sdktrace.SamplingParameters{},
			want:   sdktrace.Drop,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sampler := trace.NewCompositeSampler(tt.samplers...)
			result := sampler.ShouldSample(tt.params)

			if result.Decision != tt.want {
				t.Errorf("CompositeSampler.ShouldSample() = %v, want %v", result.Decision, tt.want)
			}
		})
	}
}

func TestCompositeSamplerDescription(t *testing.T) {
	tests := []struct {
		name     string
		samplers []sdktrace.Sampler
		want     string
	}{
		{
			name:     "no samplers",
			samplers: []sdktrace.Sampler{},
			want:     "no samplers passed in composite sampler",
		},
		{
			name: "with samplers",
			samplers: []sdktrace.Sampler{
				mock.NewSampler(sdktrace.RecordAndSample, "sampler1"),
				mock.NewSampler(sdktrace.RecordAndSample, "sampler2"),
			},
			want: "Decorates chain of samplers: sampler1\nsampler2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sampler := trace.NewCompositeSampler(tt.samplers...)
			got := sampler.Description()

			if got != tt.want {
				t.Errorf("CompositeSampler.Description() = %q, want %q", got, tt.want)
			}
		})
	}
}
