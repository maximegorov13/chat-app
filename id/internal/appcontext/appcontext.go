package appcontext

import "context"

type contextKey string

const (
	contextUserIdKey contextKey = "ContextUserIdKey"
)

func SetContextUserId(ctx context.Context, userId string) context.Context {
	return context.WithValue(ctx, contextUserIdKey, userId)
}

func GetContextUserId(ctx context.Context) string {
	return ctx.Value(contextUserIdKey).(string)
}
