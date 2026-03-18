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
	attrs := attribute.NewSet(parameters.Attributes...)

	// Check if any filter wants to drop the trace
	for _, f := range signal.RegisteredFilters() {
		if f.ShouldDrop(parameters.ParentContext, attrs) {
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
		Tracestate: psc.TraceState(),
	}
}

func (s FilterBasedSampler) Description() string {
	registered := signal.RegisteredFilters()

	descriptions := make([]string, 0, len(registered))
	for _, f := range registered {
		descriptions = append(descriptions, f.Description(context.Background()))
	}

	descriptionsStr := strings.Join(descriptions, ", ")

	return "Iterates on registered attributes based filters. Filters: " + descriptionsStr
}
