package view

import (
	"fmt"

	"github.com/codesoap/gaia/gmi"
	"github.com/codesoap/gaia/wrap"
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
)

type View struct {
	Screen      tcell.Screen
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
func (v View) Draw() {
	v.Screen.Clear()
	screenWidth, screenHeight := v.Screen.Size()
	passedLines := 0
	wrapWidth := 72
	if screenWidth < wrapWidth {
		wrapWidth = screenWidth
	}
	for _, line := range v.Page {
		for _, outLine := range wrapLine(line, wrapWidth) {
			startScreenRow := passedLines - v.CurrentLine // FIXME?!
			emitStr(v.Screen, 0, startScreenRow, tcell.StyleDefault, outLine)
			passedLines++
		}
		if passedLines >= v.CurrentLine+screenHeight {
			break // Don't try to draw below-screen content.
		}
	}
	v.Screen.Show()
}

func wrapLine(in gmi.Line, width int) []string {
	switch in.(type) {
	case gmi.TextLine:
		return wrap.Wrap(in.String(), width)
	case gmi.PreformattedLine:
		return []string{in.String()}
	case gmi.H1Line:
		return wrap.WrapWithPrefixAndIndent(in.String(), "# ", width)
	case gmi.H2Line:
		return wrap.WrapWithPrefixAndIndent(in.String(), "## ", width)
	case gmi.H3Line:
		return wrap.WrapWithPrefixAndIndent(in.String(), "### ", width)
	case gmi.LinkLine:
		return wrap.WrapWithPrefixAndIndent(in.String(), "=> ", width)
	case gmi.ListLine:
		return wrap.WrapWithPrefixAndIndent(in.String(), "* ", width)
	case gmi.QuoteLine:
		return wrap.WrapWithPrefix(in.String(), "> ", width)
	}
	panic(fmt.Errorf("unknown Line type"))
}
