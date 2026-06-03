// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package codegen

import (
	"runtime"
	"strings"
)

func eol() string {
	if runtime.GOOS == "windows" {
		return "\r\n"
	}
	return "\n"
}

// EOL is the line separator used when joining logical lines in [Node.String].
// It is "\r\n" on Windows hosts and "\n" elsewhere, fixed at package init time.
var EOL = eol()

// Callback receives a [Node] to populate, typically from [Node.Indent].
type Callback func(node Node)

// Node is a mutable buffer for building indented, multi-line text. Generators
// assemble source by appending to the current line, starting new lines, nesting
// blocks with [Node.Indent], and merging completed subtrees with [Node.AppendNode].
// Methods return the receiver so calls can be chained.
type Node interface {
	// Append adds texts to the current line without starting a new line. When the
	// current line is empty, leading indentation for the node's depth is inserted
	// before the first text.
	Append(texts ...string) Node
	// AppendLine adds texts to the current line and then begins a new line.
	AppendLine(texts ...string) Node
	// AppendNode merges nodes into the receiver. Each child's first line is
	// concatenated onto the receiver's current line; any further lines from the
	// child become subsequent lines of the receiver.
	AppendNode(nodes ...Node) Node
	// Indent runs cb with a child node indented one level deeper (four spaces
	// per level), then splices the child's content into the receiver. If the
	// receiver's current line is non-empty, a line break is inserted before the
	// merged content.
	Indent(cb Callback) Node
	// String returns the accumulated text with [EOL] between logical lines.
	String() string
}

type node struct {
	lines  []*strings.Builder
	indent int32
}

func (n *node) Append(texts ...string) Node {
	currentLine := n.lines[len(n.lines)-1]
	if len(texts) > 0 {
		// If the current line is empty, add indentation
		if currentLine.Len() == 0 {
			indentString(currentLine, n.indent)
		}
		for _, text := range texts {
			currentLine.WriteString(text)
		}
	}

	return n
}

func (n *node) AppendLine(texts ...string) Node {
	n.Append(texts...)
	n.lines = append(n.lines, &strings.Builder{})
	return n
}

func (n *node) Indent(cb Callback) Node {
	child := newNode()
	child.indent = n.indent + 1
	cb(child)
	lastLine := n.lines[len(n.lines)-1]
	if lastLine.Len() > 0 {
		n.AppendLine()
	}
	n.AppendNode(child)
	return n
}

func (n *node) AppendNode(nodes ...Node) Node {
	for _, child := range nodes {
		if genNode, ok := child.(*node); ok {
			// If we can directly access the lines, we can optimize the appending
			// We don't even need to generate the string representation
			firstLine := genNode.lines[0]
			lastLine := n.lines[len(n.lines)-1]
			// Append first line to current last line
			lastLine.WriteString(firstLine.String())
			// Append remaining lines
			n.lines = append(n.lines, genNode.lines[1:]...)
		} else {
			// Split the string representation into lines and append
			lines := splitStringLines(child.String())
			lastLine := n.lines[len(n.lines)-1]
			if len(lines) == 0 {
				continue
			}
			// Append first line to current last line
			lastLine.WriteString(lines[0])
			// Append remaining lines
			for _, line := range lines[1:] {
				sb := &strings.Builder{}
				sb.WriteString(line)
				n.lines = append(n.lines, sb)
			}
		}
	}
	return n
}

func (n *node) String() string {
	sb := &strings.Builder{}
	for i, line := range n.lines {
		sb.WriteString(line.String())
		if i < len(n.lines)-1 {
			sb.WriteString(EOL)
		}
	}
	return sb.String()
}

func splitStringLines(s string) []string {
	return strings.Split(s, EOL)
}

func indentString(currentLine *strings.Builder, indent int32) {
	for range indent {
		currentLine.WriteString("    ")
	}
}

func newNode() *node {
	return &node{
		// Initialize with one empty line
		lines:  []*strings.Builder{{}},
		indent: 0,
	}
}

// NewNode returns an empty [Node] with one blank line ready for appending.
func NewNode() Node {
	return newNode()
}
