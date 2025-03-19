package configutils

import (
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// CA:
// openssl ecparam -out cakey.pem -name secp256r1 -genkey
// openssl req -new -x509 -days 3650 -key cakey.pem -sha256 -out ca.pem
// Client:
// openssl ecparam -out key.pem -name secp256r1 -genkey
// openssl req -new -key key.pem -out csr.pem
// openssl x509 -req -days 3650 -sha256 -in csr.pem -CA ca.pem -CAkey cakey.pem -out cert.pem
// Cleanup:
// rm cakey.pem csr.pem
const ca string = `-----BEGIN CERTIFICATE-----
MIIB1jCCAX2gAwIBAgIUHeqn9ag/smJys2AMJyDVQAobXw0wCgYIKoZIzj0EAwIw
QTELMAkGA1UEBhMCU0UxEzARBgNVBAgMClNvbWUtU3RhdGUxEDAOBgNVBAoMB0N1
cGRhdGUxCzAJBgNVBAMMAkNBMB4XDTI1MDMxOTE3MTQzNVoXDTM1MDMxNzE3MTQz
NVowQTELMAkGA1UEBhMCU0UxEzARBgNVBAgMClNvbWUtU3RhdGUxEDAOBgNVBAoM
B0N1cGRhdGUxCzAJBgNVBAMMAkNBMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE
oCVYVx5I/+sOz8se2ZcGZ07CsRwsY9kb8TsS4xM2Z6AQrTvcsZVh28CcL8I2xJNV
dkj2GK8XQkps7GEyv8v7baNTMFEwHQYDVR0OBBYEFMs72Xafqu6EFpRTBT3EmuRa
vppeMB8GA1UdIwQYMBaAFMs72Xafqu6EFpRTBT3EmuRavppeMA8GA1UdEwEB/wQF
MAMBAf8wCgYIKoZIzj0EAwIDRwAwRAIgVtGsCJzcVupejmAx9hi6h7ZKAPZC4Fb5
5etd1q85owwCIDdIR5u7xqdva0ilBjcQmnsHi2LC7iDYZ9Rr4ZEK6Fcz
-----END CERTIFICATE-----`

const cert string = `-----BEGIN CERTIFICATE-----
MIIByTCCAXCgAwIBAgIUeGl/z36HycMF8sgA7jZOsR4WHcAwCgYIKoZIzj0EAwIw
QTELMAkGA1UEBhMCU0UxEzARBgNVBAgMClNvbWUtU3RhdGUxEDAOBgNVBAoMB0N1
cGRhdGUxCzAJBgNVBAMMAkNBMB4XDTI1MDMxOTE3MTcxNloXDTM1MDMxNzE3MTcx
NlowRTELMAkGA1UEBhMCU0UxEzARBgNVBAgMClNvbWUtU3RhdGUxEDAOBgNVBAoM
B0N1cGRhdGUxDzANBgNVBAMMBkNsaWVudDBZMBMGByqGSM49AgEGCCqGSM49AwEH
A0IABIF5YkSZ5/yKH3sKgRhfrl5eEaltih3KPgfmcMuGpTznEeB/SRgiVftWBM6A
5JLobTQ1Svcg3KsWpkULs+ak51GjQjBAMB0GA1UdDgQWBBSEKbBSTW8Htp1u4wPn
7E9ofkksnTAfBgNVHSMEGDAWgBTLO9l2n6ruhBaUUwU9xJrkWr6aXjAKBggqhkjO
PQQDAgNHADBEAiBKfWgK6/LVa+gtnDsykjwxzip0VWgZFmwlrUv2wBe7aQIgcatE
7rEEMp4YbA99DFsvRXerg97oj2DClYdhyluH2n0=
-----END CERTIFICATE-----`

const key string = `-----BEGIN EC PARAMETERS-----
BggqhkjOPQMBBw==
-----END EC PARAMETERS-----
-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIOnlh+ikteP7wkYOVBlx/86WAv5K7QpgFbVyQiR96eX1oAoGCCqGSM49
AwEHoUQDQgAEgXliRJnn/IofewqBGF+uXl4RqW2KHco+B+Zwy4alPOcR4H9JGCJV
+1YEzoDkkuhtNDVK9yDcqxamRQuz5qTnUQ==
-----END EC PRIVATE KEY-----`

func TestLoadTLSConfig(t *testing.T) {
	testCases := []struct {
		Name         string
		BasePath     []string
		FS           fstest.MapFS
		ExpectCA     bool
		ExpectClient bool
	}{
		{
			Name:     "Full config from host",
			BasePath: []string{"192.168.0.116"},
			FS: fstest.MapFS{
				"192.168.0.116/ca.pem":   &fstest.MapFile{Data: []byte(ca)},
				"192.168.0.116/cert.pem": &fstest.MapFile{Data: []byte(cert)},
				"192.168.0.116/key.pem":  &fstest.MapFile{Data: []byte(key)},
			},
			ExpectCA:     true,
			ExpectClient: true,
		},
		{
			Name:     "Only mTLS from host",
			BasePath: []string{"192.168.0.116"},
			FS: fstest.MapFS{
				"192.168.0.116/cert.pem": &fstest.MapFile{Data: []byte(cert)},
				"192.168.0.116/key.pem":  &fstest.MapFile{Data: []byte(key)},
			},
			ExpectCA:     false,
			ExpectClient: true,
		},
		{
			Name:     "Only CA from host",
			BasePath: []string{"192.168.0.116"},
			FS: fstest.MapFS{
				"192.168.0.116/ca.pem": &fstest.MapFile{Data: []byte(ca)},
			},
			ExpectCA:     true,
			ExpectClient: false,
		},
		{
			Name:         "Empty from host",
			BasePath:     []string{"192.168.0.116"},
			FS:           fstest.MapFS{},
			ExpectCA:     false,
			ExpectClient: false,
		},
		{
			Name:     "Full config from fallback",
			BasePath: []string{"192.168.0.116", ""},
			FS: fstest.MapFS{
				"ca.pem":   &fstest.MapFile{Data: []byte(ca)},
				"cert.pem": &fstest.MapFile{Data: []byte(cert)},
				"key.pem":  &fstest.MapFile{Data: []byte(key)},
			},
			ExpectCA:     true,
			ExpectClient: true,
		},
		{
			Name:     "Only mTLS from fallback",
			BasePath: []string{"192.168.0.116", ""},
			FS: fstest.MapFS{
				"cert.pem": &fstest.MapFile{Data: []byte(cert)},
				"key.pem":  &fstest.MapFile{Data: []byte(key)},
			},
			ExpectCA:     false,
			ExpectClient: true,
		},
		{
			Name:     "Only CA from fallback",
			BasePath: []string{"192.168.0.116", ""},
			FS: fstest.MapFS{
				"ca.pem": &fstest.MapFile{Data: []byte(ca)},
			},
			ExpectCA:     true,
			ExpectClient: false,
		},
		{
			Name:         "Empty from fallback",
			BasePath:     []string{"192.168.0.116"},
			FS:           fstest.MapFS{},
			ExpectCA:     false,
			ExpectClient: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			root := t.TempDir()
			err := os.CopyFS(root, testCase.FS)
			require.NoError(t, err)

			basePath := make([]string, len(testCase.BasePath))
			for i, b := range testCase.BasePath {
				basePath[i] = filepath.Join(root, b)
			}

			config, err := LoadTLSConfig(basePath...)
			require.NoError(t, err)

			if testCase.ExpectCA {
				require.NotNil(t, config)
				assert.NotNil(t, config.RootCAs)
			}

			if testCase.ExpectClient {
				require.NotNil(t, config)
				assert.NotNil(t, config.Certificates)
			}

			if !testCase.ExpectCA && !testCase.ExpectClient {
				assert.Nil(t, config)
			}
		})
	}
}
