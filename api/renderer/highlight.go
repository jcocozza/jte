package renderer

import (
	"regexp"
	"sort"
)

type SyntaxGroup int

const (
	SG_NORMAL SyntaxGroup = iota
	SG_NUMBER
	SG_MATCH
	SG_STRING
	SG_COMMENT
	SG_FUNCTION
	SG_IDENTIFIER
	SG_KEYWORD
)

const (
	Reset      = "\x1b[0m"
	Bold       = "\x1b[1m"
	Italic     = "\x1b[3m"
	Underline  = "\x1b[4m"
	Black      = "\x1b[30m"
	Red        = "\x1b[31m"
	Green      = "\x1b[32m"
	Yellow     = "\x1b[33m"
	Blue       = "\x1b[34m"
	Magenta    = "\x1b[35m"
	Cyan       = "\x1b[36m"
	White      = "\x1b[37m"
	Background = "\x1b[48;5;%dm" // for background color customization
)

func syntaxToColor(hl SyntaxGroup) string {
	switch hl {
	case SG_KEYWORD:
		return Red
	case SG_STRING:
		return Green
	case SG_COMMENT:
		return Cyan
	case SG_FUNCTION:
		return Magenta
	case SG_IDENTIFIER:
		return Blue
	case SG_MATCH:
		return Yellow
	case SG_NORMAL:
		fallthrough
	default:
		return White
	}
}

type Token struct {
	Text string
	SG   SyntaxGroup
}

// this is mostly AI generated
// probably trash, but seems to be working thus far
func highlightLine(line string, searchPattern string) []Token {
	var tokens []Token
	type Match struct {
		start, end int
		style      SyntaxGroup
	}

	var matches []Match

	// Define regex patterns
	keywordRegex := regexp.MustCompile(`\b(if|else|for|return|func|package)\b`)
	stringRegex := regexp.MustCompile(`"(.*?)"`)
	commentRegex := regexp.MustCompile(`//.*`)
	searchRegex := regexp.MustCompile(searchPattern) //TODO: this can panic, so we have to handle the error properly

	// Find matches for keywords
	keywordMatches := keywordRegex.FindAllStringIndex(line, -1)
	for _, match := range keywordMatches {
		matches = append(matches, Match{
			start: match[0],
			end:   match[1],
			style: SG_KEYWORD,
		})
	}

	// Find matches for strings
	stringMatches := stringRegex.FindAllStringIndex(line, -1)
	for _, match := range stringMatches {
		matches = append(matches, Match{
			start: match[0],
			end:   match[1],
			style: SG_STRING,
		})
	}

	// Find matches for comments
	commentMatches := commentRegex.FindAllStringIndex(line, -1)
	for _, match := range commentMatches {
		matches = append(matches, Match{
			start: match[0],
			end:   match[1],
			style: SG_COMMENT,
		})
	}

	if len(searchPattern) != 0 {
		searchMatches := searchRegex.FindAllStringIndex(line, -1)
		for _, match := range searchMatches {
			matches = append(matches, Match{
				start: match[0],
				end:   match[1],
				style: SG_MATCH,
			})
		}
	}

	// Sort matches by the starting position (important for correct token order)
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].start < matches[j].start
	})

	// Now build the tokens based on the sorted matches
	previousEnd := 0
	for _, match := range matches {
		// Handle plain text before the match
		if match.start > previousEnd {
			tokens = append(tokens, Token{
				Text: line[previousEnd:match.start],
				SG:   SG_NORMAL,
			})
		}

		// Add the highlighted token
		tokens = append(tokens, Token{
			Text: line[match.start:match.end],
			SG:   match.style,
		})

		previousEnd = match.end
	}

	// Handle any remaining text after the last match
	if previousEnd < len(line) {
		tokens = append(tokens, Token{
			Text: line[previousEnd:],
			SG:   SG_NORMAL,
		})
	}
	return tokens
}
