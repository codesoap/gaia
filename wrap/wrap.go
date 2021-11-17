package wrap

import (
	"fmt"
	"strings"

	"github.com/mattn/go-runewidth"
)

// TODO: WrapWithPrefix
func WrapWithPrefix(in, prefix string, width int) []string {
	prefixWidth := runewidth.StringWidth(prefix)
	if prefixWidth >= width {
		return []string{"..."}
	}
	return Wrap(in, width - prefixWidth) // FIXME
}

// TODO: doc
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
	// TODO: improve
	return r == ' ' || r == '\t'
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
