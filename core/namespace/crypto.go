package namespace

import "github.com/goto/siren/pkg/secret"

type Encryptor interface {
	Encrypt(str secret.MaskableString) (secret.MaskableString, error)
	Decrypt(str secret.MaskableString) (secret.MaskableString, error)
}
