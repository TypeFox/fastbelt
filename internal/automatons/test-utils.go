package automatons

func Expect(value any) *expectation {
	return &expectation{value: value}
}

type expectation struct {
	value any
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
