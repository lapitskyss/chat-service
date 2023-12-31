// Code generated by options-gen. DO NOT EDIT.
package serverclient

import (
	fmt461e464ebed9 "fmt"

	"github.com/getkin/kin-openapi/openapi3"
	errors461e464ebed9 "github.com/kazhuravlev/options-gen/pkg/errors"
	validator461e464ebed9 "github.com/kazhuravlev/options-gen/pkg/validator"
	"github.com/labstack/echo/v4"
	"github.com/lapitskyss/chat-service/internal/middlewares"
	clientv1 "github.com/lapitskyss/chat-service/internal/server-client/v1"
	websocketstream "github.com/lapitskyss/chat-service/internal/websocket-stream"
	"go.uber.org/zap"
)

type OptOptionsSetter func(o *Options)

func NewOptions(
	logger *zap.Logger,
	addr string,
	allowOrigins []string,
	introspector middlewares.Introspector,
	requiredResource string,
	requiredRole string,
	v1Swagger *openapi3.T,
	v1Handlers clientv1.ServerInterface,
	wsHandler *websocketstream.HTTPHandler,
	httpErrorHandler echo.HTTPErrorHandler,
	shutdown func(),
	options ...OptOptionsSetter,
) Options {
	o := Options{}

	// Setting defaults from field tag (if present)

	o.logger = logger
	o.addr = addr
	o.allowOrigins = allowOrigins
	o.introspector = introspector
	o.requiredResource = requiredResource
	o.requiredRole = requiredRole
	o.v1Swagger = v1Swagger
	o.v1Handlers = v1Handlers
	o.wsHandler = wsHandler
	o.httpErrorHandler = httpErrorHandler
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
	errs.Add(errors461e464ebed9.NewValidationError("allowOrigins", _validate_Options_allowOrigins(o)))
	errs.Add(errors461e464ebed9.NewValidationError("introspector", _validate_Options_introspector(o)))
	errs.Add(errors461e464ebed9.NewValidationError("requiredResource", _validate_Options_requiredResource(o)))
	errs.Add(errors461e464ebed9.NewValidationError("requiredRole", _validate_Options_requiredRole(o)))
	errs.Add(errors461e464ebed9.NewValidationError("v1Swagger", _validate_Options_v1Swagger(o)))
	errs.Add(errors461e464ebed9.NewValidationError("v1Handlers", _validate_Options_v1Handlers(o)))
	errs.Add(errors461e464ebed9.NewValidationError("wsHandler", _validate_Options_wsHandler(o)))
	errs.Add(errors461e464ebed9.NewValidationError("httpErrorHandler", _validate_Options_httpErrorHandler(o)))
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

func _validate_Options_allowOrigins(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.allowOrigins, "min=1"); err != nil {
		return fmt461e464ebed9.Errorf("field `allowOrigins` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_introspector(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.introspector, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `introspector` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_requiredResource(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.requiredResource, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `requiredResource` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_requiredRole(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.requiredRole, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `requiredRole` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_v1Swagger(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.v1Swagger, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `v1Swagger` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_v1Handlers(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.v1Handlers, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `v1Handlers` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_wsHandler(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.wsHandler, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `wsHandler` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_httpErrorHandler(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.httpErrorHandler, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `httpErrorHandler` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_shutdown(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.shutdown, "-"); err != nil {
		return fmt461e464ebed9.Errorf("field `shutdown` did not pass the test: %w", err)
	}
	return nil
}
