package middleware

import (
	"context"
	"net/http"

	"peso/internal/application"
	"peso/internal/domain/user"
)

type ctxKey string

const (
	userCtxKey    ctxKey = "user"
	sessionCtxKey ctxKey = "session_token"
	CookieName    string = "peso_session"
)

func SessionMiddleware(authService *application.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(CookieName)
			if err != nil || cookie.Value == "" {
				next.ServeHTTP(w, r)
				return
			}

			u, err := authService.ValidateSession(cookie.Value)
			if err != nil {
				ClearSessionCookie(w)
				next.ServeHTTP(w, r)
				return
			}

			ctx := context.WithValue(r.Context(), userCtxKey, u)
			ctx = context.WithValue(ctx, sessionCtxKey, cookie.Value)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if UserFromContext(r.Context()) == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func UserFromContext(ctx context.Context) *user.User {
	u, ok := ctx.Value(userCtxKey).(*user.User)
	if !ok {
		return nil
	}
	return u
}

func SessionTokenFromContext(ctx context.Context) string {
	token, ok := ctx.Value(sessionCtxKey).(string)
	if !ok {
		return ""
	}
	return token
}

func SetSessionCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    token,
		Path:     "/",
		MaxAge:   30 * 24 * 60 * 60, // 30 days
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func ClearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}
