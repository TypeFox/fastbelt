// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package generator

import (
	"fmt"
	"strings"

	"typefox.dev/fastbelt/generator"
	"typefox.dev/fastbelt/parser"
)

func EmitMarkdownSource(pkgName string, atn *ATN, tokenTypeNames []string) generator.Node {
	idx := make(map[*ATNState]int, len(atn.States))
	for i, s := range atn.States {
		idx[s] = i
	}

	n := generator.NewNode()

	n.AppendLine("# Runtime ATN for ", pkgName)
	n.AppendLine()

	// Group states by rule, preserving first-seen order.
	ruleOrder := []string{}
	ruleStates := map[string][]*ATNState{}
	for _, s := range atn.States {
		ruleName := (*s.Rule).Name()
		if _, seen := ruleStates[ruleName]; !seen {
			ruleOrder = append(ruleOrder, ruleName)
		}
		ruleStates[ruleName] = append(ruleStates[ruleName], s)
	}

	for _, ruleName := range ruleOrder {
		states := ruleStates[ruleName]
		n.AppendLine("## ", ruleName)
		n.AppendLine()
		n.AppendLine("```mermaid")
		n.AppendLine("flowchart TD")

		for _, s := range states {
			n.AppendLine("    ", mdNode(s, idx[s]))
		}

		n.AppendLine()

		for _, s := range states {
			src := fmt.Sprintf("%d", idx[s])
			for _, t := range s.Transitions {
				switch tr := t.(type) {
				case *EpsilonTransition:
					n.AppendLine("    q", src, " --> q", fmt.Sprintf("%d", idx[tr.Target()]))
				case *AtomTransition:
					tokenName := strings.ReplaceAll(tokenTypeNames[tr.TokenTypeId], "\"", "&quot;")
					n.AppendLine("    q", src, " -->|\"tok(", tokenName, ")\"| q", fmt.Sprintf("%d", idx[tr.Target()]))
				case *RuleTransition:
					n.AppendLine("    q", src, " -.->|\"[", (*tr.Target().Rule).Name(), "]\"| q", fmt.Sprintf("%d", idx[tr.FollowState]))
				}
			}
		}

		n.AppendLine("```")
		n.AppendLine()
	}

	return n
}

// mdNode returns the Mermaid node definition string for a state.
func mdNode(s *ATNState, i int) string {
	typShort := strings.TrimPrefix(atnStateTypeName(s.Type), "ATN")
	id := fmt.Sprintf("q%d", i)

	var label string
	if s.Type == parser.ATNRuleStart || s.Type == parser.ATNRuleStop {
		label = fmt.Sprintf("SN:%d<br/>%s", s.StateNumber, typShort)
	} else {
		decStr := ""
		if s.Decision >= 0 {
			decStr = fmt.Sprintf("<br/>dec=%d", s.Decision)
		}
		label = fmt.Sprintf("SN:%d<br/>%s<br/>%s", s.StateNumber, typShort, decStr)
	}

	switch s.Type {
	case parser.ATNRuleStart:
		return id + `(["` + label + `"])`
	case parser.ATNRuleStop:
		return id + `(["` + label + `"])`
	default:
		if s.Decision >= 0 {
			return id + `{"` + label + `"}`
		}
		return id + `["` + label + `"]`
	}
}
