package regexp

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func checkRegexp(regexp Regexp, input string, expected []int) {
	loc := regexp.FindStringIndex(input)
	if (loc == nil && expected != nil) || (loc != nil && expected == nil) {
		panic("Location mismatch")
	}
	if loc != nil && expected != nil {
		if len(loc) != len(expected) {
			panic("Location length mismatch")
		}
		for i := range loc {
			if loc[i] != expected[i] {
				panic("Location value mismatch")
			}
		}
	}
}

func TestSimple(t *testing.T) {
	regexp := MustCompilRegexp("a+")
	checkRegexp(regexp, "aaab", []int{0, 3})
	checkRegexp(regexp, "ab", []int{0, 1})
	checkRegexp(regexp, "aace", []int{0, 2})
}

func TestEmpty(t *testing.T) {
	regexp := MustCompilRegexp("")
	checkRegexp(regexp, "", []int{0, 0})
	checkRegexp(regexp, "123", []int{0, 0})
}

func TestStar(t *testing.T) {
	regexp := MustCompilRegexp("a*")
	checkRegexp(regexp, "", []int{0, 0})
	checkRegexp(regexp, "a", []int{0, 1})
	checkRegexp(regexp, "aaaabbb", []int{0, 4})
}

func TestEmail(t *testing.T) {
	regexp := MustCompilRegexp(`[\w\.+-]+@[\w\.-]+\.[\w\.-]+`)
	loc := regexp.FindStringIndex("a.b@c.de")
	if loc[0] != 0 || loc[1] != 8 {
		panic("TestEmail failed")
	}
}

func TestIP(t *testing.T) {
	regexp := MustCompilRegexp(`(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9])\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9])`)
	loc := regexp.FindStringIndex("255.241.123.10")
	if loc == nil || loc[0] != 0 || loc[1] != 14 {
		panic("TestIP failed")
	}
}

func TestGenerator(t *testing.T) {
	regexp := MustCompilRegexp(`(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9])\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9])`)
	file, err := os.Create("../example.go")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Write a string to the file
	root := NewNode()
	root.AppendLine("package main")
	root.AppendLine("// This file is generated. Do not edit.")
	root.AppendNode(regexp.(*regexpImpl).GenerateFindStringIndex("IP"))
	fmt.Fprintf(file, "%s", root)
}
