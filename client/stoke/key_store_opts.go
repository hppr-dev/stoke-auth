package stoke

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
)

type PublicKeyStoreOpt func(*BasePublicKeyStore) error

// Configures TLS to use the system certificate pool
func DefaultTLS() PublicKeyStoreOpt {
	return func (s *BasePublicKeyStore) error {
		s.SetTLSConfig(&tls.Config{})
		return nil
	}
}

// Configures TLS to use a specific ca file
func ConfigureTLS(caFile string) PublicKeyStoreOpt {
	return func (s *BasePublicKeyStore) error {
		systemPool := x509.NewCertPool()
		caFileData, err := os.ReadFile(caFile)
		if err != nil {
			return err
		}

		if !systemPool.AppendCertsFromPEM(caFileData) {
			return fmt.Errorf("Could not add %s to certificate pool", caFile)
		}

		s.SetTLSConfig(&tls.Config{
			RootCAs: systemPool,
			ClientCAs: systemPool,
		})
		return nil
	}
}

// Configures TLS to use a certificate/key pair and caFile
func ConfigureMTLS(caFile, certFile, keyFile string) PublicKeyStoreOpt {
	return func (s *BasePublicKeyStore) error {
		keyPair, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return err
		}

		systemPool := x509.NewCertPool()
		caFileData, err := os.ReadFile(caFile)
		if err != nil {
			return err
		}

		if !systemPool.AppendCertsFromPEM(caFileData) {
			return fmt.Errorf("Could not add %s to certificate pool", caFile)
		}

		s.SetTLSConfig(&tls.Config{
			RootCAs: systemPool,
			Certificates: []tls.Certificate{keyPair},
		})
		return nil
	}
}

// Set the http client to use. Used for full control over the http client used to retrieve public keys.
func SetHTTPClient(c *http.Client) PublicKeyStoreOpt {
	return func (s *BasePublicKeyStore) error {
		s.httpClient = c
		return nil
	}
}
