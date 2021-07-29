package codeexchange

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/gtank/cryptopasta"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// Repository talks to the store to read or insert data
type Repository struct {
	db            *gorm.DB
	encryptionKey *[32]byte
}

var cryptopastaEncryptor = cryptopasta.Encrypt

func encryptToken(accessToken string, encryptionKey *[32]byte) (string, error) {
	cipher, err := cryptopastaEncryptor([]byte(accessToken), encryptionKey)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(cipher), nil
}

var cryptopastaDecryptor = cryptopasta.Decrypt

func decryptToken(accessToken string, encryptionKey *[32]byte) (string, error) {
	encrypted, err := base64.StdEncoding.DecodeString(accessToken)
	if err != nil {
		return "", err
	}
	decryptedToken, err := cryptopastaDecryptor(encrypted, encryptionKey)
	if err != nil {
		return "", err
	}
	return string(decryptedToken), nil
}

// NewRepository returns repository struct
func NewRepository(db *gorm.DB, encryptionKey string) (*Repository, error) {
	secretKey := &[32]byte{}
	if len(encryptionKey) < 32 {
		return nil, errors.New("random hash should be 32 chars in length")
	}
	_, err := io.ReadFull(bytes.NewBufferString(encryptionKey), secretKey[:])
	if err != nil {
		return nil, err
	}
	return &Repository{db, secretKey}, nil
}

func (r Repository) Upsert(accessToken *AccessToken) error {
	var existingAccessToken AccessToken
	result := r.db.Where(fmt.Sprintf("workspace = '%s'", accessToken.Workspace)).Find(&existingAccessToken)
	if result.Error != nil {
		return result.Error
	}

	token, err := encryptToken(accessToken.AccessToken, r.encryptionKey)
	accessToken.AccessToken = token
	if err != nil {
		return errors.Wrap(err, "encryption failed")
	}

	if result.RowsAffected == 0 {
		result = r.db.Create(accessToken)
	} else {
		result = r.db.Where("id = ?", existingAccessToken.ID).Updates(accessToken)
	}
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r Repository) Get(workspace string) (string, error) {
	var accessToken AccessToken
	result := r.db.Where(fmt.Sprintf("workspace = '%s'", workspace)).Find(&accessToken)
	if result.Error != nil {
		return "", errors.Wrap(result.Error, "search query failed")
	}
	if result.RowsAffected == 0 {
		return "", errors.New(fmt.Sprintf("workspace not found: %s", workspace))
	}
	decryptedAccessToken, err := decryptToken(accessToken.AccessToken, r.encryptionKey)
	if err != nil {
		return "", errors.Wrap(err, "failed to decrypt token")
	}
	return decryptedAccessToken, nil
}

func (r Repository) Migrate() error {
	err := r.db.AutoMigrate(&AccessToken{})
	if err != nil {
		return err
	}
	return nil
}
