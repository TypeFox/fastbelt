// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package allstar

import (
	"fmt"
	"strings"

	"typefox.dev/fastbelt/generator"
)

func EmitMarkdownSource(pkgName string, rtn *RuntimeATN, tokenTypeNames map[int]string) generator.Node {
	idx := make(map[*RuntimeATNState]int, len(rtn.States))
	for i, s := range rtn.States {
		idx[s] = i
	}

	n := generator.NewNode()

	n.AppendLine("# Runtime ATN for ", pkgName)
	n.AppendLine()

	// Group states by rule, preserving first-seen order.
	ruleOrder := []string{}
	ruleStates := map[string][]*RuntimeATNState{}
	for _, s := range rtn.States {
		if _, seen := ruleStates[s.RuleName]; !seen {
			ruleOrder = append(ruleOrder, s.RuleName)
		}
		ruleStates[s.RuleName] = append(ruleStates[s.RuleName], s)
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
				case *RuntimeEpsilonTransition:
					n.AppendLine("    q", src, " --> q", fmt.Sprintf("%d", idx[tr.Target]))
				case *RuntimeAtomTransition:
					n.AppendLine("    q", src, " -->|\"tok(", tokenTypeNames[tr.TokenTypeID], ")\"| q", fmt.Sprintf("%d", idx[tr.Target]))
				case *RuntimeRuleTransition:
					n.AppendLine("    q", src, " -.->|\"[", tr.Target.RuleName, "]\"| q", fmt.Sprintf("%d", idx[tr.FollowState]))
				}
			}
		}

		n.AppendLine("```")
		n.AppendLine()
	}

	return n
}

// mdNode returns the Mermaid node definition string for a state.
func mdNode(s *RuntimeATNState, i int) string {
	typShort := strings.TrimPrefix(atnStateTypeName(s.Type), "ATN")
	id := fmt.Sprintf("q%d", i)

	var label string
	if s.Type == ATNRuleStart || s.Type == ATNRuleStop {
		label = fmt.Sprintf("SN:%d<br/>%s", s.StateNumber, typShort)
	} else {
		prodShort := strings.TrimPrefix(productionKindName(s.ProdKind), "Prod")
		decStr := ""
		if s.Decision >= 0 {
			decStr = fmt.Sprintf("<br/>dec=%d", s.Decision)
		}
		label = fmt.Sprintf("SN:%d<br/>%s<br/>%s #%d%s", s.StateNumber, typShort, prodShort, s.ProdIdx, decStr)
	}

	switch s.Type {
	case ATNRuleStart:
		return id + `(["` + label + `"])`
	case ATNRuleStop:
		return id + `(["` + label + `"])`
	default:
		if s.Decision >= 0 {
			return id + `{"` + label + `"}`
		}
		return id + `["` + label + `"]`
	}
}
