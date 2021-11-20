package gclient

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"net/url"
)

// A Client is used to request a single page and must be closed after
// that.
type Client struct {
	url        *url.URL
	conn       *tls.Conn
	serverCert *x509.Certificate
}

// Close closes the underlying TLS connection.
func (c *Client) Close() error {
	return c.conn.Close()
}

// Certificate returns the used TLS server certificate.
func (c *Client) ServerCertificate() *x509.Certificate {
	return c.serverCert
}

// NewClient connects to the host of u and returns a new client with
// this connection. If the returned error is nil, the returned Client
// must be closed after use.
func NewClient(u *url.URL) (*Client, error) {
	conf := &tls.Config{InsecureSkipVerify: true}
	conn, err := tls.Dial("tcp", u.Host, conf)
	if err != nil {
		return nil, err
	}
	state := conn.ConnectionState()

	// Only the used leaf certificate is of interest; it has index 0:
	relevantCert := state.PeerCertificates[0]

	return &Client{
		url:        u,
		conn:       conn,
		serverCert: relevantCert,
	}, err
}

// Get requests the resource of c from the connected host. Before doing
// that, c.Certificate() should be checked for validity.
func (c *Client) Get() (io.ReadCloser, error) {
	_, err := c.conn.Write([]byte(c.url.String() + "\r\n"))
	return inputCleaner{c.conn}, err
}

func IsKnownHost(cert *x509.Certificate) (bool, error) {
	// TODO
	return false, nil
}

func AddToKnownHosts(cert *x509.Certificate) error {
	// TODO
	return nil
}

func AddToTemporarilyKnownHosts(cert *x509.Certificate) error {
	// TODO
	return nil
}
