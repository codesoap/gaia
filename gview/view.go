package gview

import (
	"fmt"

	"github.com/codesoap/gaia/gmi"
	"github.com/codesoap/gaia/gwrap"
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
)

type View struct {
	Page        gmi.Page
	CurrentLine int // output line number (may be higher than input line number because of wrapping)
	// TODO: Highlight string
}

func emitStr(s tcell.Screen, x, y int, style tcell.Style, str string) {
	for _, c := range str {
		var comb []rune
		w := runewidth.RuneWidth(c)
		if w == 0 {
			comb = []rune{c}
			c = ' '
			w = 1
		}
		s.SetContent(x, y, c, comb, style)
		x += w
	}
}

// Scrolling to a output line number sucks when resizing...
func (v View) Draw(screen tcell.Screen) {
	screen.Clear()
	screenWidth, screenHeight := screen.Size()
	passedLines := 0
	wrapWidth := 72
	if screenWidth < wrapWidth {
		wrapWidth = screenWidth
	}
	for _, line := range v.Page {
		for _, outLine := range wrapLine(line, wrapWidth) {
			startScreenRow := passedLines - v.CurrentLine // FIXME?!
			emitStr(screen, 0, startScreenRow, tcell.StyleDefault, outLine)
			passedLines++
		}
		if passedLines >= v.CurrentLine+screenHeight {
			break // Don't try to draw below-screen content.
		}
	}
	screen.Show()
}

func wrapLine(in gmi.Line, width int) []string {
	switch in.(type) {
	case gmi.TextLine:
		return gwrap.Wrap(in.String(), width)
	case gmi.PreformattedLine:
		return []string{in.String()}
	case gmi.H1Line:
		return gwrap.WrapWithPrefixAndIndent(in.String(), "# ", width)
	case gmi.H2Line:
		return gwrap.WrapWithPrefixAndIndent(in.String(), "## ", width)
	case gmi.H3Line:
		return gwrap.WrapWithPrefixAndIndent(in.String(), "### ", width)
	case gmi.LinkLine:
		return gwrap.WrapWithPrefixAndIndent(in.String(), "=> ", width)
	case gmi.ListLine:
		return gwrap.WrapWithPrefixAndIndent(in.String(), "* ", width)
	case gmi.QuoteLine:
		return gwrap.WrapWithPrefix(in.String(), "> ", width)
	}
	panic(fmt.Errorf("unknown Line type"))
}
