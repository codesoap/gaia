package main

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"

	"github.com/codesoap/gaia/gmi"
	"github.com/codesoap/gaia/view"
	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/encoding"
)

type state struct {
	screen      tcell.Screen
	url *url.URL
	view view.View
}

type inputCleaner struct {
	io.ReadCloser
}

func (c inputCleaner) Read(p []byte) (n int, err error) {
	// Replacing with '�' instead of '?' would be nice, but woud be a lot
	// more complicated, since '�' consists of two bytes.

	origP := make([]byte, len(p))
	n, err = c.ReadCloser.Read(origP)
	for i := 0; i < n; i++ {
		if isUnwantedControlChar(origP[i]) {
			p[i] = '?'
		} else {
			p[i] = origP[i]
		}
	}
	return
}

func isUnwantedControlChar(b byte) bool {
	return b != '\t' && b != '\r' && b != '\n' && b < 32
}

func main() {
	u, err := getURL()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not get URL: %v\n", err)
		os.Exit(1)
	}

	encoding.Register()
	screen, err := tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not start tcell: %v\n", err)
		os.Exit(2)
	}
	if err = screen.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Could not start tcell: %v\n", err)
		os.Exit(2)
	}

	s := state{screen: screen, url: u}
	if err = s.loadURL(); err != nil {
		screen.Fini()
		fmt.Fprintf(os.Stderr, "Could not open URL: %v\n", err)
		os.Exit(3)
	}
	if err = s.runEventLoop(); err != nil {
		screen.Fini()
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(4)
	}
	screen.Fini()
	fmt.Println(s.url) // TODO: Cut out port if standard.
}

func (s *state) runEventLoop() error {
	for {
		switch ev := s.screen.PollEvent().(type) {
		case *tcell.EventResize:
			s.screen.Sync()
			s.view.Draw(s.screen)
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape {
				return nil
			} else if ev.Key() == tcell.KeyRune && ev.Rune() == 'j' {
				s.view.CurrentLine++
				s.view.Draw(s.screen)
			} else if ev.Key() == tcell.KeyRune && ev.Rune() == 'k' {
				s.view.CurrentLine--
				s.view.Draw(s.screen)
			}
		}
	}
}

func (s *state) loadURL() error {
	conn, err := get(s.url)
	if err != nil {
		return err
	}
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	if !scanner.Scan() {
		return fmt.Errorf("could not read header")
	}
	headerSplit := strings.SplitN(scanner.Text(), " ", 2)
	if len(headerSplit) != 2 || len(headerSplit[0]) != 2 || headerSplit[1] == "" {
		return fmt.Errorf("invalid header")
	}
	switch headerSplit[0][:1] {
	case "1":
		return fmt.Errorf("TODO: implement INPUT")
	case "2":
		page, err := gmi.ParsePage(scanner)
		if err != nil {
			return err
		}
		s.view = view.View{page, 0}
	case "3":
		return fmt.Errorf("TODO: implement REDIRECT")
	case "4":
		return fmt.Errorf("TODO: implement TMP FAIL")
	case "5":
		return fmt.Errorf("TODO: implement PERM FAIL")
	case "6":
		return fmt.Errorf("TODO: implement CERT REQUIRED")
	}
	return nil
}

func get(u *url.URL) (io.ReadCloser, error) {
	conf := &tls.Config{
		InsecureSkipVerify:    true,
		VerifyPeerCertificate: verifyServersCert,
	}
	conn, err := tls.Dial("tcp", u.Host, conf)
	if err != nil {
		return nil, err
	}
	if _, err = conn.Write([]byte(u.String() + "\r\n")); err != nil {
		conn.Close()
	}
	return inputCleaner{conn}, err
}

func getURL() (*url.URL, error) {
	if len(os.Args) != 2 {
		return &url.URL{}, fmt.Errorf("wrong argument count") // TODO: help text
	}
	u, err := url.Parse(os.Args[1])
	if err != nil {
		return &url.URL{}, err
	} else if u.Host == "" {
		u, err = url.Parse("//" + os.Args[1])
		if err != nil {
			return &url.URL{}, err
		}
	}
	if u.Scheme == "" {
		u.Scheme = "gemini"
	} else if u.Scheme != "gemini" {
		return &url.URL{}, fmt.Errorf("not a gemini URL")
	}
	if u.Port() == "" {
		u.Host += ":1965"
	}
	return u, nil
}

func verifyServersCert(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
	// for _, rawCert := range rawCerts {
	// 	cert, err := x509.ParseCertificate(rawCert)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	return nil
	// }
	// return fmt.Errorf("could not find trusted certificate")
	return nil
}
