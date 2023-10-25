// Code generated by options-gen. DO NOT EDIT.
package serverdebug

import (
	fmt461e464ebed9 "fmt"

	"github.com/getkin/kin-openapi/openapi3"
	errors461e464ebed9 "github.com/kazhuravlev/options-gen/pkg/errors"
	validator461e464ebed9 "github.com/kazhuravlev/options-gen/pkg/validator"
)

type OptOptionsSetter func(o *Options)

func NewOptions(
	addr string,
	clientSwagger *openapi3.T,
	options ...OptOptionsSetter,
) Options {
	o := Options{}

	// Setting defaults from field tag (if present)

	o.addr = addr
	o.clientSwagger = clientSwagger

	for _, opt := range options {
		opt(&o)
	}
	return o
}

func (o *Options) Validate() error {
	errs := new(errors461e464ebed9.ValidationErrors)
	errs.Add(errors461e464ebed9.NewValidationError("addr", _validate_Options_addr(o)))
	errs.Add(errors461e464ebed9.NewValidationError("clientSwagger", _validate_Options_clientSwagger(o)))
	return errs.AsError()
}

func _validate_Options_addr(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.addr, "required,hostname_port"); err != nil {
		return fmt461e464ebed9.Errorf("field `addr` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_clientSwagger(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.clientSwagger, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `clientSwagger` did not pass the test: %w", err)
	}
	return nil
}
