package api

import "context"

// iter loops over paginated apis to give a complete response.
//
// get will be called once per iteration and must be reentrant.
func iter(ctx context.Context, get func(ctx context.Context, opts ...callOption) (uint, error)) error {
	for i, total := uint(1), uint(2); i < total; i++ {
		var err error
		switch i {
		case 0, 1:
			total, err = get(ctx)
		default:
			total, err = get(ctx, filterPage(i))
		}
		if err != nil {
			return err
		}
	}
	return nil
}
