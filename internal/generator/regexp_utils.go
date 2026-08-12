// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package generator

import (
	"regexp"
	"regexp/syntax"
)

const MaxRune = 0x10FFFF

type RegexpProgram struct {
	Instructions []VmInst
}

type VmOp int

const (
	VmOpChar VmOp = iota
	VmOpJump
	VmOpSplit
	VmOpMatch
	VmOpReset
)

type VmInst struct {
	Op      VmOp
	Rune    rune
	Targets []int
}

const NL = '\n'

func escapeRegExp(value string) string {
	var re = regexp.MustCompile(`[.*+?^${}()|[\]\\]`)
	return re.ReplaceAllStringFunc(value, func(match string) string {
		return `\` + match
	})
}

func visitRegexp(program *RegexpProgram, node *syntax.Regexp) {
	switch node.Op {
	case syntax.OpAlternate:
		ends := []int{}
		program.Instructions = append(program.Instructions, VmInst{
			Op:      VmOpSplit,
			Targets: []int{},
		})
		splitIndex := len(program.Instructions) - 1
		for _, sub := range node.Sub {
			program.Instructions[splitIndex].Targets = append(program.Instructions[splitIndex].Targets, len(program.Instructions))
			visitRegexp(program, sub)
			program.Instructions = append(program.Instructions, VmInst{
				Op:      VmOpJump,
				Targets: []int{},
			})
			jumpToEndIndex := len(program.Instructions) - 1
			ends = append(ends, jumpToEndIndex)
		}
		for _, end := range ends {
			program.Instructions[end].Targets = append(program.Instructions[end].Targets, len(program.Instructions))
		}
	case syntax.OpConcat:
		for _, sub := range node.Sub {
			visitRegexp(program, sub)
		}
	case syntax.OpLiteral:
		for _, r := range node.Rune {
			program.Instructions = append(program.Instructions, VmInst{
				Op:   VmOpChar,
				Rune: r,
			})
		}
	case syntax.OpQuest:
		splitInst := VmInst{
			Op:      VmOpSplit,
			Targets: []int{len(program.Instructions) + 1},
		}
		program.Instructions = append(program.Instructions, splitInst)
		visitRegexp(program, node.Sub[0])
		splitInst.Targets = append(splitInst.Targets, len(program.Instructions))
	case syntax.OpPlus, syntax.OpStar, syntax.OpRepeat, syntax.OpAnyChar, syntax.OpAnyCharNotNL, syntax.OpCharClass:
		program.Instructions = append(program.Instructions, VmInst{
			Op: VmOpReset,
		})
	case syntax.OpCapture, syntax.OpBeginLine, syntax.OpEndLine, syntax.OpBeginText, syntax.OpEndText, syntax.OpWordBoundary, syntax.OpNoWordBoundary:
		visitRegexp(program, node.Sub[0])
	case syntax.OpEmptyMatch, syntax.OpNoMatch:
	}
}

type Part struct{ start, end string }

func GetTerminalParts(regexp string) []Part {
	pattern, err := syntax.Parse(regexp, syntax.Perl)
	if err != nil {
		return nil
	}
	program := &RegexpProgram{}
	visitRegexp(program, pattern)
	program.Instructions = append(program.Instructions, VmInst{
		Op: VmOpMatch,
	})
	return execute(program)
}

type thread struct {
	pc          int
	isHalted    bool
	isStarting  bool
	isMultiline bool
	startRegexp string
	endRegexp   string
}

func execute(program *RegexpProgram) []Part {
	result := []Part{}
	threads := []thread{{
		pc:          0,
		isHalted:    false,
		isStarting:  true,
		isMultiline: false,
		startRegexp: "",
		endRegexp:   "",
	}}
	for len(threads) > 0 {
		th := threads[len(threads)-1]
		threads = threads[:len(threads)-1]
		if th.isHalted {
			result = append(result, Part{
				start: th.startRegexp,
				end:   th.endRegexp,
			})
			continue
		}
		if th.pc >= len(program.Instructions) {
			continue
		}
		inst := program.Instructions[th.pc]
		switch inst.Op {
		case VmOpChar:
			if inst.Rune == NL {
				th.isMultiline = true
			}
			char := escapeRegExp(string(inst.Rune))
			if th.isStarting {
				th.startRegexp += char
			} else {
				th.endRegexp += char
			}
			th.pc++
			threads = append(threads, th)
		case VmOpJump:
			th.pc = inst.Targets[0]
			threads = append(threads, th)
		case VmOpSplit:
			for _, target := range inst.Targets {
				newThread := thread{
					pc:          target,
					isStarting:  th.isStarting,
					isMultiline: th.isMultiline,
					startRegexp: th.startRegexp,
					endRegexp:   th.endRegexp,
				}
				threads = append(threads, newThread)
			}
		case VmOpMatch:
			th.isHalted = true
			threads = append(threads, th)
		case VmOpReset:
			th.isStarting = false
			th.endRegexp = ""
			th.pc++
			threads = append(threads, th)
		}
	}
	return result
}
