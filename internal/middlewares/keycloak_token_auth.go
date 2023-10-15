package middlewares

import (
	"context"
	"errors"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	keycloakclient "github.com/lapitskyss/chat-service/internal/clients/keycloak"
	"github.com/lapitskyss/chat-service/internal/types"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/introspector_mock.gen.go -package=middlewaresmocks Introspector

const tokenCtxKey = "user-token"

var (
	ErrTokenNotActive         = errors.New("token not active")
	ErrNoRequiredResourceRole = errors.New("no required resource role")
)

type Introspector interface {
	IntrospectToken(ctx context.Context, token string) (*keycloakclient.IntrospectTokenResult, error)
}

// NewKeycloakTokenAuth returns a middleware that implements "active" authentication:
// each request is verified by the Keycloak server.
func NewKeycloakTokenAuth(introspector Introspector, resource, role string) echo.MiddlewareFunc {
	return middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup:  "header:Authorization",
		AuthScheme: "Bearer",
		Validator: func(tokenStr string, c echo.Context) (bool, error) {
			// Introspect JWT token in sentry
			iToken, err := introspector.IntrospectToken(c.Request().Context(), tokenStr)
			if err != nil {
				return false, err
			}
			if !iToken.Active {
				return false, ErrTokenNotActive
			}

			// Parse JWT token unverified
			jwtToken, tokenClaims, err := parseTokenUnverified(tokenStr)
			if err != nil {
				return false, err
			}

			// Validate JWT token, resource role
			if err = tokenClaims.Valid(); err != nil {
				return false, err
			}
			if err = checkTokenResourceRole(tokenClaims, resource, role); err != nil {
				return false, err
			}

			// Save token to context
			c.Set(tokenCtxKey, jwtToken)
			return true, nil
		},
	})
}

func parseTokenUnverified(tokenStr string) (*jwt.Token, *claims, error) {
	tokenClaims := &claims{}
	jwtToken, _, err := new(jwt.Parser).ParseUnverified(tokenStr, tokenClaims)
	if err != nil {
		return nil, nil, err
	}
	return jwtToken, tokenClaims, nil
}

func checkTokenResourceRole(tokenClaims *claims, resource, role string) error {
	chatClient, exist := tokenClaims.ResourceAccess[resource]
	if !exist {
		return ErrNoRequiredResourceRole
	}
	roles, exist := chatClient["roles"]
	if !exist {
		return ErrNoRequiredResourceRole
	}
	exist = false
	for _, r := range roles {
		if r == role {
			exist = true
		}
	}
	if !exist {
		return ErrNoRequiredResourceRole
	}
	return nil
}

func MustUserID(c echo.Context) types.UserID {
	uid, ok := userID(c)
	if !ok {
		panic("no user token in request context")
	}
	return uid
}

func userID(c echo.Context) (types.UserID, bool) {
	t := c.Get(tokenCtxKey)
	if t == nil {
		return types.UserIDNil, false
	}

	tt, ok := t.(*jwt.Token)
	if !ok {
		return types.UserIDNil, false
	}

	userIDProvider, ok := tt.Claims.(interface{ UserID() types.UserID })
	if !ok {
		return types.UserIDNil, false
	}
	return userIDProvider.UserID(), true
}

func userIDString(c echo.Context) string {
	uid, _ := userID(c)
	if uid.IsZero() {
		return ""
	}
	return uid.String()
}
