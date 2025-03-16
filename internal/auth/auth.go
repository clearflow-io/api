package auth

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/igorschechtel/finance-tracker-backend/internal/config"
	u "github.com/igorschechtel/finance-tracker-backend/internal/utils"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jwt"
	"github.com/sirupsen/logrus"
)

// TokenData contains information extracted from the JWT token
type TokenData struct {
	UserID string
	// Add other fields you might need from the token
}

type contextKey string

const TokenDataKey contextKey = "auth.tokenData"

// JWKS caching variables
var (
	jwksCache     jwk.Set
	jwksMutex     sync.RWMutex
	lastFetchTime time.Time
	cacheDuration = 24 * time.Hour
)

// getJWKS fetches JWKS from the URL with caching
func getJWKS(ctx context.Context, client *http.Client, jwksURL string) (jwk.Set, error) {
	// Try to use cache first
	jwksMutex.RLock()
	if jwksCache != nil && time.Since(lastFetchTime) < cacheDuration {
		defer jwksMutex.RUnlock()
		return jwksCache, nil
	}
	jwksMutex.RUnlock()

	// Cache is empty or expired, acquire write lock and update
	jwksMutex.Lock()
	defer jwksMutex.Unlock()

	// Double-check to handle race conditions
	if jwksCache != nil && time.Since(lastFetchTime) < cacheDuration {
		return jwksCache, nil
	}

	// Create a context with timeout for the fetch operation
	fetchCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Fetch the JWKS
	keySet, err := jwk.Fetch(fetchCtx, jwksURL, jwk.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JWKS: %w", err)
	}

	// Update cache
	jwksCache = keySet
	lastFetchTime = time.Now()
	logrus.Info("JWKS cache refreshed")

	return keySet, nil
}

// AuthMiddleware creates a middleware for JWT validation
func AuthMiddleware(httpClient *http.Client, cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeaders := r.Header.Get("Authorization")
			if authHeaders == "" {
				logrus.Warn("Authorization header missing")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			accessToken := strings.TrimPrefix(authHeaders, "Bearer ")
			if accessToken == "" {
				logrus.Warn("Bearer token missing")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			tokenData, err := extractTokenData(r.Context(), accessToken, httpClient, cfg)
			if err != nil {
				logrus.WithError(err).Error("Token validation failed")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Add token data to request context
			ctx := context.WithValue(r.Context(), TokenDataKey, tokenData)

			// Only log userID, not the entire token data
			logrus.WithField("userID", tokenData.UserID).Info("User authenticated")

			// Call the next handler with the updated context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// extractTokenData validates the JWT and extracts relevant data
func extractTokenData(ctx context.Context, accessToken string, httpClient *http.Client, cfg *config.Config) (*TokenData, error) {
	// Create a JWKS URL for the project
	jwksURL := fmt.Sprintf("https://api.stack-auth.com/api/v1/projects/%s/.well-known/jwks.json", cfg.StackAuth.ProjectID)
	logrus.WithField("jwksURL", jwksURL).Info("Fetching JWKS")

	// Get the JWKS key set (from cache if available)
	keySet, err := getJWKS(ctx, httpClient, jwksURL)
	if err != nil {
		return nil, err
	}

	// Parse and validate the JWT token
	token, err := jwt.Parse([]byte(accessToken), jwt.WithKeySet(keySet))
	if err != nil {
		return nil, errors.New("invalid access token")
	}

	// Validate critical claims (e.g., expiration and audience)
	if err := jwt.Validate(token, jwt.WithClock(jwt.ClockFunc(time.Now))); err != nil {
		return nil, fmt.Errorf("token validation failed: %w", err)
	}

	// Extract user ID from the token (typically in the "sub" claim)
	var userIDStr string
	if err := token.Get("sub", &userIDStr); err != nil {
		return nil, errors.New("user ID not found in token")
	}

	// Create and return the token data
	tokenData := &TokenData{
		UserID: userIDStr,
		// Extract other fields as needed
	}

	return tokenData, nil
}

func GetTokenData(r *http.Request) *TokenData {
	value := r.Context().Value(TokenDataKey)
	if value == nil {
		logrus.Warn("Token data not found in context")
		return nil
	}

	tokenData, ok := value.(*TokenData)
	if !ok {
		logrus.Error("Token data is not of expected type")
		return nil
	}

	return tokenData
}

func GetUserID(r *http.Request) (uuid.UUID, error) {
	tokenData := GetTokenData(r)
	if tokenData == nil {
		return uuid.Nil, errors.New("token data not found in context")
	}

	return u.ParseUUID(tokenData.UserID, "UserID")
}

func InitializeHTTPClient() (*http.Client, error) {
	// Create a custom transport with enhanced security settings
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			},
			PreferServerCipherSuites: true,
		},
		MaxIdleConns:    10,
		IdleConnTimeout: 30 * time.Second,
	}

	// Create HTTP client with timeout and secure transport
	httpClient := &http.Client{
		Timeout:   10 * time.Second,
		Transport: transport,
	}

	return httpClient, nil
}
