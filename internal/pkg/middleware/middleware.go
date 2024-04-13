package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/Alladan04/avito_test/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

func ParseTokenPayload(token string) (models.JwtPayload, error) {
	var claims *jwt.Token
	var err error
	isAdmin := false
	claims, err = jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_USER_SECRET")), nil
	})
	if err != nil {
		claims, err = jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("JWT_ADMIN_SECRET")), nil
		})
		if err != nil {
			return models.JwtPayload{}, err
		}
		isAdmin = true
	}
	_, err = claims.Claims.GetExpirationTime()
	if err != nil {
		return models.JwtPayload{}, err
	}

	payloadMap, ok := claims.Claims.(jwt.MapClaims)
	if !ok {
		return models.JwtPayload{}, errors.New("invalid format (claims)")
	}

	username, ok := payloadMap["usr"].(string)
	if !ok {
		return models.JwtPayload{}, errors.New("invalid format (usr)")
	}

	return models.JwtPayload{

		Username: username,
		IsAdmin:  isAdmin,
	}, nil

}

func JwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		header := r.Header.Get("token")
		if header == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		token := headerParts[1]

		payload, err := ParseTokenPayload(token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), models.PayloadContextKey, payload)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func CheckAdminPermissionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwtPayload, ok := r.Context().Value(models.PayloadContextKey).(models.JwtPayload)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if !jwtPayload.IsAdmin {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
