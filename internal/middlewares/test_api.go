package middlewares

import (
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"

	"github.com/lapitskyss/chat-service/internal/types"
)

func AuthWith(uid types.UserID) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (returnErr error) {
			SetToken(c, uid)
			return next(c)
		}
	}
}

func SetToken(c echo.Context, uid types.UserID) {
	c.Set(tokenCtxKey, &jwt.Token{Claims: claimsMock{uid: uid}, Valid: true})
}

type claimsMock struct {
	uid types.UserID
}

func (m claimsMock) Valid() error {
	return nil
}

func (m claimsMock) UserID() types.UserID {
	return m.uid
}
