package receiver

//TODO change directory to mocks and mock struct name
import (
	"bytes"
	"encoding/base64"
	"io"

	"github.com/gtank/cryptopasta"
	"github.com/odpf/siren/plugins/receivers/http"
	"github.com/pkg/errors"
)

var cryptopastaEncryptor = cryptopasta.Encrypt
var cryptopastaDecryptor = cryptopasta.Decrypt

type Transformer interface {
	PreTransform(*Receiver) error
	PostTransform(*Receiver) error
}

//go:generate mockery --name=SlackHelper -r --case underscore --with-expecter --inpackage --filename slack_helper_mock.go --output=./
type SlackHelper interface {
	Transformer
	Encrypt(string) (string, error)
	Decrypt(string) (string, error)
}

//go:generate mockery --name=Exchanger -r --case underscore --with-expecter --inpackage --filename exchanger_mock.go --output=./
type Exchanger interface {
	Exchange(string, string, string) (http.CodeExchangeHTTPResponse, error)
}

type slackHelper struct {
	exchanger     Exchanger
	encryptionKey *[32]byte
}

func NewSlackHelper(exchanger Exchanger, encryptionKey string) (*slackHelper, error) {
	secretKey := &[32]byte{}
	if len(encryptionKey) < 32 {
		return nil, errors.New("random hash should be 32 chars in length")
	}
	_, err := io.ReadFull(bytes.NewBufferString(encryptionKey), secretKey[:])
	if err != nil {
		return nil, err
	}

	return &slackHelper{
		exchanger:     exchanger,
		encryptionKey: secretKey,
	}, nil
}

func (sh *slackHelper) PreTransform(payload *Receiver) error {
	configurations := payload.Configurations
	clientId := configurations["client_id"].(string)
	clientSecret := configurations["client_secret"].(string)
	code := configurations["auth_code"].(string)

	response, err := sh.exchanger.Exchange(code, clientId, clientSecret)
	if err != nil {
		return errors.Wrap(err, "failed to exchange code with slack OAuth server")
	}

	token, err := sh.Encrypt(response.AccessToken)
	if err != nil {
		return errors.Wrap(err, "encryption failed")
	}

	newConfigurations := map[string]interface{}{}
	newConfigurations["workspace"] = response.Team.Name
	newConfigurations["token"] = token
	payload.Configurations = newConfigurations

	return nil
}

func (sh *slackHelper) PostTransform(r *Receiver) error {
	encryptedToken := r.Configurations["token"].(string)
	token, err := sh.Decrypt(encryptedToken)
	if err != nil {
		return errors.Wrap(err, "slackHelper.Decrypt")
	}
	r.Configurations["token"] = token
	return nil
}

func (sh *slackHelper) Encrypt(s string) (string, error) {
	cipher, err := cryptopastaEncryptor([]byte(s), sh.encryptionKey)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(cipher), nil
}

func (sh *slackHelper) Decrypt(s string) (string, error) {
	encrypted, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}
	decryptedToken, err := cryptopastaDecryptor(encrypted, sh.encryptionKey)
	if err != nil {
		return "", err
	}
	return string(decryptedToken), nil
}
