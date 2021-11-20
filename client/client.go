package client

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"net/url"
)

// TODO: AddKnownServerCertificate

// A Client is used to request a single page and must be closed after
// that.
type Client struct {
	url  *url.URL
	conn *tls.Conn
	cert *x509.Certificate
}

func (c *Client) Close() error {
	return c.conn.Close()
}

// NewClient connects to the host of u and returns a new client with
// this connection.
func NewClient(u *url.URL) (*Client, error) {
	conf := &tls.Config{InsecureSkipVerify: true}
	conn, err := tls.Dial("tcp", u.Host, conf)
	if err != nil {
		return nil, err
	}
	state := conn.ConnectionState()

	// Only the used leaf certificate is of interest; it has index 0:
	relevantCert := state.PeerCertificates[0]

	return &Client{url: u, conn: conn, cert: relevantCert}, err
}

// Get requests the resource of c from the connected host. Before doing
// that, c.Certificate() should be checked for validity.
func (c *Client) Get() (io.ReadCloser, error) {
	_, err := c.conn.Write([]byte(c.url.String() + "\r\n"))
	return inputCleaner{c.conn}, err
}
