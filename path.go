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

// FragmentPath describes node path descriptors pointing to a child AST node within a parent.
type FragmentPath interface {
	// Empty returns true if this descriptor is empty and equal to the empty path string, returns false otherwise.
	Empty() bool
	// Head returns field name (unique.Handle) and index of the head segment of the segment list.
	// index is expected to be negative (-1) for non-list fields.
	Head() (unique.Handle[string], int)
	// Tail returns a copy of this path descriptor with the first element dropped from the segments list.
	Tail() FragmentPath
	// String returns the corresponding fragment path string
	String() string
}

type fragmentPath []segment
type segment struct {
	name  unique.Handle[string]
	index int
}

var _ FragmentPath = (fragmentPath)(nil)

func (p fragmentPath) Empty() bool {
	return len(p) == 0
}

func (p fragmentPath) Head() (name unique.Handle[string], index int) {
	if len(p) == 0 {
		return fieldZero, -1
	} else {
		head := p[0]
		return head.name, head.index
	}
}

func (p fragmentPath) Tail() FragmentPath {
	if len(p) == 0 {
		return p
	} else {
		return p[1:]
	}
}

func (p fragmentPath) String() string {
	var sb strings.Builder
	for _, s := range p {
		sb.WriteString("/")
		sb.WriteString(s.name.Value())
		if s.index >= 0 {
			sb.WriteString("@")
			sb.WriteString(strconv.Itoa(s.index))
		}
	}
	return sb.String()
}

var fieldZero = unique.Handle[string]{}
var fieldEmpty = unique.Make("")

// PathOf composes a [FragmentPath] denoting node's path within its root container.
// It does so in a recursive manner based on node's [AstNode.ContainmentData].
// Calling [FragmentPath.String]() on the result yields a slash-separated
// path string that uniquely identifies this node within its containment hierarchy.
// e.g. "/rules@2/alternatives@0".
// Returns an empty descriptor for root nodes & nodes with no configured container,
// its [FragmentPath.String]() yields an empty string.
func PathOf(node AstNode) (fragmentPath, error) {
	container := node.Container()
	containerField, index := node.ContainmentData()

	if container == nil {
		// initialize the result fragment path with an initial capacity of 5 items (assumed common upper bound)
		return make(fragmentPath, 0, 5), nil
	} else if containerField == fieldZero || containerField == fieldEmpty {
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
	return append(parentPath, segment{containerField, index}), nil
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

func parsePath(path string) (fragmentPath, error) {
	// initialize the result fragment path with an initial capacity of 5 items (assumed common upper bound)
	result := make(fragmentPath, 0, 5)

	for curSegment := range strings.SplitSeq(path, "/") {
		if curSegment == "" {
			continue
		}
		fieldAndIndex := strings.SplitN(curSegment, "@", 2)
		field := unique.Make(fieldAndIndex[0])
		index := -1
		if len(fieldAndIndex) == 2 {
			err := error(nil)
			index, err = strconv.Atoi(fieldAndIndex[1])
			if err != nil {
				return nil, fmt.Errorf("parsePath: index '%s' is not a valid uint: %w", fieldAndIndex[1], err)
			}
		}
		result = append(result, segment{field, index})
	}
	return result, nil
}
