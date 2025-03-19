package configutils

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// LoadTLSConfig reads a TLS config from the first directory containing one or
// both of the following files:
//
//   - ca.pem
//   - cert.pem + key.pem
//
// If none are found, a nil config is returned.
func LoadTLSConfig(basePath ...string) (*tls.Config, error) {
	if len(basePath) == 0 {
		return nil, nil
	}

	if len(basePath) > 1 {
		for _, basePath := range basePath {
			tlsConfig, err := LoadTLSConfig(basePath)
			if err != nil {
				return nil, err
			}

			if tlsConfig != nil {
				return tlsConfig, nil
			}
		}
	}

	b := basePath[0]

	config := &tls.Config{}

	// Load the client certificate and key (if any)
	{
		certificate, err := tls.LoadX509KeyPair(filepath.Join(b, "cert.pem"), filepath.Join(b, "key.pem"))
		if errors.Is(err, os.ErrNotExist) {
			// Do nothing
		} else if err != nil {
			return nil, err
		} else {
			config.Certificates = []tls.Certificate{certificate}
		}
	}

	// Load the certificate authority (if any)
	{
		caFile, err := os.ReadFile(filepath.Join(b, "ca.pem"))
		if errors.Is(err, os.ErrNotExist) {
			// Do nothing
		} else if err != nil {
			return nil, err
		} else {
			rootCAs := x509.NewCertPool()
			if ok := rootCAs.AppendCertsFromPEM(caFile); !ok {
				return nil, fmt.Errorf("invalid CA PEM file")
			}
			config.RootCAs = rootCAs
		}
	}

	if config.Certificates == nil && config.RootCAs == nil {
		return nil, nil
	}

	return config, nil
}
