package auth

import (
	"context"
	"net/http"
	"os"

	"github.com/clerk/clerk-sdk-go/v2"
	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
	"github.com/sirupsen/logrus"
)

// InitializeClerk sets up the Clerk API key for use throughout the app
func InitializeClerk() {
	secretKey := os.Getenv("CLERK_SECRET_KEY")
	if secretKey == "" {
		logrus.Fatal("CLERK_SECRET_KEY environment variable is not set")
	}

	// Set the API key for the clerk package
	clerk.SetKey(secretKey)
}

// ClerkAuthMiddleware returns middleware that verifies the session token
// and adds the session claims to the request context
func ClerkAuthMiddleware() func(http.Handler) http.Handler {
	return clerkhttp.WithHeaderAuthorization()
}

// GetUserID extracts the user ID from the request context (Clerk ID)
func GetUserID(ctx context.Context) (string, bool) {
	claims, ok := clerk.SessionClaimsFromContext(ctx)
	if !ok || claims.Subject == "" {
		return "", false
	}

	return claims.Subject, true
}

// RequireAuth is a convenience wrapper that returns a 401 if authentication fails
func RequireAuth(next http.Handler) http.Handler {
	return clerkhttp.RequireHeaderAuthorization()(next)
}
