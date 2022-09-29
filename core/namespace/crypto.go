package namespace

import "github.com/odpf/siren/pkg/secret"

//go:generate mockery --name=Encryptor -r --case underscore --with-expecter --structname Encryptor --filename encryptor.go --output=./mocks
type Encryptor interface {
	Encrypt(str secret.MaskableString) (secret.MaskableString, error)
	Decrypt(str secret.MaskableString) (secret.MaskableString, error)
}
