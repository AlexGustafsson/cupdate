package webpush

import (
	"context"
	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"time"

	"vendor/golang.org/x/crypto/hkdf"

	"github.com/AlexGustafsson/cupdate/internal/webpush/vapid"
)

// Client provides methods of pushing a message to a Web Push service.
// SEE: Generic Event Delivery Using HTTP Push - https://datatracker.ietf.org/doc/html/rfc8030.
// SEE: VAPID - https://datatracker.ietf.org/doc/html/rfc8292#section-3.2.
// SEE: Message Encryption for Web Push - https://www.rfc-editor.org/rfc/rfc8291.html
type Client struct {
	Endpoint                    string
	AuthenticationSecret        []byte
	UserAgentPublicKey          *ecdh.PublicKey
	ApplicationServerPrivateKey *ecdsa.PrivateKey
}

type PushOptions struct {
	TTL         int64
	ContentType string
	Urgency     string
	Topic       string
}

func (c *Client) Push(ctx context.Context, content []byte, options *PushOptions) error {
	vapidToken, err := vapid.NewToken("https://cupdate.home.local", time.Now().Add(5*time.Minute), "https://cupdate.home.local", c.ApplicationServerPrivateKey)
	if err != nil {
		return err
	}

	privateKey, err := ecdh.P256().GenerateKey(rand.Reader)
	if err != nil {
		return err
	}

	// "WebPush: info" || 0x00 || ua_public || as_public
	info.WriteString("Content-Encoding: aes128gcm")
	info.WriteRune(0x00)
	info.WriteRune(0x01)

	hkdf := hkdf.New(sha256.New, sharedSecret, c.AuthenticationSecret, info.Bytes())

	var ikm [32]byte
	hkdf.Read(ikm[:])

	var salt [16]byte
	rand.Read(salt[:])

	req, err := http.NewRequest(http.MethodPost, c.Endpoint, nil)
	if err != nil {
		return err
	}

	req.Header.Set("TTL", "30")
	req.Header.Set("Content-Encoding", "aes128gcm")
	req.Header.Set("Crypto-Key", "dh="+base64.RawURLEncoding.EncodeToString(privateKey.PublicKey().Bytes()))
	req.Header.Set("Authorization", vapid.FormatAuthorizationHeader(vapidToken, c.ApplicationServerPrivateKey.PublicKey))

	return nil
}
