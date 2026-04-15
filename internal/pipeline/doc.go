// Package pipeline provides a composable, ordered chain of processing
// stages for log entries.
//
// Each Stage wraps a StageFn — a function that receives an Entry and
// returns either a (possibly mutated) Entry or nil to drop the entry
// from the stream.
//
// Typical usage:
//
//	redactStage, _ := pipeline.NewStage("redact", func(e pipeline.Entry) (*pipeline.Entry, error) {
//		e.Line = redactor.Apply(e.Line)
//		return &e, nil
//	})
//
//	p, errStage{redactStage, normalizeStage})
//	if err != nil { ... }
//
//	out, err := p.Run(entry)
//	if out == nil { /* entry was dropped */ }
//
// Stages are applied in the order they are supplied to New.
// The first stage that returns an error short-circuits the chain and
// the error is wrapped with the stage name for easy diagnosis.
package pipeline
