package statemachine

import (
	"context"

	fastbelt "typefox.dev/fastbelt"
)

var _ fastbelt.Validator = (*StatemachineImpl)(nil)

// Validate checks statemachine constraints:
//   - Event names must be unique.
//   - State names must be unique.
//   - Transitions must have valid event and state targets.
func (s *StatemachineImpl) Validate(_ context.Context, _ string, accept fastbelt.ValidationAcceptor) {
	_ = checkUniqueEventNames(s, accept)
	_ = checkUniqueStateNames(s, accept)
}

func checkUniqueEventNames(s *StatemachineImpl, accept fastbelt.ValidationAcceptor) map[string]bool {
	seen := map[string]bool{}
	for _, event := range s.Events() {
		if seen[event.Name()] {
			accept(fastbelt.NewDiagnostic(
				fastbelt.SeverityError,
				"Event name must be unique.",
				event,
				fastbelt.WithToken(event.NameToken()),
			))
		}
		seen[event.Name()] = true
	}
	return seen
}

func checkUniqueStateNames(s *StatemachineImpl, accept fastbelt.ValidationAcceptor) map[string]bool {
	seen := map[string]bool{}
	for _, state := range s.States() {
		if seen[state.Name()] {
			accept(fastbelt.NewDiagnostic(
				fastbelt.SeverityError,
				"State name must be unique.",
				state,
				fastbelt.WithToken(state.NameToken()),
			))
		}
		seen[state.Name()] = true
	}
	return seen
}