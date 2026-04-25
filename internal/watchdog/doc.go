// Package watchdog monitors registered log services for silence.
//
// A Watchdog tracks the last time a line was received from each service.
// When a service has been silent for longer than the configured threshold, a
// watchdog.Event is emitted on the channel returned by Run.
//
// # Basic usage
//
//	wd, err := watchdog.New(30 * time.Second)
//	if err != nil { ... }
//	wd.Register("api")
//	wd.Register("worker")
//
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//
//	for event := range wd.Run(ctx) {
//		log.Printf("service %s silent for %v", event.Service, event.SilentFor)
//	}
//
// # Integration with a log pipeline
//
// Use Runner to automatically call Ping for every log entry that flows
// through a channel, wiring watchdog monitoring into the existing pipeline
// with minimal boilerplate.
package watchdog
