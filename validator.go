package structure_validator

import "context"

// RuleValidator defines an interface for validating rules against a given type T and returning a list of RuleErrors.
type RuleValidator[T any] interface {

	// Analyze validates the input value against defined rules and returns a slice of RuleError structs for any violations.
	// It respects the timeout defined in the EngineConfig and includes a critical error if the timeout is exceeded.
	// The method runs rules concurrently, collecting errors in a result channel to preserve responsiveness.
	Analyze(T) []*RuleError
}

type ruleValidator[T any] struct {
	anz RuleEngine[T]
}

func (r *ruleValidator[T]) Analyze(value T) []*RuleError {
	var errors []*RuleError

	cfg := r.anz.Config()
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	resultChan := make(chan *RuleError)

	go func() {
		for _, rule := range r.anz.Rules() {
			select {
			case <-ctx.Done():
				return
			default:
				if result := rule.Validate(value); result != nil {
					resultChan <- result
				}
			}
		}
		close(resultChan)
	}()

	for {
		select {
		case <-ctx.Done():
			errors = append(errors, &RuleError{
				Message:    "TIMEOUT_EXCEEDED after: " + cfg.Timeout.String(),
				IsCritical: true,
			})
			return errors
		case result, ok := <-resultChan:
			if !ok {
				return errors
			}
			if result != nil {
				errors = append(errors, result)
			}

		}

	}

}

// NewRuleValidator creates a RuleValidator for the given RuleEngine, enabling validation of rules against type T.
func NewRuleValidator[T any](anz RuleEngine[T]) RuleValidator[T] {
	return &ruleValidator[T]{
		anz: anz,
	}
}
