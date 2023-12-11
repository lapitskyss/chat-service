package middlewares

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"

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
	validate := func(tokenStr string, c echo.Context) (bool, error) {
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
		if !hasTokenResourceRole(tokenClaims, resource, role) {
			return false, ErrNoRequiredResourceRole
		}

		// Save token to context
		c.Set(tokenCtxKey, jwtToken)
		return true, nil
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := getToken(c)
			if token == "" {
				return echo.NewHTTPError(http.StatusBadRequest, "invalid token")
			}

			valid, err := validate(token, c)
			if err != nil {
				return unauthorizedError(err)
			}
			if !valid {
				return echo.NewHTTPError(http.StatusBadRequest, "invalid token")
			}

			return next(c)
		}
	}
}

func unauthorizedError(err error) *echo.HTTPError {
	return &echo.HTTPError{
		Code:     http.StatusUnauthorized,
		Message:  "Unauthorized",
		Internal: err,
	}
}

func getToken(c echo.Context) string {
	token := getBearerTokenFromAuthHeader(c)
	if token != "" {
		return token
	}
	return getTokenFromWebSocketProtocolHeader(c)
}

func getBearerTokenFromAuthHeader(c echo.Context) string {
	authHeader := c.Request().Header.Get("Authorization")
	authFields := strings.Fields(authHeader)
	if len(authFields) != 2 || strings.ToLower(authFields[0]) != "bearer" {
		return ""
	}
	return authFields[1]
}

func getTokenFromWebSocketProtocolHeader(c echo.Context) string {
	header := c.Request().Header.Get("Sec-WebSocket-Protocol")
	values := strings.Split(header, ",")
	if len(values) < 2 {
		return ""
	}
	return strings.TrimSpace(values[1])
}

func parseTokenUnverified(tokenStr string) (*jwt.Token, *claims, error) {
	tokenClaims := &claims{}
	jwtToken, _, err := new(jwt.Parser).ParseUnverified(tokenStr, tokenClaims)
	if err != nil {
		return nil, nil, err
	}
	return jwtToken, tokenClaims, nil
}

func hasTokenResourceRole(tokenClaims *claims, resource, role string) bool {
	chatClient, exist := tokenClaims.ResourceAccess[resource]
	if !exist {
		return false
	}
	roles, exist := chatClient["roles"]
	if !exist {
		return false
	}
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
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
