package auth

import "context"

type principalKey struct{}

func SetPrincipal(ctx context.Context, p *Principal) context.Context {
	return context.WithValue(ctx, principalKey{}, p)
}

func GetPrincipal(ctx context.Context) *Principal {
	p, ok := ctx.Value(principalKey{}).(*Principal)
	if !ok {
		return nil
	}
	return p
}
