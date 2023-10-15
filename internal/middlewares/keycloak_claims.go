package middlewares

import (
	"errors"

	"github.com/golang-jwt/jwt"

	"github.com/lapitskyss/chat-service/internal/types"
)

var (
	ErrNoAllowedResources = errors.New("no allowed resources")
	ErrSubjectNotDefined  = errors.New(`"sub" is not defined`)
)

type claims struct {
	jwt.StandardClaims
	Audience       any                            `json:"aud,omitempty"`
	Subject        types.UserID                   `json:"sub"`
	ResourceAccess map[string]map[string][]string `json:"resource_access"`
}

// Valid returns errors:
// - from StandardClaims validation;
// - ErrNoAllowedResources, if claims doesn't contain `resource_access` map or it's empty;
// - ErrSubjectNotDefined, if claims doesn't contain `sub` field or subject is zero UUID.
func (c claims) Valid() error {
	if err := c.StandardClaims.Valid(); err != nil {
		return err
	}
	if len(c.ResourceAccess) == 0 {
		return ErrNoAllowedResources
	}
	if c.Subject.IsZero() {
		return ErrSubjectNotDefined
	}
	return nil
}

func (c claims) UserID() types.UserID {
	return c.Subject
}
