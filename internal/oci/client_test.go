package oci

import (
	"bytes"
	"context"
	"crypto"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/AlexGustafsson/cupdate/internal/cachetest"
	"github.com/AlexGustafsson/cupdate/internal/httputil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestClientGetManifest(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	client := &Client{
		Client: httputil.NewClient(cachetest.NewCache(t), 24*time.Hour),
	}

	references := []string{
		"k8s.gcr.io/pause",
		"quay.io/jetstack/cert-manager-startupapicheck:v1.16.2",
		"registry.k8s.io/kube-state-metrics/kube-state-metrics:v2.14.0",
		"gcr.io/zenika-hub/alpine-chrome:123",
		"postgres:12-alpine",
		"ghcr.io/jmbannon/ytdl-sub:2024.10.09",
		"registry.gitlab.com/arm-research/smarter/smarter-device-manager",
	}

	for _, reference := range references {
		t.Run(reference, func(t *testing.T) {
			ref, err := ParseReference(reference)
			require.NoError(t, err)

			manifest, err := client.GetManifest(context.TODO(), ref)
			require.NoError(t, err)
			fmt.Printf("%+v\n", manifest)

			// Rewrite ref to pin digest of manifst
			ref.HasTag = false
			ref.Tag = ""
			ref.HasDigest = true
			switch m := manifest.(type) {
			case *ImageManifest:
				ref.Digest = m.Digest
			case *ImageIndex:
				ref.Digest = m.Digest
			}

			// Expect it to exist
			manifest, err = client.GetManifest(context.TODO(), ref)
			require.NoError(t, err)
			fmt.Printf("%+v\n", manifest)
		})
	}
}

func TestClientGetAttestationManifest(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	client := &Client{
		Client: httputil.NewClient(cachetest.NewCache(t), 24*time.Hour),
	}

	references := []string{
		"ghcr.io/alexgustafsson/cupdate",
		"mongo:6.0.20",
	}

	for _, reference := range references {
		t.Run(reference, func(t *testing.T) {
			ref, err := ParseReference(reference)
			require.NoError(t, err)

			manifest, err := client.GetManifest(context.TODO(), ref)
			require.NoError(t, err)

			index, ok := manifest.(*ImageIndex)
			require.True(t, ok)

			for manifestDigest, attestationManifestDigest := range index.AttestationManifestDigest() {
				fmt.Println("Getting manifest", manifestDigest)
				fmt.Println("Getting attestation", attestationManifestDigest)

				attestationManifest, err := client.GetAttestationManifest(context.TODO(), ref, attestationManifestDigest)
				require.NoError(t, err)

				fmt.Printf("%+v\n", attestationManifest)

				provenancePredicateType, provenanceDigest, ok := attestationManifest.ProvenanceDigest()
				require.True(t, ok)
				fmt.Println(provenancePredicateType, provenanceDigest)

				blob, err := client.GetBlob(context.TODO(), ref, provenanceDigest, false)
				require.NoError(t, err)
				io.Copy(os.Stdout, io.LimitReader(blob, 1024))
				fmt.Println()
				blob.Close()

				sbomPredicateType, sbomDigest, ok := attestationManifest.SBOMDigest()
				require.True(t, ok)
				fmt.Println(sbomPredicateType, sbomDigest)

				blob, err = client.GetBlob(context.TODO(), ref, sbomDigest, false)
				require.NoError(t, err)
				io.Copy(os.Stdout, io.LimitReader(blob, 1024))
				fmt.Println()
				blob.Close()
			}
		})
	}
}

func TestClientHeadBlob(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	client := &Client{
		Client: httputil.NewClient(cachetest.NewCache(t), 24*time.Hour),
	}

	testCases := []struct {
		Reference string
		Digest    string
	}{
		{
			Reference: "k8s.gcr.io/pause",
			Digest:    "sha256:350b164e7ae1dcddeffadd65c76226c9b6dc5553f5179153fb0e36b78f2a5e06",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Reference, func(t *testing.T) {
			ref, err := ParseReference(testCase.Reference)
			require.NoError(t, err)

			info, err := client.HeadBlob(context.TODO(), ref, testCase.Digest)
			require.NoError(t, err)
			fmt.Printf("%+v\n", info)
		})
	}
}

