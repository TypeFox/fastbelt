package fastbelt

func DefaultLink(scope Scope, text string) (*AstNodeDescription, *ReferenceError) {
	descriptions := scope.ElementsByName(text)
	if len(descriptions) == 0 {
		return nil, NewReferenceError("Could not resolve reference to '" + text + "'.")
	} else {
		return descriptions[0], nil
	}
}

type Scope interface {
	ElementsByName(name string) []*AstNodeDescription
	AllElements() []*AstNodeDescription
}

type emptyScope struct{}

func (s *emptyScope) ElementsByName(name string) []*AstNodeDescription {
	return []*AstNodeDescription{}
}

func (s *emptyScope) AllElements() []*AstNodeDescription {
	return []*AstNodeDescription{}
}

var EmptyScope Scope = &emptyScope{}

type MapScope struct {
	elements map[string][]*AstNodeDescription
	outer    Scope
}

func NewMapScope(elements map[string][]*AstNodeDescription, outer Scope) *MapScope {
	return &MapScope{
		elements: elements,
		outer:    outer,
	}
}

func (s *MapScope) ElementsByName(name string) []*AstNodeDescription {
	if elems, ok := s.elements[name]; ok {
		return elems
	} else if s.outer != nil {
		return s.outer.ElementsByName(name)
	}
	return []*AstNodeDescription{}
}

func (s *MapScope) AllElements() []*AstNodeDescription {
	result := []*AstNodeDescription{}
	for _, elems := range s.elements {
		result = append(result, elems...)
	}
	if s.outer != nil {
		result = append(result, s.outer.AllElements()...)
	}
	return result
}
