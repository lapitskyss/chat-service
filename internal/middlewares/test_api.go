package middlewares

import (
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"

	"github.com/lapitskyss/chat-service/internal/types"
)

func SetToken(c echo.Context, uid types.UserID) {
	c.Set(tokenCtxKey, jwt.NewWithClaims(jwt.SigningMethodRS256, claims{
		StandardClaims: jwt.StandardClaims{},
		Audience:       nil,
		Subject:        uid,
		ResourceAccess: nil,
	}))
}
