package middleware

import (
	"net/http"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/labstack/echo/v4"
)

type ErrorResult struct {
	Error string `json:"error"`
}

func VerifyScope(scope string) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Request().Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)

			claims := token.CustomClaims.(*CustomClaims)
			if !claims.HasScope(scope) {
				return c.JSON(http.StatusUnauthorized, ErrorResult{Error: "user does not have sufficient permission"})
			}
			return next(c)
		}
	}
}
