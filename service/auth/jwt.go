package auth

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gonan98/go-ecommerce-api/config"
	"github.com/gonan98/go-ecommerce-api/helper"
	"github.com/gonan98/go-ecommerce-api/types"
)

type contextKey string

const UserKey contextKey = "userId"

func WithJWTAuth(apiHandler helper.APIHandler, store types.UserStore) helper.APIHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		tokenString := helper.GetTokenFromRequest(r)
		token, err := validateJWT(tokenString)
		if err != nil {
			return helper.PermissionDenied()
		}

		if !token.Valid {
			return helper.PermissionDenied()
		}

		claims := token.Claims.(jwt.MapClaims)
		str := claims["userId"].(string)
		userId, err := strconv.Atoi(str)
		if err != nil {
			return helper.PermissionDenied()
		}

		u, err := store.GetById(userId)
		if err != nil {
			return helper.PermissionDenied()
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, UserKey, u.ID)
		r = r.WithContext(ctx)

		return apiHandler(w, r)
	}
}

func GenerateJWT(secret []byte, userId int) (string, error) {
	expiration := time.Second * time.Duration(config.Envs.JWTExpirationInSeconds)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":    strconv.Itoa(userId),
		"expiresAt": time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, err
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signin method: %v", token.Header["alg"])
		}

		return []byte(config.Envs.JWTSecret), nil
	})
}

func GetUserIDFromContext(ctx context.Context) int {
	userId, ok := ctx.Value(UserKey).(int)
	if !ok {
		return -1
	}

	return userId
}
