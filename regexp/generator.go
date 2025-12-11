package regexp

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

var EOL = eol()

type Callback func(node Node)

type Node interface {
	Append(texts ...string) Node
	AppendLine(texts ...string) Node
	AppendNode(nodes ...Node) Node
	Indent(cb Callback) Node
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

func NewNode() Node {
	return newNode()
}
