// Package main is the entry point for the logdrift CLI tool.
// It wires together configuration, session management, and signal handling
// to tail and diff log streams across multiple services in real time.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/yourorg/logdrift/internal/config"
	"github.com/yourorg/logdrift/internal/session"
)

const version = "0.1.0"

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "logdrift: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	var (
		cfgPath     = flag.String("config", "logdrift.yaml", "path to configuration file")
		showVersion = flag.Bool("version", false, "print version and exit")
		verbose     = flag.Bool("verbose", false, "enable verbose logging")
	)
	flag.Parse()

	if *showVersion {
		fmt.Printf("logdrift %s\n", version)
		return nil
	}

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		return fmt.Errorf("loading config %q: %w", *cfgPath, err)
	}

	if *verbose {
		fmt.Fprintf(os.Stderr, "logdrift %s starting with %d service(s)\n",
			version, len(cfg.Services))
	}

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel()

	sess, err := session.New(cfg)
	if err != nil {
		return fmt.Errorf("creating session: %w", err)
	}

	if err := sess.Run(ctx); err != nil {
		return fmt.Errorf("session exited: %w", err)
	}

	return nil
}
