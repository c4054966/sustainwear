package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"sustainwear/internal/config"
	domainuser "sustainwear/internal/domain/user"
	"sustainwear/pkg/jwt"
)

type contextKey string

const (
	UserIDKey    contextKey = "user_id"
	UserEmailKey contextKey = "email"
	UserRoleKey  contextKey = "role"
)

// AUTH MIDDLEWARE - VALIDATES JWT TOKEN
func AuthMiddleware(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			// EXTRACT TOKEN FROM "Bearer <token>"
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
				return
			}

			token := parts[1]

			// VALIDATE TOKEN
			claims, err := jwt.ValidateToken(token, cfg.Security.JWTSecret)
			if err != nil {
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			// STORE USER INFO IN CONTEXT
			ctx := r.Context()
			ctx = context.WithValue(ctx, UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, UserEmailKey, claims.Email)
			ctx = context.WithValue(ctx, UserRoleKey, claims.Role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// ROLE-BASED AUTHORIZATION MIDDLEWARE
func RequireRole(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role := GetUserRole(r)

			allowed := false
			for _, allowedRole := range allowedRoles {
				if role == allowedRole {
					allowed = true
					break
				}
			}

			if !allowed {
				http.Error(w, "Forbidden: insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// GETS USER ID FROM CONTEXT
func GetUserID(r *http.Request) uint {
	userID, ok := r.Context().Value(UserIDKey).(uint)
	if !ok {
		return 0
	}
	return userID
}

// GETS USER EMAIL FROM CONTEXT
func GetUserEmail(r *http.Request) string {
	email, ok := r.Context().Value(UserEmailKey).(string)
	if !ok {
		return ""
	}
	return email
}

// GETS USER ROLE FROM CONTEXT
func GetUserRole(r *http.Request) string {
	role, ok := r.Context().Value(UserRoleKey).(string)
	if !ok {
		return ""
	}
	return role
}

// GETS ORG ID BASED ON USER ROLE & QUERIES DB TO GET ORG ID (SECURITY PURPOSES)
func GetOrgIDForRequest(r *http.Request, getUserByID func(uint) (*domainuser.User, error)) (uint, error) {
	userID := GetUserID(r)
	role := GetUserRole(r)

	if role == "admin" {
		orgIDStr := r.URL.Query().Get("org_id")
		if orgIDStr == "" {
			return 0, fmt.Errorf("admin must specify org_id parameter")
		}
		orgID, err := strconv.ParseUint(orgIDStr, 10, 32)
		if err != nil {
			return 0, fmt.Errorf("invalid org_id")
		}
		return uint(orgID), nil
	}

	if role == "charity_staff" {
		appUser, err := getUserByID(userID)
		if err != nil {
			return 0, fmt.Errorf("failed to get user info")
		}
		if appUser.OrganisationID == nil {
			return 0, fmt.Errorf("charity staff must be assigned to an organisation")
		}
		return *appUser.OrganisationID, nil
	}

	return 0, fmt.Errorf("access denied: insufficient permissions")
}
