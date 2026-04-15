// Package routing implements pattern-based log-line routing for logdrift.
//
// A Router is constructed from a map of destination labels to regular-expression
// patterns. At runtime, each incoming service name is tested against every
// registered pattern; all destinations whose pattern matches are returned by
// Match, allowing a single log line to be fanned out to multiple output groups.
//
// Example usage:
//
//	r, err := routing.New(map[string]string{
//		"frontend": "^(api|web)-",
//		"backend":  "^(worker|queue)-",
//		"all":      ".*",
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	destinations := r.Match("api-gateway") // ["frontend", "all"]
package routing
