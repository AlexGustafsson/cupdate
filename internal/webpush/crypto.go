package webpush

import (
	"bytes"
	"crypto/ecdh"
	"crypto/hmac"
	"crypto/sha256"
	"vendor/golang.org/x/crypto/hkdf"
)

func DeriveInputKeyingMaterial(senderPrivateKey *ecdh.PrivateKey, recipientPublicKey *ecdh.PublicKey, authenticationSecret []byte) ([]byte, error) {
	// SEE: https://www.rfc-editor.org/rfc/rfc8291.html#section-3.4
	sharedSecret, err := senderPrivateKey.ECDH(recipientPublicKey)
	if err != nil {
		return nil, err
	}
	prkKey := hkdf.Extract(sha256.New, sharedSecret, authenticationSecret)

	hkdf.New()

	// "WebPush: info" || 0x00 || ua_public || as_public
	var info bytes.Buffer
	info.WriteString("WebPush: info")
	info.WriteRune(0x00)
	info.Write(recipientPublicKey.Bytes())
	info.Write(senderPrivateKey.PublicKey().Bytes())
	info.WriteRune(0x01)

	return hmac.New(sha256.New, prkKey).Sum(info.Bytes()), nil
}

func DeriveContentEncryptionKey(ikm []byte, salt []byte) []byte {
	prk := hkdf.Extract(sha256.New, ikm, salt)

	// "Content-Encoding: aes128gcm" || 0x00 || 0x01
	var info bytes.Buffer
	info.WriteString("Content-Encoding: aes128gcm")
	info.WriteRune(0x00)
	info.WriteRune(0x01)

	cek := hmac.New(sha256.New, prk).Sum(info.Bytes())

	hkdf.Expand(sha256.New, prk)

	hkdf := hkdf.New(sha256.New, ikm, salt, info.Bytes())

	var cek [32]byte
	hkdf.Read(cek[:])

	return cek[:]
}

func DeriveNonce(prk []byte) {
	// "Content-Encoding: nonce" || 0x00 || 0x01
	var info bytes.Buffer
	info.WriteString("Content-Encoding: nonce")
	info.WriteRune(0x00)
	info.WriteRune(0x01)

	hkdf := hkdf.New(sha256.New, prk, salt, info.Bytes())

	var cek [32]byte
	hkdf.Read(cek[:])

	return cek[:]
}
