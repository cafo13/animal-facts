package auth

import (
	"context"
	"net/http"

	"firebase.google.com/go/v4/auth"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type MockAuthMiddleware struct {
	AuthClient *auth.Client
}

// MockAuthMiddleware can be used in the local environments for development. It doesn't depends on firebase.
func NewMockAuthMiddleware() AuthMiddleware {
	return MockAuthMiddleware{AuthClient: nil}
}

// MockAuthMiddleware is used in the local environment (which doesn't depend on Firebase)
func (a MockAuthMiddleware) Middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var claims jwt.MapClaims
		token, err := request.ParseFromRequest(
			ctx.Request,
			request.AuthorizationHeaderExtractor,
			func(token *jwt.Token) (i interface{}, e error) {
				return []byte("mock_secret"), nil
			},
			request.WithClaims(&claims),
		)
		if err != nil {
			log.Error(errors.Wrap(err, "unable to get jwt"))
			ctx.String(http.StatusUnauthorized, "Unauthorized")
			return
		}

		if !token.Valid {
			log.Error(errors.New("invalid jwt"))
			ctx.String(http.StatusUnauthorized, "Unauthorized")
			return
		}

		ctx.Request = ctx.Request.WithContext(context.WithValue(ctx, userContextKey, User{
			UUID:        claims["user_uuid"].(string),
			Email:       claims["email"].(string),
			Role:        claims["role"].(string),
			DisplayName: claims["name"].(string),
		}))

		ctx.Next()
	}
}
