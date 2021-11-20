package gclient

import "io"

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