func TestClientGetAnnotations(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	client := &Client{
		Client: httputil.NewClient(cachetest.NewCache(t), 24*time.Hour),
	}

	references := []string{
		"quay.io/jetstack/cert-manager-startupapicheck:v1.16.2",
		"homeassistant/home-assistant",
		"ghcr.io/jmbannon/ytdl-sub",
	}

	for _, reference := range references {
		t.Run(reference, func(t *testing.T) {
			ref, err := ParseReference(reference)
			require.NoError(t, err)

			annotations, err := client.GetAnnotations(context.TODO(), ref, nil)
			require.NoError(t, err)
			assert.NotNil(t, annotations)
			fmt.Printf("%+v\n", annotations)
		})
	}
}

func TestClientGetTags(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	client := &Client{
		Client: httputil.NewClient(cachetest.NewCache(t), 24*time.Hour),
	}

	references := []string{
		"k8s.gcr.io/pause",
		"quay.io/jetstack/cert-manager-startupapicheck",
		"registry.k8s.io/kube-state-metrics/kube-state-metrics",
		"gcr.io/zenika-hub/alpine-chrome",
		"mongo",
	}

	for _, reference := range references {
		t.Run(reference, func(t *testing.T) {
			ref, err := ParseReference(reference)
			require.NoError(t, err)

			tags, err := client.GetTags(context.TODO(), ref, nil)
			require.NoError(t, err)
			assert.NotNil(t, tags)
			fmt.Printf("%+v\n", tags)
		})
	}
}

