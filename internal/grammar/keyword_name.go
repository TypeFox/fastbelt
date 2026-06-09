// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package grammar

import "strings"

// KeywordValue returns the string value of k stripped of its surrounding quotes.
func KeywordValue(k Keyword) string {
	return k.Value()[1 : len(k.Value())-1]
}

// KeywordName converts the keyword's text to a PascalCase identifier suitable
// for use in generated Go constant names.
func KeywordName(k Keyword) string {
	sb := &strings.Builder{}
	for _, r := range KeywordValue(k) {
		switch r {
		case '(':
			sb.WriteString("LeftParen")
		case ')':
			sb.WriteString("RightParen")
		case '{':
			sb.WriteString("LeftBrace")
		case '}':
			sb.WriteString("RightBrace")
		case '[':
			sb.WriteString("LeftBracket")
		case ']':
			sb.WriteString("RightBracket")
		case '_':
			sb.WriteString("Underscore")
		case '$':
			sb.WriteString("Dollar")
		case '%':
			sb.WriteString("Percent")
		case '#':
			sb.WriteString("Hash")
		case '@':
			sb.WriteString("At")
		case '!':
			sb.WriteString("Exclamation")
		case '^':
			sb.WriteString("Caret")
		case '&':
			sb.WriteString("Ampersand")
		case '*':
			sb.WriteString("Asterisk")
		case '-':
			sb.WriteString("Dash")
		case '+':
			sb.WriteString("Plus")
		case '=':
			sb.WriteString("Equals")
		case '<':
			sb.WriteString("LessThan")
		case '>':
			sb.WriteString("GreaterThan")
		case '?':
			sb.WriteString("Question")
		case '/':
			sb.WriteString("Slash")
		case '\\':
			sb.WriteString("Backslash")
		case '|':
			sb.WriteString("Pipe")
		case '~':
			sb.WriteString("Tilde")
		case '`':
			sb.WriteString("Backtick")
		case '.':
			sb.WriteString("Dot")
		case ',':
			sb.WriteString("Comma")
		case ':':
			sb.WriteString("Colon")
		case ';':
			sb.WriteString("Semicolon")
		case ' ':
			sb.WriteString("Space")
		case '\t':
			sb.WriteString("Tab")
		case '\n':
			sb.WriteString("Newline")
		case '\r':
			sb.WriteString("CarriageReturn")
		default:
			sb.WriteRune(r)
		}
	}
	return sb.String()
}
