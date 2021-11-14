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
	// TODO: Remove tabs on non-preformatted?!
	p := Page{}
	for s.Scan() {
		raw := s.Text()
		switch {
		case len(raw) >= 2 && raw[:2] == "=>":
			l, err := NewLinkLine(raw)
			if err != nil {
				return p, err
			}
			p = append(p, l)
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
		case len(raw) >= 1 && raw[:1] == "#":
			if len(raw) >= 4 && raw[:4] == "####" {
				return p, fmt.Errorf("encountered heading with more than three '#'")
			}
			p = append(p, HeadingLine{raw})
		case len(raw) >= 2 && raw[:2] == "* ":
			p = append(p, ListLine{raw})
		case len(raw) >= 1 && raw[:1] == ">":
			p = append(p, QuoteLine{raw})
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
	url  *url.URL
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
		return fmt.Sprintf("=> [%s] %s", l.url.Scheme, l.name)
	}
	return l.name
}

type TextLine struct{ raw string }
type PreformattedLine struct{ raw string }
type HeadingLine struct{ raw string }
type ListLine struct{ raw string }
type QuoteLine struct{ raw string }

func (t TextLine) String() string         { return t.raw }
func (q QuoteLine) String() string        { return q.raw }
func (p PreformattedLine) String() string { return p.raw }
func (h HeadingLine) String() string      { return h.raw }
func (l ListLine) String() string         { return l.raw }
