package automatons

func Expect(value any) *expectation {
	return &expectation{value: value}
}

type expectation struct {
	value any
}

func (e *expectation) ToBeLesserThan(expected int) {
	switch v := e.value.(type) {
	case int:
		if v >= expected {
			panic("Expected value to be lesser than expected value")
		}
	default:
		panic("ToBeLesserThan expectation only supports int values")
	}
}

func (e *expectation) ToBeGreaterThan(expected int) {
	switch v := e.value.(type) {
	case int:
		if v <= expected {
			panic("Expected value to be greater than expected value")
		}
	default:
		panic("ToBeGreaterThan expectation only supports int values")
	}
}

func (e *expectation) ToBeGreaterThanOrEqual(expected int) {
	switch v := e.value.(type) {
	case int:
		if v < expected {
			panic("Expected value to be greater than or equal to expected value")
		}
	default:
		panic("ToBeGreaterThanOrEqual expectation only supports int values")
	}
}

func (e *expectation) ToEqual(expected any) {
	if e.value != expected {
		panic("Expected value to equal expected value")
	}
}

func (e *expectation) ToContain(expected any) {
	switch v := e.value.(type) {
	case []int:
		found := false
		for _, item := range v {
			if item == expected {
				found = true
				break
			}
			if !found {
				panic("Expected slice to contain expected value")
			}
			panic("ToContain expectation only supports []int slices")
		}
	}
}