func TestIntegrationClientZotBasicAuth(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	zotConfig := `{
  "distSpecVersion": "1.0.1",
  "storage": {
    "rootDirectory": "/tmp/zot/storage"
  },
  "log": {
    "level": "debug"
  },
  "http": {
    "address": "0.0.0.0",
    "port": "9090",
    "auth": {
      "htpasswd": {
        "path": "/etc/zot/htpasswd"
      },
      "apikey": true
    }
  }
}
`
	// htpasswd -bBn username password
	htpasswd := `username:$2y$05$uAaTn0C6bgrELr8oDCGkaurrdjctUCnb2LHl7M6GMhRmw45fSJEHe`

	zotRequest := testcontainers.ContainerRequest{
		Image:        "ghcr.io/project-zot/zot:v2.1.10",
		ExposedPorts: []string{"9090/tcp"},
		Files: []testcontainers.ContainerFile{
			{
				Reader:            strings.NewReader(zotConfig),
				ContainerFilePath: "/etc/zot/config.json",
			},
			{
				Reader:            strings.NewReader(htpasswd),
				ContainerFilePath: "/etc/zot/htpasswd",
			},
		},
		WaitingFor: wait.ForHTTP("/v2/").WithPort("9090").WithStatusCodeMatcher(func(status int) bool { return true }),
		LogConsumerCfg: &testcontainers.LogConsumerConfig{
			Consumers: []testcontainers.LogConsumer{
				&testcontainers.StdoutLogConsumer{},
			},
		},
	}

	zot, err := testcontainers.GenericContainer(context.TODO(), testcontainers.GenericContainerRequest{
		ContainerRequest: zotRequest,
		Started:          true,
	})
	testcontainers.CleanupContainer(t, zot)
	require.NoError(t, err)

	zotHost, err := zot.Host(context.TODO())
	require.NoError(t, err)

	zotPort, err := zot.MappedPort(context.TODO(), "9090/tcp")
	require.NoError(t, err)

	authMux := httputil.NewAuthMux()
	authMux.Handle(fmt.Sprintf("%s:%s", zotHost, zotPort.Port()), httputil.BasicAuthHandler{
		Username: "username",
		Password: "password",
	})

	client := &Client{
		Client:   httputil.NewClient(cachetest.NewCache(t), 24*time.Hour),
		AuthFunc: authMux.HandleAuth,
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s:%s/v2/", zotHost, zotPort.Port()), nil)
	require.NoError(t, err)

	res, err := client.Do(req)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestIntegrationClientZotBearerAuth(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	// SEE: https://github.com/project-zot/zot/blob/2a4edde637509ae05da81a8e84c72ebc687e5c0c/pkg/api/bearer.go
	// SEE: https://github.com/project-zot/zot/blob/2a4edde637509ae05da81a8e84c72ebc687e5c0c/pkg/api/authn.go#L404
	// SEE: https://github.com/project-zot/zot/blob/2a4edde637509ae05da81a8e84c72ebc687e5c0c/pkg/api/authn.go#L905
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	require.NoError(t, err)

	template := x509.Certificate{
		SerialNumber: big.NewInt(0),
		Subject: pkix.Name{
			Organization: []string{"Localhost IDP"},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(1 * time.Hour),

		KeyUsage:              x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,

		DNSNames: []string{"localhost"},
	}

	der, err := x509.CreateCertificate(rand.Reader, &template, &template, publicKey, privateKey)
	require.NoError(t, err)

	certificate := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: der,
	})

	authServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// username:password
		if assert.Equal(t, "Basic dXNlcm5hbWU6cGFzc3dvcmQ=", r.Header.Get("Authorization")) {
			var result struct {
				Token     string    `json:"token"`
				ExpiresIn int       `json:"expires_in"`
				IssuedAt  time.Time `json:"issued_at"`
			}

			header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"EdDSA","typ":"JWT"}`))
			payload := base64.RawURLEncoding.EncodeToString([]byte(fmt.Sprintf(`{"iat":%d,"access":[{"type":"registry","name":"catalog","action":"*"}]}`, time.Now().Unix())))

			jwt := header + "." + payload

			signature, err := privateKey.Sign(rand.Reader, []byte(jwt), crypto.Hash(0))
			require.NoError(t, err)
			jwt += "." + base64.RawURLEncoding.EncodeToString(signature)

			result.Token = jwt
			result.ExpiresIn = int(time.Now().Add(1 * time.Hour).Unix())
			result.IssuedAt = time.Now()

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(&result)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}
	}))
	defer authServer.Close()

	zotConfig := fmt.Sprintf(`{
  "distSpecVersion": "1.0.1",
  "storage": {
    "rootDirectory": "/tmp/zot/storage"
  },
  "log": {
    "level": "debug"
  },
  "http": {
    "address": "0.0.0.0",
    "port": "9090",
    "auth": {
      "bearer": {
        "realm": "%s",
        "service": "zot",
        "cert": "/etc/zot/auth-service.crt"
      },
      "apikey": true
    }
  }
}
`, authServer.URL)

	zotRequest := testcontainers.ContainerRequest{
		Image:        "ghcr.io/project-zot/zot:v2.1.10",
		ExposedPorts: []string{"9090/tcp"},
		Files: []testcontainers.ContainerFile{
			{
				Reader:            strings.NewReader(zotConfig),
				ContainerFilePath: "/etc/zot/config.json",
			},
			{
				Reader:            bytes.NewReader(certificate),
				ContainerFilePath: "/etc/zot/auth-service.crt",
			},
		},
		WaitingFor: wait.ForHTTP("/v2/").WithPort("9090").WithStatusCodeMatcher(func(status int) bool { return true }),
		LogConsumerCfg: &testcontainers.LogConsumerConfig{
			Consumers: []testcontainers.LogConsumer{
				&testcontainers.StdoutLogConsumer{},
			},
		},
	}

	zot, err := testcontainers.GenericContainer(context.TODO(), testcontainers.GenericContainerRequest{
		ContainerRequest: zotRequest,
		Started:          true,
	})
	testcontainers.CleanupContainer(t, zot)
	require.NoError(t, err)

	zotHost, err := zot.Host(context.TODO())
	require.NoError(t, err)

	zotPort, err := zot.MappedPort(context.TODO(), "9090/tcp")
	require.NoError(t, err)

	authMux := httputil.NewAuthMux()
	authMux.Handle(authServer.URL, httputil.BasicAuthHandler{
		Username: "username",
		Password: "password",
	})

	client := &Client{
		Client:   httputil.NewClient(cachetest.NewCache(t), 24*time.Hour),
		AuthFunc: authMux.HandleAuth,
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s:%s/v2/", zotHost, zotPort.Port()), nil)
	require.NoError(t, err)

	res, err := client.Do(req)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
}
