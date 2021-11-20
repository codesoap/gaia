package client

import (
	"crypto/x509"
	"fmt"
)

type ErrUnknownServerCertificate struct {
	cert *x509.Certificate
	// TODO: Fingerprint
}

func (e ErrUnknownServerCertificate) Error() string {
	return fmt.Sprintf("unknown server certificate")
	// TODO: Fingerprint
}
