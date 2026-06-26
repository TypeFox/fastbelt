// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package generator

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"typefox.dev/fastbelt/internal/grammar"
	"typefox.dev/fastbelt/workspace"
)

type textMateGrammar struct {
	repository        repository
	scopeName         string
	patterns          []pattern
	injections        map[string]pattern //optional
	injectionSelector string             //optional
	fileTypes         []string           //optional
	name              string             //optional
	firstLineMatch    string             //optional
}

type repository = map[string]pattern

type pattern struct {
	id                  int        //optional
	include             string     //optional
	name                string     //optional
	contentName         string     //optional
	match               string     //optional
	captures            captures   //optional
	begin               string     //optional
	beginCaptures       captures   //optional
	end                 string     //optional
	endCaptures         captures   //optional
	while               string     //optional
	whileCaptures       captures   //optional
	patterns            []pattern  //optional
	repository          repository //optional
	applyEndPatternLast bool       //optional
}

type captures = map[string]pattern

type TextMateGeneratorConfig struct {
	Id              workspace.LanguageID
	FileExtensions  workspace.FileExtensions
	CaseInsensitive bool
}

func GenerateTextMate(grammar grammar.Grammar, config TextMateGeneratorConfig) string {
	obj := textMateGrammar{
		name:       string(config.Id),
		scopeName:  fmt.Sprintf("source.%s", config.Id),
		fileTypes:  []string(config.FileExtensions),
		patterns:   GetPatterns(grammar, config),
		repository: GetRepository(grammar, config),
	}
	result, _ := json.MarshalIndent(obj, "", "  ")
	return string(result)
}

func GetPatterns(grammar grammar.Grammar, config TextMateGeneratorConfig) []pattern {
	patterns := []pattern{}
	patterns = append(patterns, pattern{
		include: "#comments",
	})
	patterns = append(patterns, GetControlKeywords(grammar, config))
	patterns = append(patterns, GetStringPatterns(grammar, config)...)

	return patterns
}

func GetRepository(grammar grammar.Grammar, config TextMateGeneratorConfig) repository {
	repository := repository{}
	commentPatterns := []pattern{}
	var stringEscapePattern *pattern
	for _, terminal := range grammar.Terminals() {
		if terminal.Type() == "comment" {
			parts := GetTerminalParts(terminal.Regexp())
			for _, part := range parts {
				if part.end != "" {
					commentPatterns = append(commentPatterns, pattern{
						name:  fmt.Sprintf(`comment.block.%s`, config.Id),
						begin: part.start,
						beginCaptures: captures{
							"0": pattern{
								name: fmt.Sprintf(`punctuation.definition.comment.%s`, config.Id),
							},
						},
						end: part.end,
						endCaptures: captures{
							"0": pattern{
								name: fmt.Sprintf(`punctuation.definition.comment.%s`, config.Id),
							},
						},
					})
				} else {
					commentPatterns = append(commentPatterns, pattern{
						begin: part.start,
						beginCaptures: captures{
							"1": pattern{
								name: fmt.Sprintf(`punctuation.whitespace.comment.leading.%s`, config.Id),
							},
						},
						end:  `(?=$)`,
						name: fmt.Sprintf(`comment.line.%s`, config.Id),
					})
				}
			}
		} else if strings.ToLower(terminal.Name()) == "string" {
			stringEscapePattern = &pattern{
				name:  fmt.Sprintf(`constant.character.escape.%s`, config.Id),
				match: `\\(x[0-9A-Fa-f]{2}|u[0-9A-Fa-f]{4}|u\{[0-9A-Fa-f]+\}|[0-2][0-7]{0,2}|3[0-6][0-7]?|37[0-7]?|[4-7][0-7]?|.|$)`,
			}
		}

		if len(commentPatterns) > 0 {
			repository["comments"] = pattern{
				patterns: commentPatterns,
			}
		}
		if stringEscapePattern != nil {
			repository["string-character-escape"] = *stringEscapePattern
		}
	}
	return repository
}

