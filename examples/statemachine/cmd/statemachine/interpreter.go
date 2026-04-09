// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strings"

	"typefox.dev/fastbelt/examples/statemachine"
)

// PrintModelSummary writes a human-readable sketch of the state machine (states, transitions, declared commands)
// to Stdout using the parsed model on the runner.
func (r *Runner) PrintModelSummary() error {
	if r.Stdout == nil {
		return fmt.Errorf("runner: Stdout must be non-nil")
	}
	if r.sm == nil {
		return fmt.Errorf("runner: no parsed model; call ParseAndValidate first")
	}
	sm := r.sm
	w := r.Stdout

	_, _ = fmt.Fprintf(w, "State machine: %q\n", sm.Name())
	_, _ = fmt.Fprintf(w, "Events: %s\n", joinNames(eventNames(sm)))
	_, _ = fmt.Fprintf(w, "Commands: %s\n", joinNames(commandNames(sm)))
	initRef := sm.Init()
	if initRef == nil {
		_, _ = fmt.Fprintf(w, "Initial state: <missing>\n")
	} else {
		_, _ = fmt.Fprintf(w, "Initial state: %q\n", initRef.Text())
	}
	_, _ = fmt.Fprintf(w, "States:\n")
	for _, st := range sm.States() {
		var actionParts []string
		for _, a := range st.Actions() {
			if a == nil {
				continue
			}
			actionParts = append(actionParts, a.Text())
		}
		actions := "(no actions block)"
		if len(actionParts) > 0 {
			actions = strings.Join(actionParts, ", ")
		}
		_, _ = fmt.Fprintf(w, "  - %q  actions: %s\n", st.Name(), actions)
		for _, tr := range st.Transitions() {
			ev := ""
			if tr.Event() != nil {
				ev = tr.Event().Text()
			}
			to := ""
			if tr.State() != nil {
				to = tr.State().Text()
			}
			_, _ = fmt.Fprintf(w, "      %s => %s\n", ev, to)
		}
	}
	return nil
}

// Interpret runs a tiny event loop: each non-empty line from EventInput is an event name; the current state's
// transitions are searched and the machine moves when a matching transition exists. After each successful step,
// commands linked from the new state's actions block are reported on Stdout (demo "runtime" only).
func (r *Runner) Interpret() error {
	if r.Stdout == nil {
		return fmt.Errorf("runner: Stdout must be non-nil")
	}
	if r.EventInput == nil {
		return fmt.Errorf("runner: EventInput must be non-nil")
	}
	if r.sm == nil {
		return fmt.Errorf("runner: no parsed model; call ParseAndValidate first")
	}

	ctx := context.Background()
	sm := r.sm
	init := sm.Init()
	if init == nil {
		return fmt.Errorf("no initialState in document")
	}
	current := init.Ref(ctx)
	if current == nil {
		return fmt.Errorf("initial state reference did not resolve (see linker diagnostics)")
	}
	_, _ = fmt.Fprintf(r.Stdout, "Start state: %q\n", current.Name())

	sc := bufio.NewScanner(r.EventInput)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		next, stepped := stepTransition(ctx, current, line)
		if !stepped {
			_, _ = fmt.Fprintf(r.Stdout, "  no transition on %q from %q\n", line, current.Name())
			continue
		}
		current = next
		_, _ = fmt.Fprintf(r.Stdout, "  %q -> %q\n", line, current.Name())
		emitCommandsForState(r.Stdout, current)
	}
	return sc.Err()
}

func eventNames(sm statemachine.Statemachine) []string {
	var names []string
	for _, e := range sm.Events() {
		names = append(names, e.Name())
	}
	return names
}

func commandNames(sm statemachine.Statemachine) []string {
	var names []string
	for _, c := range sm.Commands() {
		names = append(names, c.Name())
	}
	return names
}

func joinNames(names []string) string {
	if len(names) == 0 {
		return "(none)"
	}
	return strings.Join(names, ", ")
}

func emitCommandsForState(w io.Writer, st statemachine.State) {
	for _, a := range st.Actions() {
		if a == nil {
			continue
		}
		_, _ = fmt.Fprintf(w, "  emit command %q\n", a.Text())
	}
}

func stepTransition(ctx context.Context, st statemachine.State, event string) (statemachine.State, bool) {
	for _, tr := range st.Transitions() {
		ev := tr.Event()
		if ev == nil || ev.Text() != event {
			continue
		}
		target := tr.State()
		if target == nil {
			return nil, false
		}
		next := target.Ref(ctx)
		if next == nil {
			return nil, false
		}
		return next, true
	}
	return nil, false
}
