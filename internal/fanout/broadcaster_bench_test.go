package fanout_test

import (
	"context"
	"testing"

	"github.com/user/logdrift/internal/fanout"
	"github.com/user/logdrift/internal/tail"
)

func BenchmarkBroadcast_TwoSubs(b *testing.B) {
	src := make(chan tail.Entry, b.N+1)
	for i := 0; i < b.N; i++ {
		src <- tail.Entry{Service: "bench", Line: "x"}
	}
	close(src)

	br, _ := fanout.New(src, b.N+1)
	s1 := br.Subscribe()
	s2 := br.Subscribe()

	b.ResetTimer()
	ctx := context.Background()
	go br.Run(ctx)

	for i := 0; i < b.N; i++ {
		<-s1
		<-s2
	}
}
