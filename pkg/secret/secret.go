package secret

import (
	"bytes"
	"encoding/base64"
	"io"

	"github.com/gtank/cryptopasta"
	"github.com/odpf/siren/pkg/errors"
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

func (sec *Crypto) Encrypt(str MaskableString) (MaskableString, error) {
	cipher, err := cryptopasta.Encrypt([]byte(str), sec.encryptionKey)
	if err != nil {
		return "", err
	}
	return MaskableString(base64.StdEncoding.EncodeToString(cipher)), nil
}

func (sec *Crypto) Decrypt(str MaskableString) (MaskableString, error) {
	encrypted, err := base64.StdEncoding.DecodeString(str.UnmaskedString())
	if err != nil {
		return "", err
	}
	decryptedToken, err := cryptopasta.Decrypt(encrypted, sec.encryptionKey)
	if err != nil {
		return "", err
	}
	return MaskableString(decryptedToken), nil
}
