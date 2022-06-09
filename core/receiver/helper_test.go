package receiver_test

// var cryptopastaEncryptor = cryptopasta.Encrypt
// var cryptopastaDecryptor = cryptopasta.Decrypt

// type SlackHelperTestSuite struct {
// 	suite.Suite
// 	exchangerMock *mocks.Exchanger
// 	key           *[32]byte
// 	keyString     string
// }

// func TestSlackHelper(t *testing.T) {
// 	suite.Run(t, new(SlackHelperTestSuite))
// }

// func (s *SlackHelperTestSuite) SetupTest() {
// 	s.exchangerMock = &mocks.Exchanger{}
// 	s.keyString = "abcdefghijklmnopqrstuvwxyzabcdef"

// 	secretKey := &[32]byte{}
// 	_, err := io.ReadFull(bytes.NewBufferString(s.keyString), secretKey[:])
// 	s.NoError(err)
// 	s.key = secretKey
// }

// func (s *SlackHelperTestSuite) TestSlackHelper_PreTransform() {
// 	configurations := make(map[string]interface{})
// 	configurations["client_id"] = "foo"
// 	configurations["client_secret"] = "bar"
// 	configurations["auth_code"] = "foo"

// 	responseConfigurations := make(map[string]interface{})
// 	responseConfigurations["workspace"] = "test-name"
// 	responseConfigurations["token"] = "rz+RFvt1Uo6mfJ0DsilnyYqeTI2bJD6OcCS14+AhesxN9sSVePr7jc56toOI"
// 	response := &receiver.Receiver{
// 		Configurations: responseConfigurations,
// 	}

// 	codeExchangeHTTPResponse := http.CodeExchangeHTTPResponse{
// 		AccessToken: "test-access-token",
// 		Team: struct {
// 			Name string `json:"name"`
// 		}{
// 			Name: "test-name",
// 		},
// 	}

// 	s.Run("should transform payload on successful code exchange", func() {
// 		slackHelper, err := receiver.NewSlackHelper(s.exchangerMock, s.keyString)
// 		s.NoError(err)

// 		var oldCryptopastaEncryptor = cryptopastaEncryptor
// 		defer func() {
// 			cryptopastaEncryptor = oldCryptopastaEncryptor
// 		}()
// 		cryptopastaEncryptor = func(_ []byte, _ *[32]byte) ([]byte, error) {
// 			return []byte("bar"), nil
// 		}

// 		payload := &receiver.Receiver{
// 			Configurations: configurations,
// 		}
// 		s.exchangerMock.On("Exchange", "foo", "foo", "bar").
// 			Return(codeExchangeHTTPResponse, nil).Once()

// 		err = slackHelper.PreTransform(payload)
// 		s.Nil(err)
// 		s.Equal(payload, response)
// 		s.exchangerMock.AssertCalled(s.T(), "Exchange", "foo", "foo", "bar")
// 	})

// 	s.Run("should return error if code exchange failed", func() {
// 		slackHelper, err := receiver.NewSlackHelper(s.exchangerMock, "abcdefghijklmnopqrstuvwxyzabcdef")
// 		s.NoError(err)

// 		payload := &receiver.Receiver{
// 			Configurations: configurations,
// 		}
// 		s.exchangerMock.On("Exchange", "foo", "foo", "bar").
// 			Return(http.CodeExchangeHTTPResponse{}, errors.New("random error")).Once()

// 		err = slackHelper.PreTransform(payload)
// 		s.EqualError(err, "failed to exchange code with slack OAuth server: random error")
// 	})

// 	s.Run("should return error if access token encryption failed", func() {
// 		slackHelper, err := receiver.NewSlackHelper(s.exchangerMock, "abcdefghijklmnopqrstuvwxyzabcdef")
// 		s.NoError(err)

// 		var oldCryptopastaEncryptor = cryptopastaEncryptor
// 		defer func() {
// 			cryptopastaEncryptor = oldCryptopastaEncryptor
// 		}()
// 		cryptopastaEncryptor = func(_ []byte, _ *[32]byte) ([]byte, error) {
// 			return nil, errors.New("random error")
// 		}
// 		payload := &receiver.Receiver{
// 			Configurations: configurations,
// 		}
// 		s.exchangerMock.On("Exchange", "foo", "foo", "bar").
// 			Return(codeExchangeHTTPResponse, nil).Once()

// 		err = slackHelper.PreTransform(payload)
// 		s.EqualError(err, "encryption failed: random error")
// 	})
// }

// func (s *SlackHelperTestSuite) TestSlackHelper_PostTransform() {

// 	response := &receiver.Receiver{
// 		Configurations: map[string]interface{}{
// 			"token": "test-token",
// 		},
// 	}

// 	s.Run("should transform payload on successful decrypt", func() {
// 		configurations := make(map[string]interface{})
// 		configurations["token"] = "YmFy"
// 		payload := &receiver.Receiver{
// 			Configurations: configurations,
// 		}

// 		slackHelper, err := receiver.NewSlackHelper(s.exchangerMock, "abcdefghijklmnopqrstuvwxyzabcdef")
// 		s.NoError(err)

// 		var oldCryptopastaDecryptor = cryptopastaDecryptor
// 		defer func() {
// 			cryptopastaEncryptor = oldCryptopastaDecryptor
// 		}()
// 		cryptopastaDecryptor = func(_ []byte, _ *[32]byte) ([]byte, error) {
// 			return []byte("test-token"), nil
// 		}

// 		err = slackHelper.PostTransform(payload)
// 		s.Nil(err)
// 		s.Equal(payload, response)
// 	})

// 	s.Run("should return error if slack token decryption failed", func() {
// 		configurations := make(map[string]interface{})
// 		configurations["token"] = "YmFy"
// 		payload := &receiver.Receiver{
// 			Configurations: configurations,
// 		}

// 		slackHelper, err := receiver.NewSlackHelper(s.exchangerMock, "abcdefghijklmnopqrstuvwxyzabcdef")
// 		s.NoError(err)

// 		var oldCryptopastaDecryptor = cryptopastaDecryptor
// 		defer func() {
// 			cryptopastaEncryptor = oldCryptopastaDecryptor
// 		}()
// 		cryptopastaDecryptor = func(_ []byte, _ *[32]byte) ([]byte, error) {
// 			return nil, errors.New("random error")
// 		}

// 		err = slackHelper.PostTransform(payload)
// 		s.EqualError(err, "slackHelper.Decrypt: random error")
// 	})
// }
