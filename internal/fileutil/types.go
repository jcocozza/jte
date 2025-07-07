package fileutil

import (
	"path/filepath"
	"regexp"
	"sort"
)

type SyntaxType int

const (
	NORMAL SyntaxType = iota
	NUMBER
	STRING
	COMMENT
	FUNCTION
	IDENTIFIER
	KEYWORD
)
// colors
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

func SyntaxTypeToColor(st SyntaxType) string {
	switch st {
	case KEYWORD:
		return Red
	case STRING:
		return Green
	case COMMENT:
		return Cyan
	case FUNCTION:
		return Magenta
	case IDENTIFIER:
		return Blue
		//case MATCH:
		//	return Yellow
	case NORMAL:
		fallthrough
	default:
		return White
	}
}

// different file types
//
// to add a file type
// 1. add a new const
// 2. add to the FileTypes list
// 3. add to the Registry
type FileType int

func (f *FileType) String() string {
	return FileTypes[*f]
}

func DetermineFileType(path string) FileType {
	ext := filepath.Ext(path)
	ext = ext[1:] // remove the "."
	switch ext {
	case "go":
		return Go
	case "py":
		return Python
	default:
		return Unknown
	}
}

const (
	Unknown FileType = iota
	Go
	Python
)

var FileTypes = [...]string{
	Unknown: "unknown",
	Go:      "go",
	Python:  "python",
}

var FileTypeRegistry = map[FileType]map[SyntaxType]*regexp.Regexp{
	Go: {
		KEYWORD: regexp.MustCompile(`\b(if|else|for|return|func|package|panic|map|var)\b`),
		STRING:  regexp.MustCompilePOSIX(`"(.*?)"`),
		COMMENT: regexp.MustCompile(`//.*`),
	},
	Python: {},
}

type Match struct {
	start int
	end   int
	style SyntaxType
}

type Token struct {
	Text  string
	Color string
}

func GetMatches(fileType FileType, line string) []Token {
	matches := []Match{}
	for s, r := range FileTypeRegistry[fileType] {
		allMatches := r.FindAllStringIndex(line, -1)
		for _, match := range allMatches {
			matches = append(matches, Match{
				start: match[0],
				end:   match[1],
				style: s,
			})
		}
	}
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].start < matches[j].start
	})
	// Now build the tokens based on the sorted matches
	tokens := []Token{}
	previousEnd := 0
	for _, match := range matches {
		// Handle plain text before the match
		if match.start > previousEnd {
			tokens = append(tokens, Token{
				Text:  line[previousEnd:match.start],
				Color: SyntaxTypeToColor(NORMAL),
			})
		}

		// Add the highlighted token
		tokens = append(tokens, Token{
			Text:  line[match.start:match.end],
			Color: SyntaxTypeToColor(match.style),
		})

		previousEnd = match.end
	}
	// Handle any remaining text after the last match
	if previousEnd < len(line) {
		tokens = append(tokens, Token{
			Text:  line[previousEnd:],
			Color: SyntaxTypeToColor(NORMAL),
		})
	}
	return tokens
}
