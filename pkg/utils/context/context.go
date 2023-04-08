package contextUtils

import (
	"context"
	"errors"
)

func CheckContext(ctx context.Context) error {
	select {
	case <-ctx.Done():
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return ctx.Err()
		}
	default:
	}
	return nil
}
