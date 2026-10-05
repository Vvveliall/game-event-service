package retry

import (
	"context"
	"time"
)

func Do(
	ctx context.Context,
	attempts int,
	delay time.Duration,
	operation func(context.Context) error,
) error {
	if attempts <= 0 {
		return nil
	}

	var err error

	for attempt := 1; attempt <= attempts; attempt++ {
		err = operation(ctx)
		if err == nil {
			return nil
		}

		if attempt == attempts {
			break
		}

		timer := time.NewTimer(delay)

		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
		}
	}

	return err
}
