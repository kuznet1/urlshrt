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

var CookieName = "token"

type Auth struct {
	repo   repository.Repo
	cfg    config.Config
	logger *zap.Logger
}

func NewAuth(repo repository.Repo, cfg config.Config, logger *zap.Logger) *Auth {
	return &Auth{repo: repo, cfg: cfg, logger: logger}
}

type Claims struct {
	jwt.RegisteredClaims
	UserID int
}

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
