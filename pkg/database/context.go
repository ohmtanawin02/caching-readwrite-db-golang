package database

import "context"

type usePrimaryKey struct{}

// WithPrimary returns a context that signals the repository to read from writeDB (primary).
// Use this after a write operation to avoid replication lag on the replica.
func WithPrimary(ctx context.Context) context.Context {
	return context.WithValue(ctx, usePrimaryKey{}, true)
}

// UsePrimary reports whether the context requests a primary DB read.
func UsePrimary(ctx context.Context) bool {
	v, _ := ctx.Value(usePrimaryKey{}).(bool)
	return v
}
