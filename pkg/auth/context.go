package auth

import "context"

type contextKey string

const userIDKey contextKey = "user_id"

func SetUserID(ctx context.Context, id uint) context.Context {
	return context.WithValue(ctx, userIDKey, id)
}

func GetUserID(ctx context.Context) (uint, bool) {
	id, ok := ctx.Value(userIDKey).(uint)
	return id, ok
}
