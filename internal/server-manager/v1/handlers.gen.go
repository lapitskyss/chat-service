// Code generated by options-gen. DO NOT EDIT.
package managerv1

import (
	fmt461e464ebed9 "fmt"

	errors461e464ebed9 "github.com/kazhuravlev/options-gen/pkg/errors"
	validator461e464ebed9 "github.com/kazhuravlev/options-gen/pkg/validator"
)

type OptOptionsSetter func(o *Options)

func NewOptions(
	canReceiveProblems canReceiveProblemsUseCase,
	options ...OptOptionsSetter,
) Options {
	o := Options{}

	// Setting defaults from field tag (if present)

	o.canReceiveProblems = canReceiveProblems

	for _, opt := range options {
		opt(&o)
	}
	return o
}

func (o *Options) Validate() error {
	errs := new(errors461e464ebed9.ValidationErrors)
	errs.Add(errors461e464ebed9.NewValidationError("canReceiveProblems", _validate_Options_canReceiveProblems(o)))
	return errs.AsError()
}

func _validate_Options_canReceiveProblems(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.canReceiveProblems, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `canReceiveProblems` did not pass the test: %w", err)
	}
	return nil
}
