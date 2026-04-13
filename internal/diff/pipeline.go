package diff

import (
	"context"

	"github.com/user/logdrift/internal/tail"
)

// PairConfig describes which two services to diff against each other.
type PairConfig struct {
	Left  string
	Right string
}

// Pipeline reads from a fan-in channel, pairs entries by service name,
// and emits Results on the returned channel.
func Pipeline(ctx context.Context, in <-chan tail.LogLine, pair PairConfig, d *Differ) <-chan Result {
	out := make(chan Result, 16)

	go func() {
		defer close(out)

		buf := make(map[string][]string)

		for {
			select {
			case <-ctx.Done():
				return
			case line, ok := <-in:
				if !ok {
					return
				}
				svc := line.Service
				if svc != pair.Left && svc != pair.Right {
					continue
				}

				buf[svc] = append(buf[svc], line.Text)

				// Emit a result whenever both sides have a pending line.
				for len(buf[pair.Left]) > 0 && len(buf[pair.Right]) > 0 {
					a := Entry{Service: pair.Left, Line: buf[pair.Left][0]}
					b := Entry{Service: pair.Right, Line: buf[pair.Right][0]}
					buf[pair.Left] = buf[pair.Left][1:]
					buf[pair.Right] = buf[pair.Right][1:]

					res := d.Compare(a, b)
					select {
					case out <- res:
					case <-ctx.Done():
						return
					}
				}
			}
		}
	}()

	return out
}
