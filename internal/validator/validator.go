package validator

import (
	"github.com/go-playground/validator/v10"
	optsGenValidator "github.com/kazhuravlev/options-gen/pkg/validator"
)

var Validator = validator.New()

func init() {
	optsGenValidator.Set(Validator)
}
