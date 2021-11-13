package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net/url"
	"os"
)

/*
- Read header line. If MIME != text/* or encoding is not UTF-8: do something else
- gmi is handled by the browser itself.
- text is viewn in less?!
- Read the tcp input into a buffer until \n is received.
	- Discard \r if directly before \n.
	- Replace control bytes by '�'.
	- Limit buffer to 10k?!
*/

func main() {
	u, err := getURL()
	if err != nil {
		panic(err)
	}
	conf := &tls.Config{
		InsecureSkipVerify:    true,
		VerifyPeerCertificate: verifyServersCert,
	}
	conn, err := tls.Dial("tcp", u.Host+":1965", conf)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	if _, err = conn.Write([]byte(u.String() + "\r\n")); err != nil {
		panic(err)
	}
	if _, err = io.Copy(os.Stdout, conn); err != nil {
		panic(err)
	}
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
	return u, nil
}

func verifyServersCert(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
	for _, rawCert := range rawCerts {
		cert, err := x509.ParseCertificate(rawCert)
		if err != nil {
			// TODO: Think about this not being a reason to quit.
			return err
		}
		fmt.Println("NotBefore         :", cert.NotBefore.String())
		fmt.Println("NotAfter          :", cert.NotAfter.String())
		fmt.Println("Subject           :", cert.Subject)
		fmt.Println("SignatureAlgorithm:", cert.SignatureAlgorithm)
		return nil
	}
	return fmt.Errorf("could not find trusted certificate")
}
