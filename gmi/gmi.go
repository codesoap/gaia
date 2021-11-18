package gmi

import (
	"bufio"
	"fmt"
	"net/url"
	"regexp"
)

// FIXME: The specification does not really say whether trailing spaces
// with a missing name are legal, so I'm choosing to be strict.
var linkLineRegexp = regexp.MustCompile(`^=>[ \t]*([^ \t]+)(?:$|[ \t]+([^ \t].*)$)`)

type Page []Line

func ParsePage(s *bufio.Scanner) (Page, error) {
	linkIndex := 0
	p := Page{}
	for s.Scan() {
		raw := s.Text()
		switch {
		case len(raw) >= 2 && raw[:2] == "=>":
			l, err := NewLinkLine(raw, linkIndex)
			if err != nil {
				return p, err
			}
			p = append(p, l)
			linkIndex++
		case len(raw) >= 3 && raw[:3] == "```":
			for s.Scan() {
				raw := s.Text()
				if len(raw) >= 3 && raw[:3] == "```" {
					break
				}
				p = append(p, PreformattedLine{raw})
			}
			if s.Err() != nil {
				return p, s.Err()
			}
		case len(raw) >= 2 && raw[:2] == "# ":
			p = append(p, H1Line{raw[2:]})
		case len(raw) >= 3 && raw[:3] == "## ":
			p = append(p, H2Line{raw[3:]})
		case len(raw) >= 4 && raw[:4] == "### ":
			p = append(p, H3Line{raw[4:]})
		case len(raw) >= 2 && raw[:2] == "* ":
			p = append(p, ListLine{raw[2:]})
		case len(raw) >= 1 && raw[:1] == ">":
			p = append(p, QuoteLine{raw[1:]})
		default:
			p = append(p, TextLine{raw})
		}
	}
	return p, s.Err()
}

type Line interface {
	String() string // Returns a human readable representation of the line.
}

type LinkLine struct {
	url   *url.URL
	name  string
	index int
}

func NewLinkLine(raw string, index int) (LinkLine, error) {
	submatches := linkLineRegexp.FindStringSubmatch(raw)
	if len(submatches) != 3 {
		return LinkLine{}, fmt.Errorf("line does not match specification", raw)
	}
	url, err := url.Parse(submatches[1])
	return LinkLine{url: url, name: submatches[2], index: index}, err
}

func (l LinkLine) String() string {
	if l.name == "" {
		return l.url.String()
	}
	if l.url.Scheme != "" && l.url.Scheme != "gemini" {
		return fmt.Sprintf("#%d [%s] %s", l.index, l.url.Scheme, l.name)
	}
	return fmt.Sprintf("#%d %s", l.index, l.name)
}

type TextLine struct{ raw string }
type PreformattedLine struct{ raw string }
type H1Line struct{ raw string }
type H2Line struct{ raw string }
type H3Line struct{ raw string }
type ListLine struct{ raw string }
type QuoteLine struct{ raw string }

func (t TextLine) String() string         { return t.raw }
func (q QuoteLine) String() string        { return q.raw }
func (p PreformattedLine) String() string { return p.raw }
func (h H1Line) String() string           { return h.raw }
func (h H2Line) String() string           { return h.raw }
func (h H3Line) String() string           { return h.raw }
func (l ListLine) String() string         { return l.raw }
