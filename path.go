// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package fastbelt

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unique"
)

// FragmentPath describes instance of node path descriptors pointing to a child AST node within a parent.
type FragmentPath interface {
	// Empty returns [true] if this descriptor is empty and equal to the empty path string, returns [false] otherwise.
	Empty() bool
	// Head returns field name (unique.Handle) and index of the head segment of the segment list.
	// index is expected to be negative (-1) for non-list fields.
	Head() (unique.Handle[string], int)
	// Tail returns a copy of this path descriptor with the first element dropped from the segments list.
	Tail() FragmentPath
	// String returns the corresponding fragment path string
	String() string
}

// Base implementation of [FragmentPath]
type fragmentPathImpl struct {
	segments []struct {
		name  unique.Handle[string]
		index int
	}
}

var _ FragmentPath = (*fragmentPathImpl)(nil)

func (p *fragmentPathImpl) push(name unique.Handle[string], index int) *fragmentPathImpl {
	p.segments = append(p.segments, struct {
		name  unique.Handle[string]
		index int
	}{name, index})
	return p
}

func (p *fragmentPathImpl) Empty() bool {
	return len(p.segments) == 0
}

func (p *fragmentPathImpl) Head() (name unique.Handle[string], index int) {
	count := len(p.segments)
	if count == 0 {
		return fieldZero, -1
	} else {
		head := p.segments[0]
		return head.name, head.index
	}
}

func (p *fragmentPathImpl) Tail() FragmentPath {
	clone := *p
	clone.segments = clone.segments[1:]
	return &clone
}

func (p *fragmentPathImpl) String() string {
	var sb strings.Builder
	for _, s := range p.segments {
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

// PathOf determines node's path within its root container in a recursive manner
// based on node's [AstNode.ContainmentData]. It returns a slash-separated
// path string that uniquely identifies this node within its document tree,
// e.g. "/rules@2/alternatives@0".
// Returns "" for the root node (no container).
func PathOf(node AstNode) (*fragmentPathImpl, error) {
	container := node.Container()
	containerField, index := node.ContainmentData()

	if container == nil {
		return &fragmentPathImpl{}, nil
	} else if containerField == fieldZero || containerField.Value() == "" {
		return nil, errors.New("cannot determine node path, 'containerField' is empty")
	}

	parentPath, err := PathOf(container)
	if err != nil {
		if errors.Unwrap(err) == nil {
			return nil, fmt.Errorf(
				"PathOf: error within node of type %T: %w", container, err)
		} else {
			return nil, fmt.Errorf(
				"PathOf: error within container of type %T:\n %w", container, err)
		}
	}
	return parentPath.push(containerField, index), nil
}

// Resolve returns the (deeply) contained child of node denoted by path.
// Leading slashes of path are ignored, the field names are evaluated child by child starting with node.
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

func parsePath(path string) (*fragmentPathImpl, error) {
	result := &fragmentPathImpl{}
	path = strings.TrimLeft(path, "/")

	for segment := range strings.SplitSeq(path, "/") {
		if segment == "" {
			continue
		}
		fieldAndIndex := strings.SplitN(segment, "@", 2)
		field := unique.Make(fieldAndIndex[0])
		if len(fieldAndIndex) == 1 {
			result.push(field, -1)
		} else {
			index, err := strconv.Atoi(fieldAndIndex[1])
			if err != nil {
				return nil, fmt.Errorf("parsePath: index '%s' is not a valid uint: %w", fieldAndIndex[1], err)
			}
			result.push(field, index)
		}
	}
	return result, nil
}
