package api

import "context"

// iter loops over paginated apis to give a complete response.
//
// get will be called once per iteration and must be reentrant.
func iter(ctx context.Context, get func(ctx context.Context, opts ...callOption) (uint, error)) error {
	for i, total := uint(0), uint(1); i < total; i++ {
		var err error
		total, err = get(ctx, filterPage(i+1))
		if err != nil {
			return err
		}
	}
	return nil
}
