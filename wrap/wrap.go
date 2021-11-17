package wrap

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/mattn/go-runewidth"
)

// WrapWithPrefix wraps in so that none of the returned lines are longer
// than width. prefix is added to the front of every returned line. It
// will be attempted to break lines on whitespaces.
//
// All space runes will be replaced by 0x20 in the process.
func WrapWithPrefix(in, prefix string, width int) []string {
	prefixWidth := runewidth.StringWidth(prefix)
	if prefixWidth >= width {
		return []string{"..."}
	}
	lines := Wrap(in, width - prefixWidth) // FIXME
	if prefixWidth == 0 {
		return lines
	}
	prefixedLines := make([]string, 0, len(lines))
	for _, line := range lines {
		prefixedLines = append(prefixedLines, prefix + line)
	}
	return prefixedLines
}

// Wrap wraps in so that none of the returned lines are longer than
// width. It will be attempted to break lines on whitespaces.
//
// All space runes will be replaced by 0x20 in the process.
func Wrap(in string, width int) []string {
	ret := []string{}
	split := strings.FieldsFunc(in, isSplitRune)
	split = splitLongWords(split, width)
	outLine := []string{}
	outLineWidth := 0
	for _, word := range split {
		wordWidth := runewidth.StringWidth(word)
		if outLineWidth + wordWidth <= width {
			outLine = append(outLine, word)
			outLineWidth += 1 + wordWidth // 1 for the space.
		} else {
			ret = append(ret, strings.Join(outLine, " "))
			outLine = []string{word}
			outLineWidth = wordWidth
		}
	}
	if len(outLine) > 0 {
		ret = append(ret, strings.Join(outLine, " "))
	}
	return ret
}

func isSplitRune(r rune) bool {
	// \u00A0 is NBSP
	return r != '\u00A0' && unicode.IsSpace(r)
}

func splitLongWords(words []string, lineWidth int) []string {
	if lineWidth == 0 {
		panic(fmt.Errorf("cannot wrap to width 0"))
	} else if len(words) == 0 {
		return words
	}
	ret := []string{}
	for _, word := range words {
		ret = append(ret, splitLongWord(word, lineWidth)...)
	}
	return ret
}

func splitLongWord(word string, lineWidth int) []string {
	if runewidth.StringWidth(word) <= lineWidth {
		return []string{word}
	}
	ret := []string{}
	for runewidth.StringWidth(word) > lineWidth {
		i := 2
		for runewidth.StringWidth(word[:i]) <= lineWidth {
			i++
		}
		ret = append(ret, word[:i-1])
		word = word[i-1:]
	}
	if runewidth.StringWidth(word) > 0 {
		ret = append(ret, word)
	}
	return ret
}
