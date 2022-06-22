package secret

import (
	"bytes"
	"encoding/base64"
	"errors"
	"io"

	"github.com/gtank/cryptopasta"
)

type Crypto struct {
	encryptionKey *[32]byte
}

func New(encryptionKey string) (*Crypto, error) {
	secretKey := &[32]byte{}
	if len(encryptionKey) < 32 {
		return nil, errors.New("random hash should be 32 chars in length")
	}
	_, err := io.ReadFull(bytes.NewBufferString(encryptionKey), secretKey[:])
	if err != nil {
		return nil, err
	}

	return &Crypto{
		encryptionKey: secretKey,
	}, nil
}

func (sec *Crypto) Encrypt(str string) (string, error) {
	cipher, err := cryptopasta.Encrypt([]byte(str), sec.encryptionKey)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(cipher), nil
}

func (sec *Crypto) Decrypt(str string) (string, error) {
	encrypted, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return "", err
	}
	decryptedToken, err := cryptopasta.Decrypt(encrypted, sec.encryptionKey)
	if err != nil {
		return "", err
	}
	return string(decryptedToken), nil
}