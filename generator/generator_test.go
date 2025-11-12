package generator

import "testing"

func TestGeneratorIndent(t *testing.T) {
	root := NewNode()
	root.AppendLine("line 1")
	root.Indent(func(n Node) {
		n.AppendLine("indented line 1")
		n.AppendLine("indented line 2")
	})
	root.AppendLine("line 2")
	expected := "line 1" + EOL +
		"    indented line 1" + EOL +
		"    indented line 2" + EOL +
		"line 2" + EOL
	actual := root.String()
	if actual != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, actual)
	}
}
func TestGeneratorAppend(t *testing.T) {
	root := NewNode()
	root.Append("Hello", " ", "World")
	expected := "Hello World"
	actual := root.String()
	if actual != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, actual)
	}
}

func TestGeneratorAppendLine(t *testing.T) {
	root := NewNode()
	root.AppendLine("Line 1")
	root.AppendLine("Line 2")
	expected := "Line 1" + EOL + "Line 2" + EOL
	actual := root.String()
	if actual != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, actual)
	}
}

func TestGeneratorAppendNode(t *testing.T) {
	root := NewNode()
	child := NewNode()
	child.Append("Child content")
	root.Append("Parent ")
	root.AppendNode(child)
	expected := "Parent Child content"
	actual := root.String()
	if actual != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, actual)
	}
}

func TestGeneratorNestedIndent(t *testing.T) {
	root := NewNode()
	root.AppendLine("root")
	root.Indent(func(n Node) {
		n.AppendLine("level 1")
		n.Indent(func(n2 Node) {
			n2.AppendLine("level 2")
		})
	})
	expected := "root" + EOL +
		"    level 1" + EOL +
		"        level 2" + EOL
	actual := root.String()
	if actual != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, actual)
	}
}

func TestGeneratorEmptyNode(t *testing.T) {
	root := NewNode()
	expected := ""
	actual := root.String()
	if actual != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, actual)
	}
}

func TestGeneratorMultipleAppends(t *testing.T) {
	root := NewNode()
	root.Append("A").Append("B").AppendLine("C")
	root.Append("D").AppendLine("E")
	expected := "ABC" + EOL + "DE" + EOL
	actual := root.String()
	if actual != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, actual)
	}
}
