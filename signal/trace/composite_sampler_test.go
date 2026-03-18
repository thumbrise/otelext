// Copyright 2026 thumbrise
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
