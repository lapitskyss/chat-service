// Code generated by options-gen. DO NOT EDIT.
package server

import (
	fmt461e464ebed9 "fmt"
	"net/http"

	errors461e464ebed9 "github.com/kazhuravlev/options-gen/pkg/errors"
	validator461e464ebed9 "github.com/kazhuravlev/options-gen/pkg/validator"
	"go.uber.org/zap"
)

type OptOptionsSetter func(o *Options)

func NewOptions(
	logger *zap.Logger,
	addr string,
	handler http.Handler,
	shutdown func(),
	options ...OptOptionsSetter,
) Options {
	o := Options{}

	// Setting defaults from field tag (if present)

	o.logger = logger
	o.addr = addr
	o.handler = handler
	o.shutdown = shutdown

	for _, opt := range options {
		opt(&o)
	}
	return o
}

func (o *Options) Validate() error {
	errs := new(errors461e464ebed9.ValidationErrors)
	errs.Add(errors461e464ebed9.NewValidationError("logger", _validate_Options_logger(o)))
	errs.Add(errors461e464ebed9.NewValidationError("addr", _validate_Options_addr(o)))
	errs.Add(errors461e464ebed9.NewValidationError("handler", _validate_Options_handler(o)))
	errs.Add(errors461e464ebed9.NewValidationError("shutdown", _validate_Options_shutdown(o)))
	return errs.AsError()
}

func _validate_Options_logger(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.logger, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `logger` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_addr(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.addr, "required,hostname_port"); err != nil {
		return fmt461e464ebed9.Errorf("field `addr` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_handler(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.handler, "-"); err != nil {
		return fmt461e464ebed9.Errorf("field `handler` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_shutdown(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.shutdown, "-"); err != nil {
		return fmt461e464ebed9.Errorf("field `shutdown` did not pass the test: %w", err)
	}
	return nil
}
