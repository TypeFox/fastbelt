package lexer

import core "typefox.dev/fastbelt"

type TokenTypeUsage struct {
	TokenType *core.TokenType
	PushMode  int
	PopMode   bool
	Modifier  int
}

type TokenMode struct {
	Name       string
	TokenTypes []*TokenTypeUsage
	TokenMap   [][]*TokenTypeUsage
}

func UseTokenType(tokenType *core.TokenType) *TokenTypeUsage {
	return &TokenTypeUsage{
		TokenType: tokenType,
		PushMode:  -1,
		PopMode:   false,
		Modifier:  0,
	}
}

func (ttu *TokenTypeUsage) WithPushMode(mode int) *TokenTypeUsage {
	ttu.PushMode = mode
	return ttu
}

func (ttu *TokenTypeUsage) WithPopMode() *TokenTypeUsage {
	ttu.PopMode = true
	return ttu
}

func (ttu *TokenTypeUsage) WithSetMode(mode int) *TokenTypeUsage {
	ttu.PopMode = true
	ttu.PushMode = mode
	return ttu
}

func (ttu *TokenTypeUsage) WithModifier(modifier int) *TokenTypeUsage {
	ttu.Modifier = modifier
	return ttu
}

// IsSkipped reports whether t is routed to the skipped-token group.
func (t *TokenTypeUsage) IsSkipped() bool {
	return t.Modifier == core.SkippedModifier
}

// IsComment reports whether t is routed to the comment-token group.
func (t *TokenTypeUsage) IsComment() bool {
	return t.Modifier == core.CommentModifier
}

// NewTokenMode returns a token mode that recognizes the given token types.
// At each position the longest match wins; among equal-length matches, the
// first argument wins.
func NewTokenMode(name string, tokenTypeUsages ...*TokenTypeUsage) *TokenMode {
	tokenMap := make([][]*TokenTypeUsage, maxChar)
	for i := range maxChar {
		tokenMap[i] = make([]*TokenTypeUsage, 0)
	}
	for _, tokenTypeUsage := range tokenTypeUsages {
		for _, r := range tokenTypeUsage.TokenType.StartChars {
			index := int(r) % maxChar
			tokenMap[index] = append(tokenMap[index], tokenTypeUsage)
		}
	}
	return &TokenMode{
		Name:       name,
		TokenTypes: tokenTypeUsages,
		TokenMap:   tokenMap,
	}
}

// TokenModeStack tracks which [TokenMode] a lexer run is currently in. The
// bottom entry is the mode the run started in and is never removed, so the stack
// is never empty. TokenModeStack is not safe for concurrent use; each
// [DefaultLexer.Exec] uses its own.
type TokenModeStack struct {
	modes []*TokenMode
}

// NewTokenModeStack returns a stack holding defaultMode as its only entry.
func NewTokenModeStack(defaultMode *TokenMode) *TokenModeStack {
	return &TokenModeStack{
		modes: []*TokenMode{defaultMode},
	}
}

// Push makes mode the active mode, keeping the previous one for a later [Pop].
func (s *TokenModeStack) Push(mode *TokenMode) {
	s.modes = append(s.modes, mode)
}

// Pop removes the active mode and returns it, making the mode below it active
// again. Popping the bottom entry is a no-op that returns the start mode: input
// with more pops than pushes stays in the start mode rather than failing.
func (s *TokenModeStack) Pop() *TokenMode {
	if len(s.modes) <= 1 {
		return s.modes[0]
	}
	mode := s.modes[len(s.modes)-1]
	s.modes = s.modes[:len(s.modes)-1]
	return mode
}

// Peek returns the active mode.
func (s *TokenModeStack) Peek() *TokenMode {
	return s.modes[len(s.modes)-1]
}

// SetMode replaces the active mode with mode, leaving the rest of the stack
// untouched. This is what a `mode(X)` command does: unlike [Push] it does not
// deepen the stack, so it cannot be undone by a later [Pop].
func (s *TokenModeStack) SetMode(mode *TokenMode) {
	s.modes[len(s.modes)-1] = mode
}
