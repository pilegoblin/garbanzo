package context

import "context"

type ContextKey string

const emailContextKey ContextKey = "email"

func GetEmail(ctx context.Context) (string, bool) {
	email, ok := ctx.Value(emailContextKey).(string)
	return email, ok
}

func SetEmail(ctx context.Context, email string) context.Context {
	return context.WithValue(ctx, emailContextKey, email)
}