func GetControlKeywords(grammr grammar.Grammar, config TextMateGeneratorConfig) pattern {
	regex := regexp.MustCompile(`[A-Za-z]`)
	controlKeywords := []grammar.Keyword{}
	for _, keyword := range GetAllKeywords(grammr) {
		if regex.MatchString(keyword.Text()) {
			controlKeywords = append(controlKeywords, keyword)
		}
	}
	groups := GroupKeywords(controlKeywords)
	return pattern{
		name:  fmt.Sprintf(`keyword.control.%s`, config.Id),
		match: fmt.Sprintf(`%s%s`, ifThenElse(config.CaseInsensitive, `(?i)`, ``), strings.Join(groups, "|")),
	}
}

func ifThenElse(condition bool, trueValue string, falseValue string) string {
	if condition {
		return trueValue
	} else {
		return falseValue
	}
}

type group struct {
	letter       []string
	leftSpecial  []string
	rightSpecial []string
	special      []string
}

func GroupKeywords(keywords []grammar.Keyword) []string {
	groups := group{
		letter:       []string{},
		leftSpecial:  []string{},
		rightSpecial: []string{},
		special:      []string{},
	}

	for _, keyword := range keywords {
		keywordPattern := regexp.QuoteMeta(keyword.Text())
		if regexp.MustCompile(`\w`).MatchString(string(keyword.Text()[0])) {
			if regexp.MustCompile(`\w`).MatchString(string(keyword.Text()[len(keyword.Text())-1])) {
				groups.letter = append(groups.letter, keywordPattern)
			} else {
				groups.rightSpecial = append(groups.rightSpecial, keywordPattern)
			}
		} else {
			if regexp.MustCompile(`\w`).MatchString(string(keyword.Text()[len(keyword.Text())-1])) {
				groups.leftSpecial = append(groups.leftSpecial, keywordPattern)
			} else {
				groups.special = append(groups.special, keywordPattern)
			}
		}
	}

	res := []string{}
	if len(groups.letter) > 0 {
		res = append(res, fmt.Sprintf(`\b(%s)\b`, strings.Join(groups.letter, "|")))
	}
	if len(groups.leftSpecial) > 0 {
		res = append(res, fmt.Sprintf(`\B(%s)\b`, strings.Join(groups.leftSpecial, "|")))
	}
	if len(groups.rightSpecial) > 0 {
		res = append(res, fmt.Sprintf(`\b(%s)\B`, strings.Join(groups.rightSpecial, "|")))
	}
	if len(groups.special) > 0 {
		res = append(res, fmt.Sprintf(`\B(%s)\B`, strings.Join(groups.special, "|")))
	}
	return res
}

func GetStringPatterns(grammr grammar.Grammar, config TextMateGeneratorConfig) []pattern {
	terminals := grammr.Terminals()
	var stringTerminal grammar.Token
	for _, terminal := range terminals {
		if strings.ToLower(terminal.Name()) == "string" {
			stringTerminal = terminal
			break
		}
	}
	stringPatterns := []pattern{}
	if stringTerminal != nil {
		parts := GetTerminalParts(stringTerminal.Regexp())
		for _, part := range parts {
			if part.end != "" {
				stringPatterns = append(stringPatterns, pattern{
					name:  fmt.Sprintf(`string.quoted.%s.%s`, delimiterName(part.start), config.Id),
					begin: part.start,
					end:   part.end,
					patterns: []pattern{
						{
							include: "#string-character-escape",
						},
					},
				})
			}
		}
	}
	return stringPatterns
}

func delimiterName(delimiter string) string {
	if delimiter == "'" {
		return "single"
	} else if delimiter == `"` {
		return "double"
	} else if delimiter == "`" {
		return "backtick"
	} else {
		return "delimiter"
	}
}
