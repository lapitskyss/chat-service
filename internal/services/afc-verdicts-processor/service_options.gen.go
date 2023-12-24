// Code generated by options-gen. DO NOT EDIT.
package afcverdictsprocessor

import (
	fmt461e464ebed9 "fmt"
	"time"

	errors461e464ebed9 "github.com/kazhuravlev/options-gen/pkg/errors"
	validator461e464ebed9 "github.com/kazhuravlev/options-gen/pkg/validator"
)

type OptOptionsSetter func(o *Options)

func NewOptions(
	brokers []string,
	consumers int,
	consumerGroup string,
	verdictsTopic string,
	readerFactory KafkaReaderFactory,
	dlqWriter KafkaDLQWriter,
	txtor transactor,
	msgRepo messagesRepository,
	outBox outboxService,
	options ...OptOptionsSetter,
) Options {
	o := Options{}

	// Setting defaults from field tag (if present)
	o.backoffInitialInterval, _ = time.ParseDuration("100ms")
	o.backoffMaxElapsedTime, _ = time.ParseDuration("5s")
	o.processBatchSize = 1

	o.brokers = brokers
	o.consumers = consumers
	o.consumerGroup = consumerGroup
	o.verdictsTopic = verdictsTopic
	o.readerFactory = readerFactory
	o.dlqWriter = dlqWriter
	o.txtor = txtor
	o.msgRepo = msgRepo
	o.outBox = outBox

	for _, opt := range options {
		opt(&o)
	}
	return o
}

func WithBackoffInitialInterval(opt time.Duration) OptOptionsSetter {
	return func(o *Options) {
		o.backoffInitialInterval = opt
	}
}

func WithBackoffMaxElapsedTime(opt time.Duration) OptOptionsSetter {
	return func(o *Options) {
		o.backoffMaxElapsedTime = opt
	}
}

func WithVerdictsSignKey(opt string) OptOptionsSetter {
	return func(o *Options) {
		o.verdictsSignKey = opt
	}
}

func WithProcessBatchSize(opt int) OptOptionsSetter {
	return func(o *Options) {
		o.processBatchSize = opt
	}
}

func (o *Options) Validate() error {
	errs := new(errors461e464ebed9.ValidationErrors)
	errs.Add(errors461e464ebed9.NewValidationError("backoffInitialInterval", _validate_Options_backoffInitialInterval(o)))
	errs.Add(errors461e464ebed9.NewValidationError("backoffMaxElapsedTime", _validate_Options_backoffMaxElapsedTime(o)))
	errs.Add(errors461e464ebed9.NewValidationError("brokers", _validate_Options_brokers(o)))
	errs.Add(errors461e464ebed9.NewValidationError("consumers", _validate_Options_consumers(o)))
	errs.Add(errors461e464ebed9.NewValidationError("consumerGroup", _validate_Options_consumerGroup(o)))
	errs.Add(errors461e464ebed9.NewValidationError("verdictsTopic", _validate_Options_verdictsTopic(o)))
	errs.Add(errors461e464ebed9.NewValidationError("processBatchSize", _validate_Options_processBatchSize(o)))
	errs.Add(errors461e464ebed9.NewValidationError("readerFactory", _validate_Options_readerFactory(o)))
	errs.Add(errors461e464ebed9.NewValidationError("dlqWriter", _validate_Options_dlqWriter(o)))
	errs.Add(errors461e464ebed9.NewValidationError("txtor", _validate_Options_txtor(o)))
	errs.Add(errors461e464ebed9.NewValidationError("msgRepo", _validate_Options_msgRepo(o)))
	errs.Add(errors461e464ebed9.NewValidationError("outBox", _validate_Options_outBox(o)))
	return errs.AsError()
}

func _validate_Options_backoffInitialInterval(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.backoffInitialInterval, "min=50ms,max=1s"); err != nil {
		return fmt461e464ebed9.Errorf("field `backoffInitialInterval` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_backoffMaxElapsedTime(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.backoffMaxElapsedTime, "min=500ms,max=1m"); err != nil {
		return fmt461e464ebed9.Errorf("field `backoffMaxElapsedTime` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_brokers(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.brokers, "min=1"); err != nil {
		return fmt461e464ebed9.Errorf("field `brokers` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_consumers(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.consumers, "min=1,max=16"); err != nil {
		return fmt461e464ebed9.Errorf("field `consumers` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_consumerGroup(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.consumerGroup, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `consumerGroup` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_verdictsTopic(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.verdictsTopic, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `verdictsTopic` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_processBatchSize(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.processBatchSize, "min=1"); err != nil {
		return fmt461e464ebed9.Errorf("field `processBatchSize` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_readerFactory(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.readerFactory, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `readerFactory` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_dlqWriter(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.dlqWriter, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `dlqWriter` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_txtor(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.txtor, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `txtor` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_msgRepo(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.msgRepo, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `msgRepo` did not pass the test: %w", err)
	}
	return nil
}

func _validate_Options_outBox(o *Options) error {
	if err := validator461e464ebed9.GetValidatorFor(o).Var(o.outBox, "required"); err != nil {
		return fmt461e464ebed9.Errorf("field `outBox` did not pass the test: %w", err)
	}
	return nil
}
