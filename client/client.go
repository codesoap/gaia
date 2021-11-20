package client

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net/url"
)

// TODO: AddKnownServerCertificate

// A Client is used to request a single page and must be closed after
// that.
type Client struct {
	url  *url.URL
	conn *tls.Conn
}

func (c *Client) Close() error {
	return c.conn.Close()
}

// NewClient connects to the host of u and returns a new client with
// this connection.
//
// ErrUnknownServerCertificate will be returned, if the servers
// certificate is unknown.
func NewClient(u *url.URL) (*Client, error) {
	conf := &tls.Config{
		InsecureSkipVerify:    true,
		VerifyPeerCertificate: verifyServersCert,
	}
	conn, err := tls.Dial("tcp", u.Host, conf)
	if err != nil {
		return nil, err
	}
	return &Client{url: u, conn: conn}, nil
}

func (c *Client) Get() (io.ReadCloser, error) {
	_, err := c.conn.Write([]byte(c.url.String() + "\r\n"))
	return inputCleaner{c.conn}, err
}

func verifyServersCert(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
	// The first received certificate is used, others are discarded.
	if len(rawCerts) > 0 {
		cert, err := x509.ParseCertificate(rawCerts[0])
		if err != nil {
			return err
		}
		if isKnownCertificate(cert) {
			return nil
		}
		return ErrUnknownServerCertificate{cert}
	}
	return fmt.Errorf("could not find trusted certificate")
}

func isKnownCertificate(cert *x509.Certificate) bool {
	return true // FIXME
}
