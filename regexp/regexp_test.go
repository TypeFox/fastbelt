package regexp

import (
	"regexp"
	"testing"
)

func TestSimpleCustom(t *testing.T) {
	regexp := MustCompilRegexp("a+")
	loc := regexp.FindStringIndex("aaab")
	if loc[0] != 0 || loc[1] != 3 {
		panic("TestSimple failed")
	}
}

func TestSimpleOriginal(t *testing.T) {
	regexp := regexp.MustCompile("a+")
	loc := regexp.FindStringIndex("aaab")
	if loc[0] != 0 || loc[1] != 3 {
		panic("TestSimple failed")
	}
}
