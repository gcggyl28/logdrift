// Package session wires together the tail, diff, and render subsystems into
// a single cohesive runtime for a logdrift invocation.
//
// Typical usage:
//
//	cfg, _ := config.Load("logdrift.yaml")
//	s, _ := session.New(cfg, os.Stdout)
//	s.Run(ctx)
//
// Session reads the list of services from the config, opens a file tailer for
// each one, fans the streams together, runs them through the diff pipeline,
// and writes any detected drift to the configured output renderer.
package session
