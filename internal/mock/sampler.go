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
