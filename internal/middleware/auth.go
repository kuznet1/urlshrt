package middleware

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/kuznet1/urlshrt/internal/config"
	"github.com/kuznet1/urlshrt/internal/repository"
	"go.uber.org/zap"
	"net/http"
)

// CookieName is the name of the cookie that carries the JWT with the user identity.
var CookieName = "token"

// Auth issues and validates per-user JWT cookies and stores the user id in the request context.
// If an incoming request has no valid cookie, a new user is created via the repository and a token is set.
type Auth struct {
	repo   repository.Repo
	cfg    config.Config
	logger *zap.Logger
}

// NewAuth creates the authentication middleware using the provided config, repository and logger.
func NewAuth(repo repository.Repo, cfg config.Config, logger *zap.Logger) *Auth {
	return &Auth{repo: repo, cfg: cfg, logger: logger}
}

// Claims contains the JWT payload used by the authentication middleware.
// It includes the numeric UserID and standard registered claims.
type Claims struct {
	jwt.RegisteredClaims
	UserID int
}

// Authentication is an HTTP middleware that authenticates the request using a JWT cookie.
// On success it injects the user id into the context and refreshes the cookie when needed.
func (auth *Auth) Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(CookieName)
		if err != nil && err != http.ErrNoCookie {
			auth.internalError("unable to get cookie", err, w)
			return
		}

		var userID int

		if err == nil {
			claims, err := auth.parseToken(cookie.Value)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			userID = claims.UserID
		}

		if err == http.ErrNoCookie {
			userID, err = auth.repo.CreateUser(r.Context())
			if err != nil {
				auth.internalError("unable to create user", err, w)
			}

			token, err := auth.createToken(Claims{UserID: userID})
			if err != nil {
				auth.internalError("unable to create token", err, w)
			}

			http.SetCookie(w, &http.Cookie{
				Name:     CookieName,
				Value:    token,
				Path:     "/",
				HttpOnly: true,
			})
		}

		ctx := context.WithValue(r.Context(), repository.UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (auth *Auth) internalError(msg string, err error, w http.ResponseWriter) {
	auth.logger.Error(msg, zap.Error(err))
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (auth *Auth) parseToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if t.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Method.Alg())
			}
			return []byte(auth.cfg.SecretKey), nil
		})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func (auth *Auth) createToken(claims Claims) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(auth.cfg.SecretKey))
}
