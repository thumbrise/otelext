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

// ShouldSample determines if a trace should be sampled. It iterates through the
// configured samplers in order. If any sampler returns sdktrace.Drop, the
// decision is immediately Drop. If a sampler returns sdktrace.RecordOnly,
// that decision is recorded, and subsequent samplers cannot upgrade it to
// sdktrace.RecordAndSample. If a sampler returns sdktrace.RecordAndSample,
// it is only considered if no prior sampler decided sdktrace.RecordOnly.
// This ensures that restrictions are respected and the behavior is "AND-like"
// where any sampler can restrict the decision, but not upgrade it beyond
// what a previous sampler allowed.
func (r CompositeSampler) ShouldSample(parameters sdktrace.SamplingParameters) sdktrace.SamplingResult {
       if len(r.samplers) == 0 {
               return sdktrace.SamplingResult{
                       Decision:   sdktrace.Drop,
                       Attributes: nil,
                       Tracestate: trace.TraceState{},
               }
       }

       var finalResult sdktrace.SamplingResult
       recordOnlyEncountered := false

       for _, sampler := range r.samplers {
               res := sampler.ShouldSample(parameters)
               if res.Decision == sdktrace.Drop {
                       // If any sampler decides to drop, the final decision is Drop.
                       return res
               }
               if res.Decision == sdktrace.RecordOnly {
                       // If a sampler decides RecordOnly, we record it and continue.
                       // Subsequent samplers cannot upgrade this to RecordAndSample.
                       finalResult = res
                       recordOnlyEncountered = true
               } else if res.Decision == sdktrace.RecordAndSample {
                       // If a sampler decides RecordAndSample, we only consider it if no prev RecordOnly
                       // decision has been encountered.
                       if !recordOnlyEncountered {
                               finalResult = res
                       }
               }
       }

       // If RecordOnly was encountered at any point, ensure final decision is RecordOnly.
       if recordOnlyEncountered {
               finalResult.Decision = sdktrace.RecordOnly
       }

       return finalResult
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
