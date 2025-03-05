package auth

import (
	"context"
	"net/http"
	"newsletter/internal/utils"
)

// Key type for context values
type contextKey string

// Context keys
const (
	UserIDKey   contextKey = "user_id"
	UsernameKey contextKey = "username"
)

// AuthMiddleware is a middleware that checks if the request has a valid JWT token
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract token from request
		tokenString, err := ExtractTokenFromRequest(r)
		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized: "+err.Error())
			return
		}

		// Validate token
		claims, err := ValidateToken(tokenString)
		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized: "+err.Error())
			return
		}

		// Add user information to request context
		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, UsernameKey, claims.Username)

		// Call the next handler with the updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserIDFromContext retrieves the user ID from the request context
func GetUserIDFromContext(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(UserIDKey).(int64)
	return userID, ok
}

// GetUsernameFromContext retrieves the username from the request context
func GetUsernameFromContext(ctx context.Context) (string, bool) {
	username, ok := ctx.Value(UsernameKey).(string)
	return username, ok
}

// RequireAuth is a helper function to create a handler that requires authentication
func RequireAuth(handler http.HandlerFunc) http.Handler {
	return AuthMiddleware(http.HandlerFunc(handler))
}
