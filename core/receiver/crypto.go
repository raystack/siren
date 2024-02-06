package receiver

type Encryptor interface {
	Encrypt(str string) (string, error)
	Decrypt(str string) (string, error)
}
