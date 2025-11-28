package auth

import "context"

type contextKey string

const userIDKey contextKey = "user_id"

// setUserID adds user ID to context
func setUserID(ctx context.Context, userID int) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// GetUserID retrieves user ID from context
func GetUserID(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value(userIDKey).(int)
	return userID, ok
}
