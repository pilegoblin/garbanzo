package context

import "context"

type ContextKey string

const userIDContextKey ContextKey = "userID"
const authIDContextKey ContextKey = "authID"

func GetUserID(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value(userIDContextKey).(int)
	return userID, ok
}

func SetUserID(ctx context.Context, userID int) context.Context {
	return context.WithValue(ctx, userIDContextKey, userID)
}

func GetAuthID(ctx context.Context) (string, bool) {
	authID, ok := ctx.Value(authIDContextKey).(string)
	return authID, ok
}

func SetAuthID(ctx context.Context, authID string) context.Context {
	return context.WithValue(ctx, authIDContextKey, authID)
}
