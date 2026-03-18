package trace_test

import (
	"testing"

	"github.com/thumbrise/otelext/signal/trace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type staticSampler struct {
	d    sdktrace.SamplingDecision
	name string
}

func (s staticSampler) ShouldSample(parameters sdktrace.SamplingParameters) sdktrace.SamplingResult {
	return sdktrace.SamplingResult{Decision: s.d}
}

func (s staticSampler) Description() string {
	return s.name
}

func TestCompositeSampler_MixedDecisions_RecordOnlyThenRecordAndSample(t *testing.T) {
	s1 := staticSampler{d: sdktrace.RecordOnly, name: "RecordOnly"}
	s2 := staticSampler{d: sdktrace.RecordAndSample, name: "RecordAndSample"}
	cs := trace.NewCompositeSampler(s1, s2)

	res := cs.ShouldSample(sdktrace.SamplingParameters{})
	if res.Decision != sdktrace.RecordOnly {
		t.Fatalf("expected RecordOnly, got %v", res.Decision)
	}
}

func TestCompositeSampler_MixedDecisions_RecordAndSampleThenRecordOnly(t *testing.T) {
	s1 := staticSampler{d: sdktrace.RecordAndSample, name: "RecordAndSample"}
	s2 := staticSampler{d: sdktrace.RecordOnly, name: "RecordOnly"}
	cs := trace.NewCompositeSampler(s1, s2)

	res := cs.ShouldSample(sdktrace.SamplingParameters{})
	if res.Decision != sdktrace.RecordOnly {
		t.Fatalf("expected RecordOnly, got %v", res.Decision)
	}
}

func TestCompositeSampler_MixedDecisions_RecordAndSampleBoth(t *testing.T) {
	s1 := staticSampler{d: sdktrace.RecordAndSample, name: "R&A1"}
	s2 := staticSampler{d: sdktrace.RecordAndSample, name: "R&A2"}
	cs := trace.NewCompositeSampler(s1, s2)

	res := cs.ShouldSample(sdktrace.SamplingParameters{})
	if res.Decision != sdktrace.RecordAndSample {
		t.Fatalf("expected RecordAndSample, got %v", res.Decision)
	}
}

func TestCompositeSampler_MixedDecisions_WithDrop(t *testing.T) {
	s1 := staticSampler{d: sdktrace.RecordAndSample, name: "R&A"}
	s2 := staticSampler{d: sdktrace.Drop, name: "Drop"}
	cs := trace.NewCompositeSampler(s1, s2)

	res := cs.ShouldSample(sdktrace.SamplingParameters{})
	if res.Decision != sdktrace.Drop {
		t.Fatalf("expected Drop, got %v", res.Decision)
	}
}

func TestCompositeSampler_NoSamplers(t *testing.T) {
	cs := trace.NewCompositeSampler()

	res := cs.ShouldSample(sdktrace.SamplingParameters{})
	if res.Decision != sdktrace.Drop {
		t.Fatalf("expected Drop, got %v", res.Decision)
	}
}

func TestCompositeSampler_Description(t *testing.T) {
	t.Run("no samplers", func(t *testing.T) {
		cs := trace.NewCompositeSampler()
		got := cs.Description()
		want := "no samplers passed in composite sampler"
		if got != want {
			t.Fatalf("Description() = %q, want %q", got, want)
		}
	})

	t.Run("with samplers", func(t *testing.T) {
		s1 := staticSampler{d: sdktrace.RecordAndSample, name: "sampler1"}
		s2 := staticSampler{d: sdktrace.RecordAndSample, name: "sampler2"}
		cs := trace.NewCompositeSampler(s1, s2)
		got := cs.Description()
		want := "Decorates chain of samplers: sampler1\nsampler2"
		if got != want {
			t.Fatalf("Description() = %q, want %q", got, want)
		}
	})
}

// TestSamplingDecisionOrdering asserts the OTel SDK enum ordering that
// CompositeSampler's restrict-only logic depends on:
// Drop(0) < RecordOnly(1) < RecordAndSample(2).
func TestSamplingDecisionOrdering(t *testing.T) {
	if !(sdktrace.Drop < sdktrace.RecordOnly && sdktrace.RecordOnly < sdktrace.RecordAndSample) {
		t.Fatalf("unexpected SamplingDecision ordering: Drop=%d, RecordOnly=%d, RecordAndSample=%d",
			sdktrace.Drop, sdktrace.RecordOnly, sdktrace.RecordAndSample)
	}
}
