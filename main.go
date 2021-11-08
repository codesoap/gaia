package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"os"
)

func main() {
	conf := &tls.Config{
		InsecureSkipVerify:    true,
		VerifyPeerCertificate: verifyServersCert,
	}
	conn, err := tls.Dial("tcp", "gemini.circumlunar.space:1965", conf)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	if _, err = conn.Write([]byte("gemini://gemini.circumlunar.space/docs/specification.gmi\r\n")); err != nil {
		panic(err)
	}
	if _, err = io.Copy(os.Stdout, conn); err != nil {
		panic(err)
	}
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
