package retry

import (
	"context"
	"testing"
)

func BenchmarkDoSuccess(b *testing.B) {
	ctx := context.Background()

	b.ResetTimer()

	for b.Loop() {
		_ = Do(
			ctx,
			3,
			0,
			func(ctx context.Context) error {
				return nil
			},
		)
	}
}
