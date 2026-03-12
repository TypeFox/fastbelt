// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package allstar

import parser "typefox.dev/fastbelt/parser"

// parserBridge wraps LLStarLookahead to satisfy parser.LookaheadStrategy.
// *parser.ParserState satisfies allstar.TokenSource via its LA method.
type parserBridge struct{ *LLStarLookahead }

func (b *parserBridge) Predict(src *parser.ParserState, key string) int {
	return b.PredictAlternation(src, key)
}

func (b *parserBridge) PredictOpt(src *parser.ParserState, key string) bool {
	return b.PredictOptional(src, key)
}

// AsParserStrategy returns a parser.LookaheadStrategy backed by this ALL(*) engine.
func (s *LLStarLookahead) AsParserStrategy() parser.LookaheadStrategy {
	return &parserBridge{s}
}
