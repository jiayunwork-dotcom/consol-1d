package model

import "context"

func abortSolveContext() error {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := ctx.Err(); err != nil {
		return err
	}
	return nil
}
