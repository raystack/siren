package namespace

//go:generate mockery --name=Encryptor -r --case underscore --with-expecter --structname Encryptor --filename encryptor.go --output=./mocks
type Encryptor interface {
	Encrypt(str string) (string, error)
	Decrypt(str string) (string, error)
}
