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

type TextMateGrammar struct {
	Repository        Repository         "json:\"repository\""
	ScopeName         string             "json:\"scopeName\""
	Patterns          []Pattern          "json:\"patterns\""
	Injections        map[string]Pattern "json:\"injections,omitempty\""
	InjectionSelector string             "json:\"injectionSelector,omitempty\""
	FileTypes         []string           "json:\"fileTypes,omitempty\""
	Name              string             "json:\"name,omitempty\""
	FirstLineMatch    string             "json:\"firstLineMatch,omitempty\""
}

type Repository = map[string]Pattern

type Pattern struct {
	Id                  int        "json:\"id,omitempty\""
	Include             string     "json:\"include,omitempty\""
	Name                string     "json:\"name,omitempty\""
	ContentName         string     "json:\"contentName,omitempty\""
	Match               string     "json:\"match,omitempty\""
	Captures            Captures   "json:\"captures,omitempty\""
	Begin               string     "json:\"begin,omitempty\""
	BeginCaptures       Captures   "json:\"beginCaptures,omitempty\""
	End                 string     "json:\"end,omitempty\""
	EndCaptures         Captures   "json:\"endCaptures,omitempty\""
	While               string     "json:\"while,omitempty\""
	WhileCaptures       Captures   "json:\"whileCaptures,omitempty\""
	Patterns            []Pattern  "json:\"patterns,omitempty\""
	Repository          Repository "json:\"repository,omitempty\""
	ApplyEndPatternLast bool       "json:\"applyEndPatternLast,omitempty\""
}

type Captures = map[string]Pattern

type TextMateGeneratorConfig struct {
	Id              workspace.LanguageID
	FileExtensions  workspace.FileExtensions
	CaseInsensitive bool
}

func GenerateTextMate(grammar grammar.Grammar, config TextMateGeneratorConfig) string {
	obj := TextMateGrammar{
		Name:       string(config.Id),
		ScopeName:  fmt.Sprintf("source.%s", config.Id),
		FileTypes:  []string(config.FileExtensions),
		Patterns:   GetPatterns(grammar, config),
		Repository: GetRepository(grammar, config),
	}
	result, _ := json.MarshalIndent(obj, "", "  ")
	return string(result)
}

func GetPatterns(grammar grammar.Grammar, config TextMateGeneratorConfig) []Pattern {
	patterns := []Pattern{}
	patterns = append(patterns, Pattern{
		Include: "#comments",
	})
	patterns = append(patterns, GetControlKeywords(grammar, config))
	patterns = append(patterns, GetStringPatterns(grammar, config)...)

	return patterns
}

func GetRepository(grammar grammar.Grammar, config TextMateGeneratorConfig) Repository {
	repository := Repository{}
	commentPatterns := []Pattern{}
	var stringEscapePattern *Pattern
	for _, terminal := range grammar.Terminals() {
		if terminal.Type() == "comment" {
			regexPattern := terminal.Regexp()
			regexPattern = regexPattern[1 : len(regexPattern)-1]
			parts := GetTerminalParts(regexPattern)
			for _, part := range parts {
				if part.end != "" {
					commentPatterns = append(commentPatterns, Pattern{
						Name:  fmt.Sprintf(`comment.block.%s`, config.Id),
						Begin: part.start,
						BeginCaptures: Captures{
							"0": Pattern{
								Name: fmt.Sprintf(`punctuation.definition.comment.%s`, config.Id),
							},
						},
						End: part.end,
						EndCaptures: Captures{
							"0": Pattern{
								Name: fmt.Sprintf(`punctuation.definition.comment.%s`, config.Id),
							},
						},
					})
				} else {
					commentPatterns = append(commentPatterns, Pattern{
						Begin: part.start,
						BeginCaptures: Captures{
							"1": Pattern{
								Name: fmt.Sprintf(`punctuation.whitespace.comment.leading.%s`, config.Id),
							},
						},
						End:  `(?=$)`,
						Name: fmt.Sprintf(`comment.line.%s`, config.Id),
					})
				}
			}
		} else if strings.ToLower(terminal.Name()) == "string" {
			stringEscapePattern = &Pattern{
				Name:  fmt.Sprintf(`constant.character.escape.%s`, config.Id),
				Match: `\\(x[0-9A-Fa-f]{2}|u[0-9A-Fa-f]{4}|u\{[0-9A-Fa-f]+\}|[0-2][0-7]{0,2}|3[0-6][0-7]?|37[0-7]?|[4-7][0-7]?|.|$)`,
			}
		}

		if len(commentPatterns) > 0 {
			repository["comments"] = Pattern{
				Patterns: commentPatterns,
			}
		}
		if stringEscapePattern != nil {
			repository["string-character-escape"] = *stringEscapePattern
		}
	}
	return repository
}

func GetControlKeywords(grammr grammar.Grammar, config TextMateGeneratorConfig) Pattern {
	regex := regexp.MustCompile(`[A-Za-z]`)
	controlKeywords := []grammar.Keyword{}
	for _, keyword := range GetAllKeywords(grammr) {
		if regex.MatchString(keyword.Text()) {
			controlKeywords = append(controlKeywords, keyword)
		}
	}
	groups := GroupKeywords(controlKeywords)
	return Pattern{
		Name:  fmt.Sprintf(`keyword.control.%s`, config.Id),
		Match: fmt.Sprintf(`%s%s`, ifThenElse(config.CaseInsensitive, `(?i)`, ``), strings.Join(groups, "|")),
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

func GetStringPatterns(grammr grammar.Grammar, config TextMateGeneratorConfig) []Pattern {
	terminals := grammr.Terminals()
	var stringTerminal grammar.Token
	for _, terminal := range terminals {
		if strings.ToLower(terminal.Name()) == "string" {
			stringTerminal = terminal
			break
		}
	}
	stringPatterns := []Pattern{}
	if stringTerminal != nil {
		regexPattern := stringTerminal.Regexp()
		regexPattern = regexPattern[1 : len(regexPattern)-1]
		parts := GetTerminalParts(regexPattern)
		for _, part := range parts {
			if part.end != "" {
				stringPatterns = append(stringPatterns, Pattern{
					Name:  fmt.Sprintf(`string.quoted.%s.%s`, delimiterName(part.start), config.Id),
					Begin: part.start,
					End:   part.end,
					Patterns: []Pattern{
						{
							Include: "#string-character-escape",
						},
					},
				})
			}
		}
	}
	return stringPatterns
}

func delimiterName(delimiter string) string {
	switch delimiter {
	case "'":
		return "single"
	case `"`:
		return "double"
	case "`":
		return "backtick"
	default:
		return "delimiter"
	}
}
