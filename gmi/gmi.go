package gmi

import (
	"regexp"
	"fmt"
	"io"
	"net/url"
)

// FIXME: The specification does not really say whether trailing spaces
// with a missing name are legal, so I'm choosing to be strict.
var linkLineRegexp = regexp.MustCompile(`^=>[ \t]*([^ \t]+)(?:$|[ \t]+([^ \t].*)$)`)

type Page []Line

func NewPage(io.Reader) (Page, error) {
	// TODO: Parse
	return Page{}, nil
}

type Line interface {
	String() string // Returns a human readable representation of the line.
}

type TextLine struct {
	raw string
}

func NewTextLine(raw string) TextLine {
	return TextLine{raw}
}

func (t TextLine) String() string {
	return t.raw
}

type LinkLine struct {
	url *url.URL
	name string
}

func NewLinkLine(raw string) (LinkLine, error) {
	submatches := linkLineRegexp.FindStringSubmatch(raw)
	if len(submatches) != 3 {
		return LinkLine{}, fmt.Errorf("line does not match specification", raw)
	}
	url, err := url.Parse(submatches[1])
	return LinkLine{url: url, name: submatches[2]}, err
}

func (l LinkLine) String() string {
	if l.name == "" {
		return l.url.String()
	}
	if l.url.Scheme != "" && l.url.Scheme != "gemini" {
		return fmt.Sprintf("%s (%s)", l.name, l.url.Scheme)
	}
	return l.name
}

type PreformattedLine struct {
	raw string
}

func (p PreformattedLine) String() string {
	return p.raw
}

type H1Line struct {
	raw string
}

func (h H1Line) String() string {
	return h.raw
}

type H2Line struct {
	raw string
}

func (h H2Line) String() string {
	return h.raw
}

type H3Line struct {
	raw string
}

func (h H3Line) String() string {
	return h.raw
}

type ListLine struct {
	raw string
}

func (l ListLine) String() string {
	return l.raw
}

type QuoteLine struct {
	raw string
}

func (q QuoteLine) String() string {
	return q.raw
}
