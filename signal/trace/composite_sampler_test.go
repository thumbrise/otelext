package trace

import (
"testing"

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
cs := NewCompositeSampler(s1, s2)
res := cs.ShouldSample(sdktrace.SamplingParameters{})
if res.Decision != sdktrace.RecordOnly {
t.Fatalf("expected RecordOnly, got %v", res.Decision)
}
}

func TestCompositeSampler_MixedDecisions_RecordAndSampleThenRecordOnly(t *testing.T) {
s1 := staticSampler{d: sdktrace.RecordAndSample, name: "RecordAndSample"}
s2 := staticSampler{d: sdktrace.RecordOnly, name: "RecordOnly"}
cs := NewCompositeSampler(s1, s2)
res := cs.ShouldSample(sdktrace.SamplingParameters{})
if res.Decision != sdktrace.RecordOnly {
t.Fatalf("expected RecordOnly, got %v", res.Decision)
}
}

func TestCompositeSampler_MixedDecisions_RecordAndSampleBoth(t *testing.T) {
s1 := staticSampler{d: sdktrace.RecordAndSample, name: "R&A"}
s2 := staticSampler{d: sdktrace.RecordAndSample, name: "R&A2"}
cs := NewCompositeSampler(s1, s2)
res := cs.ShouldSample(sdktrace.SamplingParameters{})
if res.Decision != sdktrace.RecordAndSample {
t.Fatalf("expected RecordAndSample, got %v", res.Decision)
}
}

func TestCompositeSampler_MixedDecisions_WithDrop(t *testing.T) {
s1 := staticSampler{d: sdktrace.RecordAndSample, name: "R&A"}
s2 := staticSampler{d: sdktrace.Drop, name: "Drop"}
cs := NewCompositeSampler(s1, s2)
res := cs.ShouldSample(sdktrace.SamplingParameters{})
if res.Decision != sdktrace.Drop {
t.Fatalf("expected Drop, got %v", res.Decision)
}
}
