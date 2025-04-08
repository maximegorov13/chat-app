package appcontext

import "context"

type contextKey string

const (
	contextUserIDKey contextKey = "ContextUserIDKey"
)

func SetContextUserID(ctx context.Context, userId string) context.Context {
	return context.WithValue(ctx, contextUserIDKey, userId)
}

func GetContextUserID(ctx context.Context) string {
	return ctx.Value(contextUserIDKey).(string)
}
