package auth

import (
	"context"
	"net/http"
	"strings"

	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type AuthMiddleware interface {
	Middleware() gin.HandlerFunc
}

type FirebaseAuthMiddleware struct {
	AuthClient *auth.Client
}

func NewFirebaseAuthMiddleware(authClient *auth.Client) AuthMiddleware {
	return FirebaseAuthMiddleware{AuthClient: authClient}
}

func (a FirebaseAuthMiddleware) Middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bearerToken := a.tokenFromHeader(ctx.Request)
		if bearerToken == "" {
			log.Error(errors.New("empty bearer token"))
			ctx.String(http.StatusUnauthorized, "Unauthorized")
			ctx.Abort()
			return
		}

		token, err := a.AuthClient.VerifyIDToken(ctx, bearerToken)
		if err != nil {
			log.Error(errors.Wrap(err, "unable to verify jwt"))
			ctx.String(http.StatusUnauthorized, "Unauthorized")
			ctx.Abort()
			return
		}

		ctx.Request = ctx.Request.WithContext(context.WithValue(ctx, userContextKey, User{
			UUID:  token.UID,
			Email: token.Claims["email"].(string),
		}))

		ctx.Next()
	}
}

func (a FirebaseAuthMiddleware) tokenFromHeader(r *http.Request) string {
	headerValue := r.Header.Get("Authorization")

	if len(headerValue) > 7 && strings.ToLower(headerValue[0:6]) == "bearer" {
		return headerValue[7:]
	}

	return ""
}

type User struct {
	UUID  string
	Email string
	Role  string

	DisplayName string
}

type ctxKey int

const (
	userContextKey ctxKey = iota
)

var (
	ErrNoUserInContext = errors.New("auth error: no user in context")
)

func UserFromCtx(ctx context.Context) (User, error) {
	u, ok := ctx.Value(userContextKey).(User)
	if ok {
		return u, nil
	}

	return User{}, ErrNoUserInContext
}
