package common

import "context"

type TransactionStep struct {
	Name       string
	Execute    func(ctx context.Context) error
	Compensate func(ctx context.Context) error
}

type Transaction struct {
	steps         []TransactionStep
	executedSteps []TransactionStep
}

func NewTransaction() *Transaction {
	return &Transaction{
		steps:         make([]TransactionStep, 0),
		executedSteps: make([]TransactionStep, 0),
	}
}

func (t *Transaction) Add(step TransactionStep) *Transaction {
	t.steps = append(t.steps, step)
	return t
}

func (t *Transaction) Run(ctx context.Context, name string, execute func(context.Context) error, compensate func(context.Context) error) error {
	log := NewAppLogger(ctx, "Transaction.Run")
	step := TransactionStep{Name: name, Execute: execute, Compensate: compensate}

	log.Info().Str("step", step.Name).Msg("executing step")

	if err := step.Execute(ctx); err != nil {
		log.Error().Err(err).Str("step", step.Name).Msg("step failed, initiating rollback")
		t.rollback(ctx)
		return err
	}

	t.executedSteps = append(t.executedSteps, step)
	log.Info().Str("step", step.Name).Msg("step completed")
	return nil
}

func (t *Transaction) Execute(ctx context.Context) error {
	log := NewAppLogger(ctx, "Transaction.Execute")
	log.Info().Int("total_steps", len(t.steps)).Msg("starting transaction")

	for _, step := range t.steps {
		log.Info().Str("step", step.Name).Msg("executing step")

		if err := step.Execute(ctx); err != nil {
			log.Error().Err(err).Str("step", step.Name).Msg("step failed, initiating rollback")
			t.rollback(ctx)
			return err
		}

		t.executedSteps = append(t.executedSteps, step)
		log.Info().Str("step", step.Name).Msg("step completed")
	}

	log.Info().Msg("all steps completed")
	return nil
}

func (t *Transaction) Commit(ctx context.Context) {
	log := NewAppLogger(ctx, "Transaction.Commit")
	log.Info().Int("steps_completed", len(t.executedSteps)).Msg("transaction committed")
	t.executedSteps = nil
}

func (t *Transaction) rollback(ctx context.Context) {
	log := NewAppLogger(ctx, "Transaction.Rollback")
	log.Warn().Int("steps_to_compensate", len(t.executedSteps)).Msg("starting rollback")

	for i := len(t.executedSteps) - 1; i >= 0; i-- {
		step := t.executedSteps[i]
		log.Warn().Str("step", step.Name).Msg("compensating step")

		if err := step.Compensate(ctx); err != nil {
			log.Error().Err(err).Str("step", step.Name).Msg("compensation failed — system may be inconsistent")
		} else {
			log.Warn().Str("step", step.Name).Msg("step compensated")
		}
	}
}
