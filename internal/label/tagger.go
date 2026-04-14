package label

import "fmt"

// Entry is a log line annotated with a resolved canonical service label.
type Entry struct {
	Service string
	Line    string
}

// Tagger resolves raw service identifiers to canonical labels before
// forwarding entries downstream.
type Tagger struct {
	reg *Registry
}

// NewTagger constructs a Tagger backed by the given Registry.
func NewTagger(reg *Registry) *Tagger {
	return &Tagger{reg: reg}
}

// Tag resolves nameOrAlias to its canonical label and returns an Entry.
// Returns an error if the name cannot be resolved.
func (t *Tagger) Tag(nameOrAlias, line string) (Entry, error) {
	canon, err := t.reg.Resolve(nameOrAlias)
	if err != nil {
		return Entry{}, fmt.Errorf("tagger: %w", err)
	}
	return Entry{Service: canon, Line: line}, nil
}

// TagAll applies Tag to each (service, line) pair, skipping unresolvable
// services and collecting errors. All valid entries are returned alongside
// any accumulated errors.
func (t *Tagger) TagAll(pairs [][2]string) ([]Entry, []error) {
	var entries []Entry
	var errs []error
	for _, p := range pairs {
		e, err := t.Tag(p[0], p[1])
		if err != nil {
			errs = append(errs, err)
			continue
		}
		entries = append(entries, e)
	}
	return entries, errs
}
