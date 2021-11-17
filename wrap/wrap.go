package wrap

import (
	"fmt"
	"strings"

	"github.com/mattn/go-runewidth"
)

// TODO: doc
func WrapToWidth(in, prefix string, width int) []string {
	prefixWidth := runewidth.StringWidth(prefix)
	if prefixWidth >= width {
		return []string{"..."}
	}
	ret := []string{}
	split := strings.Split(in, " \t") // TODO: Improve with other whitespaces and dashes.
	split = furtherSplitLongWords(split, prefixWidth, width)
	outLine := []string{}
	outLineWidth := 0
	for _, word := range split {
		wordWidth := runewidth.StringWidth(word)
		if prefixWidth + outLineWidth + wordWidth <= width {
			outLine = append(outLine, word)
			outLineWidth += 1 + wordWidth // 1 for the space.
		} else {
			ret = append(ret, prefix+strings.Join(outLine, " "))
			outLine = []string{word}
			outLineWidth = wordWidth
		}
	}
	if len(outLine) > 0 {
		ret = append(ret, prefix+strings.Join(outLine, " "))
	}
	return ret
}

func furtherSplitLongWords(words []string, prefixWidth, lineWidth int) []string {
	if lineWidth == 0 {
		panic(fmt.Errorf("cannot wrap to width 0"))
	} else if prefixWidth >= lineWidth {
		panic(fmt.Errorf("prefix is longer than line"))
	} else if len(words) == 0 {
		return words
	}
	ret := []string{}
	for _, word := range words {
		ret = append(ret, furtherSplitLongWord(word, prefixWidth, lineWidth)...)
	}
	return ret
}

func furtherSplitLongWord(word string, prefixWidth, lineWidth int) []string {
	ret := []string{}
	if runewidth.StringWidth(word) > lineWidth-prefixWidth {
		for runewidth.StringWidth(word) > lineWidth-prefixWidth {
			i := 2
			for runewidth.StringWidth(word[:i]) <= lineWidth-prefixWidth {
				i++
			}
			ret = append(ret, word[:i-1])
			word = word[i-1:]
		}
		if runewidth.StringWidth(word) > 0 {
			ret = append(ret, word)
		}
	} else {
		ret = append(ret, word)
	}
	return ret
}
