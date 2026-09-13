// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package railroad

import (
	"strconv"
	"strings"
)

// svgAttr is a single, order-preserved SVG attribute.
type svgAttr struct {
	name  string
	value string
}

// svgElement is a minimal SVG element builder: just enough to emit the
// <svg>/<g>/<path>/<rect>/<text>/<style> elements this package needs, with
// deterministic attribute order (useful for tests).
type svgElement struct {
	tag      string
	attrs    []svgAttr
	children []*svgElement
	text     string
	hasText  bool
}

func el(tag string) *svgElement {
	return &svgElement{tag: tag}
}

func (e *svgElement) attr(name, value string) *svgElement {
	e.attrs = append(e.attrs, svgAttr{name, value})
	return e
}

func (e *svgElement) attrf(name string, value float64) *svgElement {
	return e.attr(name, formatNum(value))
}

func (e *svgElement) add(children ...*svgElement) *svgElement {
	e.children = append(e.children, children...)
	return e
}

func (e *svgElement) setText(s string) *svgElement {
	e.text = s
	e.hasText = true
	return e
}

func (e *svgElement) writeTo(b *strings.Builder) {
	b.WriteByte('<')
	b.WriteString(e.tag)
	for _, a := range e.attrs {
		b.WriteByte(' ')
		b.WriteString(a.name)
		b.WriteString(`="`)
		b.WriteString(escapeAttr(a.value))
		b.WriteByte('"')
	}
	if !e.hasText && len(e.children) == 0 {
		b.WriteString("/>")
		return
	}
	b.WriteByte('>')
	if e.hasText {
		b.WriteString(escapeText(e.text))
	}
	for _, c := range e.children {
		c.writeTo(b)
	}
	b.WriteString("</")
	b.WriteString(e.tag)
	b.WriteByte('>')
}

func (e *svgElement) String() string {
	var b strings.Builder
	e.writeTo(&b)
	return b.String()
}

// formatNum formats a coordinate without trailing zeros or exponents. Every
// value in this package is built from sums, differences, and halvings of
// small integers, all of which are exact in float64, so no rounding is
// needed to avoid binary-fraction noise.
func formatNum(v float64) string {
	return strconv.FormatFloat(v, 'f', -1, 64)
}

func escapeAttr(s string) string {
	var b strings.Builder
	for _, r := range s {
		switch r {
		case '&':
			b.WriteString("&amp;")
		case '"':
			b.WriteString("&quot;")
		case '<':
			b.WriteString("&lt;")
		case '>':
			b.WriteString("&gt;")
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}

func escapeText(s string) string {
	var b strings.Builder
	for _, r := range s {
		switch r {
		case '&':
			b.WriteString("&amp;")
		case '<':
			b.WriteString("&lt;")
		case '>':
			b.WriteString("&gt;")
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}
