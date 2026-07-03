package fastbelt

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unique"
)

type FragmentPath interface {
	Empty() bool
	Push(name unique.Handle[string], index int) FragmentPath
	Shift() (unique.Handle[string], int)
	String() string
}

type FragmentPathImpl struct {
	segments []struct {
		name  unique.Handle[string]
		index int
	}
}

func (p *FragmentPathImpl) Empty() bool {
	return len(p.segments) == 0
}

func (p *FragmentPathImpl) Push(name unique.Handle[string], index int) FragmentPath {
	p.segments = append(p.segments, struct {
		name  unique.Handle[string]
		index int
	}{name, index})
	return p
}

func (p *FragmentPathImpl) Shift() (name unique.Handle[string], index int) {
	len := len(p.segments)
	if len == 0 {
		return fieldZero, -1
	} else {
		head := p.segments[0]
		p.segments = p.segments[1:]
		return head.name, head.index
	}
}

func (ps *FragmentPathImpl) String() string {
	var sb strings.Builder
	for _, s := range ps.segments {
		sb.WriteString("/")
		sb.WriteString(s.name.Value())
		if s.index >= 0 {
			sb.WriteString("@")
			sb.WriteString(strconv.Itoa(int(s.index)))
		}
	}
	return sb.String()
}

var fieldZero = unique.Handle[string]{}

// NodePath determines node's path within its root container in a recursive manner,
// based on node's [AstNode.ContainmentData]. It returns a slash-separated
// path string that uniquely identifies this node within its document tree,
// e.g. "/rules@2/alternatives@0".
// Returns "" for the root node (no container).
func PathOf(node AstNode) (FragmentPath, error) {
	container := node.Container()
	containerField, index := node.ContainmentData()

	b := &FragmentPathImpl{}
	if container == nil {
		return b, nil
	} else if containerField == fieldZero || containerField.Value() == "" {
		return b, errors.New("cannot determine node path, 'containerField' is empty")
	}

	parentPath, err := PathOf(container)
	if err != nil {
		if errors.Unwrap(err) == nil {
			return b, fmt.Errorf(
				"AstNodeBase.NodePath: error within node of type %T: %w", container, err)
		} else {
			return b, fmt.Errorf(
				"AstNodeBase.NodePath: error within container of type %T:\n %w", container, err)
		}
	}
	return parentPath.Push(containerField, index), nil
}

func Resolve(path string, node AstNode) (AstNode, error) {
	if path == "" {
		return node, nil
	}
	segments, err := parsePath(path)
	if err != nil {
		return nil, err
	}
	return node.Resolve(segments)
}

func parsePath(path string) (FragmentPath, error) {
	result := &FragmentPathImpl{}
	path = strings.TrimLeft(path, "/")

	for segment := range strings.SplitSeq(path, "/") {
		if segment == "" {
			continue
		}
		fieldAndIndex := strings.SplitN(segment, "@", 2)
		field := unique.Make(fieldAndIndex[0])
		if len(fieldAndIndex) == 1 {
			result.Push(field, -1)
		} else {
			index, err := strconv.Atoi(fieldAndIndex[1])
			if err != nil {
				return nil, fmt.Errorf("ParsePath: index '%s' is not a valid uint: %w", fieldAndIndex[1], err)
			}
			result.Push(field, index)
		}
	}
	return result, nil
}
