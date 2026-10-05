package retry

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestDoSuccessOnFirstAttempt(t *testing.T) {
	attempts := 0

	err := Do(
		context.Background(),
		3,
		time.Millisecond,
		func(ctx context.Context) error {
			attempts++
			return nil
		},
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if attempts != 1 {
		t.Fatalf("unexpected attempts: got %d, want 1", attempts)
	}
}

func TestDoRetriesAfterError(t *testing.T) {
	attempts := 0
	expectedErr := errors.New("temporary error")

	err := Do(
		context.Background(),
		3,
		time.Millisecond,
		func(ctx context.Context) error {
			attempts++

			if attempts < 3 {
				return expectedErr
			}

			return nil
		},
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if attempts != 3 {
		t.Fatalf("unexpected attempts: got %d, want 3", attempts)
	}
}

func TestDoReturnsLastError(t *testing.T) {
	attempts := 0
	expectedErr := errors.New("permanent error")

	err := Do(
		context.Background(),
		3,
		time.Millisecond,
		func(ctx context.Context) error {
			attempts++
			return expectedErr
		},
	)

	if !errors.Is(err, expectedErr) {
		t.Fatalf("unexpected error: got %v, want %v", err, expectedErr)
	}

	if attempts != 3 {
		t.Fatalf("unexpected attempts: got %d, want 3", attempts)
	}
}

func TestDoStopsWhenContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	attempts := 0

	err := Do(
		ctx,
		3,
		time.Second,
		func(ctx context.Context) error {
			attempts++
			return errors.New("temporary error")
		},
	)

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("unexpected error: got %v, want %v", err, context.Canceled)
	}

	if attempts != 1 {
		t.Fatalf("unexpected attempts: got %d, want 1", attempts)
	}
}
