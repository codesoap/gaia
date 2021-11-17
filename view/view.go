package view

import (
	"strings"

	"github.com/mattn/go-runewidth"
	"github.com/gdamore/tcell"
)

type View struct {
	Screen      tcell.Screen
	Page        gmi.Page
	CurrentLine int // output line number (may be higher than input line number because of wrapping)
	// TODO: Highlight string
}

// Scrolling to a output line number sucks when resizing...
func (v View) Draw() {
	screenWidth, screenHeight := v.Screen.Size()
	passedLines := 0
	wrapWidth := 72
	if screenWidth < wrapWidth {
		wrapWidth = screenWidth
	}
	for _, line := range v.Page {
		for _, outLine := range wrap(line, wrapWidth) {
			startScreenRow := passedLines - v.CurrentLine // FIXME?!
			drawText(v.Screen, style, outLine, startScreenRow)
			passedLines++
		}
		if passedLines >= v.CurrentLine+screenHeight {
			break // Don't try to draw below-screen content.
		}
	}
}

func wrap(in gmi.Line, width int) []string {
	switch in.(type) {
	case TextLine:
		return wrapWithPrefix(in.String(), "")
	case gmi.PreformattedLine:
		return []string{in.String()}
	case ListLine, H1Line:
		return wrapWithPrefix(in.String(), "  ")
	case LinkLine, H2Line:
		return wrapWithPrefix(in.String(), "   ")
	case H3Line:
		return wrapWithPrefix(in.String(), "    ")
	case QuoteLine:
		txt := in.String()
		prefix := ">"
		if len(txt) >= 2 && txt[1] == ' ' {
			prefix = "> "
		}
		return wrapWithPrefix(in.String(), prefix, width)
	}
}

func wrapWidth(in, prefix string, width int) []string {
	if len(prefix) >= width {
		return []string{"..."}
	}
	ret := []string{}
	prefixWidth := runewidth.StringWidth(prefix)
	split := strings.Split(in, " \t") // TODO: Improve with other whitespaces and dashes.
	                                  //       See strings.FieldsFunc
	split = furtherSplitLongWords(split, prefixWidth, width)
	outLineWidth := 0
	for _, word := range split {
		wordWidth := runewidth.StringWidth(word)
		if outLineWidth + wordWidth > width {
			red = append(ret, outLine.String())
			outLine.Reset()
			outLine.WriteString(prefix)
			outLineWidth = prefixWidth
		}
		outLine.WriteString(word)
	}
}

func furtherSplitLongWords(words []string, prefixWidth, lineWidth int) []string {
	if len(words) = 0 {
		return words
	}
	ret := []string{}

	// First word is special, because it does not get a prefix.
	if runewidth.StringWidth(words[0]) > lineWidth {
		remainder := words[0]
		for runewidth.StringWidth(remainder) > lineWidth {
			i := 2
			for runewidth.StringWidth(remainder[:i]) <= lineWidth {
				i++
			}
			ret = append(ret, remainder[:i]
			remainder = remainder[i:]
		}
		if runewidth.StringWidth(remainder) > 0 {
			ret = append(ret, remainder)
		}
	} else {
		ret = append(ret, words[0])
	}

	for _, word := range words[1:] {
		if runewidth.StringWidth(word) > lineWidth - prefixWidth {
			for runewidth.StringWidth(word) > lineWidth - prefixWidth {
				i := 2
				for runewidth.StringWidth(word[:i]) <= lineWidth - prefixWidth {
					i++
				}
				ret = append(ret, word[:i]
				word = word[i:]
			}
			if runewidth.StringWidth(word) > 0 {
				ret = append(ret, word)
			}
		} else {
			ret = append(ret, word)
		}
	}
	return ret
}
