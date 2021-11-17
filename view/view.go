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
