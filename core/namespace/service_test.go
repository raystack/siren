package namespace_test

import (
	"strings"
	"testing"
	"time"

	"github.com/odpf/siren/core/namespace"
	"github.com/odpf/siren/core/namespace/mocks"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_ListNamespaces(t *testing.T) {
	labels := map[string]string{"foo": "bar"}

	t.Run("should call repository List method and return result in decrypted format", func(t *testing.T) {
		repositoryMock := &mocks.NamespaceRepository{}
		transformerMock := &mocks.EncryptorDecryptor{}
		dummyService, err := namespace.NewService(repositoryMock, transformerMock)
		require.NoError(t, err)

		dummyNamespaces := []*namespace.EncryptedNamespace{
			{
				Namespace: &namespace.Namespace{
					Id:        1,
					Provider:  1,
					Name:      "foo",
					Labels:    labels,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Credentials: `encrypted-text-1`,
			},
			{
				Namespace: &namespace.Namespace{
					Id:        2,
					Provider:  1,
					Name:      "foo",
					Labels:    labels,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Credentials: `encrypted-text-2`,
			},
		}
		repositoryMock.On("List").Return(dummyNamespaces, nil).Once()
		transformerMock.On("Decrypt", "encrypted-text-1").
			Return(`{"foo":"bar"}`, nil).Once()
		transformerMock.On("Decrypt", "encrypted-text-2").
			Return(`{"bar":"baz"}`, nil).Once()
		result, err := dummyService.ListNamespaces()
		assert.Nil(t, err)
		assert.Equal(t, len(dummyNamespaces), len(result))
		assert.Equal(t, `bar`, result[0].Credentials["foo"])
		assert.Equal(t, `baz`, result[1].Credentials["bar"])
		repositoryMock.AssertCalled(t, "List")
		transformerMock.AssertExpectations(t)
	})

	t.Run("should call repository List method and return error if any", func(t *testing.T) {
		repositoryMock := &mocks.NamespaceRepository{}
		transformerMock := &mocks.EncryptorDecryptor{}
		dummyService, err := namespace.NewService(repositoryMock, transformerMock)
		require.NoError(t, err)

		repositoryMock.On("List").Return(nil, errors.New("random error")).Once()
		result, err := dummyService.ListNamespaces()
		assert.Nil(t, result)
		assert.EqualError(t, err, "s.repository.List: random error")
		repositoryMock.AssertCalled(t, "List")
		transformerMock.AssertExpectations(t)
	})

	t.Run("should decrypt the repository response and return error if any", func(t *testing.T) {
		repositoryMock := &mocks.NamespaceRepository{}
		transformerMock := &mocks.EncryptorDecryptor{}
		dummyService, err := namespace.NewService(repositoryMock, transformerMock)
		require.NoError(t, err)

		dummyNamespaces := []*namespace.EncryptedNamespace{
			{
				Namespace: &namespace.Namespace{
					Id:        1,
					Provider:  1,
					Name:      "foo",
					Labels:    labels,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Credentials: `encrypted-text-1`,
			},
		}
		repositoryMock.On("List").Return(dummyNamespaces, nil).Once()
		transformerMock.On("Decrypt", "encrypted-text-1").
			Return(`{"bar":"baz"}`, errors.New("random error")).Once()
		result, err := dummyService.ListNamespaces()
		assert.Nil(t, result)
		assert.EqualError(t, err, "s.transformer.Decrypt: random error")
		repositoryMock.AssertCalled(t, "List")
		transformerMock.AssertExpectations(t)
	})

	t.Run("should unmarshal decrypted response and return error if any", func(t *testing.T) {
		repositoryMock := &mocks.NamespaceRepository{}
		transformerMock := &mocks.EncryptorDecryptor{}
		dummyService, err := namespace.NewService(repositoryMock, transformerMock)
		require.NoError(t, err)

		dummyNamespaces := []*namespace.EncryptedNamespace{
			{
				Namespace: &namespace.Namespace{
					Id:        1,
					Provider:  1,
					Name:      "foo",
					Labels:    labels,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Credentials: `encrypted-text-1`,
			},
		}
		repositoryMock.On("List").Return(dummyNamespaces, nil).Once()
		transformerMock.On("Decrypt", "encrypted-text-1").
			Return(`abcd`, nil).Once()
		result, err := dummyService.ListNamespaces()
		assert.Nil(t, result)
		assert.True(t, strings.Contains(err.Error(), `json.Unmarshal: invalid character`))
		repositoryMock.AssertCalled(t, "List")
		transformerMock.AssertExpectations(t)
	})
}

func TestService_CreateNamespaces(t *testing.T) {
	credentials := make(map[string]interface{})
	credentials["foo"] = "bar"
	labels := map[string]string{"foo": "bar"}
	timeNow := time.Now()
	input := &namespace.Namespace{
		Provider:    1,
		Name:        "foo",
		Credentials: credentials,
		Labels:      labels,
		CreatedAt:   timeNow,
		UpdatedAt:   timeNow,
	}

	t.Run("should call repository Create method and return result in domain's type", func(t *testing.T) {
		repositoryMock := &mocks.NamespaceRepository{}
		transformerMock := &mocks.EncryptorDecryptor{}
		expectedID := uint64(2)
		dummyService, err := namespace.NewService(repositoryMock, transformerMock)
		require.NoError(t, err)

		repositoryMock.On("Create", mock.AnythingOfType("*namespace.EncryptedNamespace")).
			Run(func(args mock.Arguments) {
				rarg := args.Get(0)
				r := rarg.(*namespace.EncryptedNamespace)
				assert.Equal(t, uint64(0), r.Namespace.Id)
				r.Namespace.Id = expectedID // simulate ID creation from repository
				assert.Equal(t, "foo", r.Name)
				assert.Equal(t, uint64(1), r.Provider)
			}).Return(nil).Once()
		transformerMock.On("Encrypt", `{"foo":"bar"}`).
			Return("encrypted-text-1", nil).Once()
		err = dummyService.CreateNamespace(input)
		assert.Nil(t, err)
		assert.Equal(t, expectedID, input.Id)
		repositoryMock.AssertExpectations(t)
		transformerMock.AssertExpectations(t)
	})

	t.Run("should encrypt credentials before passing to repository and return decrypted creds", func(t *testing.T) {
		repositoryMock := &mocks.NamespaceRepository{}
		dummyTransformer, _ := namespace.NewTransformer("abcdefghijklmnopqrstuvwxyzabcdef")
		dummyService, err := namespace.NewService(repositoryMock, dummyTransformer)
		require.NoError(t, err)

		repositoryMock.On("Create", mock.AnythingOfType("*namespace.EncryptedNamespace")).
			Run(func(args mock.Arguments) {
				param := args.Get(0).(*namespace.EncryptedNamespace)
				assert.NotEmpty(t, param.Credentials)
				assert.NotEqual(t, input.Credentials, param.Credentials)
			}).Return(nil).Once()

		err = dummyService.CreateNamespace(input)
		assert.Nil(t, err)
		assert.Equal(t, credentials, input.Credentials)
		repositoryMock.AssertExpectations(t)
	})

	t.Run("should call repository Create method and return error if any", func(t *testing.T) {
		repositoryMock := &mocks.NamespaceRepository{}
		transformerMock := &mocks.EncryptorDecryptor{}
		dummyService, err := namespace.NewService(repositoryMock, transformerMock)
		require.NoError(t, err)

		repositoryMock.On("Create", mock.AnythingOfType("*namespace.EncryptedNamespace")).
			Return(errors.New("random error")).Once()
		transformerMock.On("Encrypt", `{"foo":"bar"}`).
			Return("encrypted-text-1", nil).Once()
		err = dummyService.CreateNamespace(input)
		assert.EqualError(t, err, "s.repository.Create: random error")
		repositoryMock.AssertExpectations(t)
	})

	t.Run("should encrypt credentials and return error if any", func(t *testing.T) {
		repositoryMock := &mocks.NamespaceRepository{}
		transformerMock := &mocks.EncryptorDecryptor{}
		dummyService, err := namespace.NewService(repositoryMock, transformerMock)
		require.NoError(t, err)

		transformerMock.On("Encrypt", `{"foo":"bar"}`).
			Return("encrypted-text-1", errors.New("random error")).Once()
		err = dummyService.CreateNamespace(input)
		assert.EqualError(t, err, "s.transformer.Encrypt: random error")
		transformerMock.AssertExpectations(t)
	})

	t.Run("should marshal credentials and return error if any", func(t *testing.T) {
		repositoryMock := &mocks.NamespaceRepository{}
		transformerMock := &mocks.EncryptorDecryptor{}
		dummyService, err := namespace.NewService(repositoryMock, transformerMock)
		require.NoError(t, err)

		badNamespace := &namespace.Namespace{
			Id:       2,
			Provider: 1,
			Name:     "foo",
			Credentials: map[string]interface{}{
				"foo": make(chan int),
			},
			Labels:    labels,
			CreatedAt: timeNow,
			UpdatedAt: timeNow,
		}
		err = dummyService.CreateNamespace(badNamespace)
		assert.True(t, strings.Contains(err.Error(), `json.Marshal: json: unsupported type: chan int`))
	})
}

func TestService_GetNamespace(t *testing.T) {
	labels := map[string]string{"foo": "bar"}

	t.Run("should call repository Get method and return result in decrypted format", func(t *testing.T) {
		repositoryMock := &mocks.NamespaceRepository{}
		transformerMock := &mocks.EncryptorDecryptor{}
		dummyService, err := namespace.NewService(repositoryMock, transformerMock)
		require.NoError(t, err)

		dummyNamespace := &namespace.EncryptedNamespace{
			Namespace: &namespace.Namespace{
				Id:        1,
				Provider:  1,
				Name:      "foo",
				Labels:    labels,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Credentials: `encrypted-text-1`,
		}
		repositoryMock.On("Get", uint64(1)).Return(dummyNamespace, nil).Once()
		transformerMock.On("Decrypt", "encrypted-text-1").
			Return(`{"foo":"bar"}`, nil).Once()
		result, err := dummyService.GetNamespace(uint64(1))
		assert.Nil(t, err)
		assert.Equal(t, `bar`, result.Credentials["foo"])
		repositoryMock.AssertCalled(t, "Get", uint64(1))
		transformerMock.AssertExpectations(t)
	})

	t.Run("should return decrypted namespace credentials", func(t *testing.T) {
		repositoryMock := &mocks.NamespaceRepository{}
		dummyTransformer, _ := namespace.NewTransformer("abcdefghijklmnopqrstuvwxyzabcdef")
		dummyService, err := namespace.NewService(repositoryMock, dummyTransformer)
		require.NoError(t, err)

		expectedCredentialsStr := `{"foo":"bar"}`
		expectedCredentials := map[string]interface{}{"foo": "bar"}
		expectedEncryptedCredentials, _ := dummyTransformer.Encrypt(expectedCredentialsStr)
		expectedEncruptedNamespace := &namespace.EncryptedNamespace{
			Namespace:   &namespace.Namespace{},
			Credentials: expectedEncryptedCredentials,
		}
		repositoryMock.On("Get", mock.AnythingOfType("uint64")).Return(expectedEncruptedNamespace, nil).Once()

		result, err := dummyService.GetNamespace(uint64(1))
		assert.Nil(t, err)
		assert.Equal(t, expectedCredentials, result.Credentials)
	})

	t.Run("should call repository Get method and return nil if namespace does not exist", func(t *testing.T) {
		repositoryMock := &mocks.NamespaceRepository{}
		transformerMock := &mocks.EncryptorDecryptor{}
		dummyService, err := namespace.NewService(repositoryMock, transformerMock)
		require.NoError(t, err)

		repositoryMock.On("Get", uint64(1)).Return(nil, nil).Once()
		result, err := dummyService.GetNamespace(uint64(1))
		assert.Nil(t, err)
		assert.Nil(t, result)
	})

	t.Run("should call repository Get method and handle error if any", func(t *testing.T) {
		repositoryMock := &mocks.NamespaceRepository{}
		transformerMock := &mocks.EncryptorDecryptor{}
		dummyService, err := namespace.NewService(repositoryMock, transformerMock)
		require.NoError(t, err)

		repositoryMock.On("Get", uint64(1)).
			Return(nil, errors.New("random error")).Once()
		result, err := dummyService.GetNamespace(uint64(1))
		assert.Nil(t, result)
		assert.EqualError(t, err, "s.repository.Get: random error")
		repositoryMock.AssertCalled(t, "Get", uint64(1))
	})

	t.Run("should decrypt credentials and return error if any", func(t *testing.T) {
		repositoryMock := &mocks.NamespaceRepository{}
		transformerMock := &mocks.EncryptorDecryptor{}
		dummyService, err := namespace.NewService(repositoryMock, transformerMock)
		require.NoError(t, err)

		dummyNamespace := &namespace.EncryptedNamespace{
			Namespace: &namespace.Namespace{
				Id:        1,
				Provider:  1,
				Name:      "foo",
				Labels:    labels,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Credentials: `encrypted-text-1`,
		}
		repositoryMock.On("Get", uint64(1)).Return(dummyNamespace, nil).Once()
		transformerMock.On("Decrypt", "encrypted-text-1").
			Return(`{"foo":"bar"}`, errors.New("random error")).Once()
		result, err := dummyService.GetNamespace(uint64(1))
		assert.Nil(t, result)
		assert.EqualError(t, err, "s.transformer.Decrypt: random error")
		repositoryMock.AssertCalled(t, "Get", uint64(1))
		transformerMock.AssertCalled(t, "Decrypt", "encrypted-text-1")
	})
}

func TestService_UpdateNamespaces(t *testing.T) {
	credentials := make(map[string]interface{})
	credentials["foo"] = "bar"
	labels := map[string]string{"foo": "bar"}
	timeNow := time.Now()
	dummyNamespace := &namespace.Namespace{
		Id:          2,
		Provider:    1,
		Name:        "foo",
		Credentials: credentials,
		Labels:      labels,
		CreatedAt:   timeNow,
		UpdatedAt:   timeNow,
	}

	t.Run("should call repository Update method and return result in domain's type", func(t *testing.T) {
		repositoryMock := &mocks.NamespaceRepository{}
		transformerMock := &mocks.EncryptorDecryptor{}
		dummyService, err := namespace.NewService(repositoryMock, transformerMock)
		require.NoError(t, err)

		repositoryMock.On("Update", mock.AnythingOfType("*namespace.EncryptedNamespace")).
			Run(func(args mock.Arguments) {
				rarg := args.Get(0)
				r := rarg.(*namespace.EncryptedNamespace)
				assert.Equal(t, "foo", r.Name)
				assert.Equal(t, uint64(2), r.Id)
				assert.Equal(t, uint64(1), r.Provider)
			}).Return(nil).Once()
		transformerMock.On("Encrypt", `{"foo":"bar"}`).
			Return("encrypted-text-1", nil).Once()
		err = dummyService.UpdateNamespace(dummyNamespace)
		assert.Nil(t, err)
		repositoryMock.AssertExpectations(t)
		transformerMock.AssertExpectations(t)
	})

	t.Run("should call repository Update method and return error if any", func(t *testing.T) {
		repositoryMock := &mocks.NamespaceRepository{}
		transformerMock := &mocks.EncryptorDecryptor{}
		dummyService, err := namespace.NewService(repositoryMock, transformerMock)
		require.NoError(t, err)

		repositoryMock.On("Update", mock.AnythingOfType("*namespace.EncryptedNamespace")).
			Return(errors.New("random error")).Once()
		transformerMock.On("Encrypt", `{"foo":"bar"}`).
			Return("encrypted-text-1", nil).Once()
		err = dummyService.UpdateNamespace(dummyNamespace)
		assert.EqualError(t, err, "s.repository.Update: random error")
		repositoryMock.AssertExpectations(t)
	})

	t.Run("should encrypt credentials and return error if any", func(t *testing.T) {
		repositoryMock := &mocks.NamespaceRepository{}
		transformerMock := &mocks.EncryptorDecryptor{}
		dummyService, err := namespace.NewService(repositoryMock, transformerMock)
		require.NoError(t, err)

		transformerMock.On("Encrypt", `{"foo":"bar"}`).
			Return("encrypted-text-1", errors.New("random error")).Once()
		err = dummyService.UpdateNamespace(dummyNamespace)
		assert.EqualError(t, err, "s.transformer.Encrypt: random error")
		transformerMock.AssertExpectations(t)
	})

	t.Run("should marshal credentials and return error if any", func(t *testing.T) {
		repositoryMock := &mocks.NamespaceRepository{}
		transformerMock := &mocks.EncryptorDecryptor{}
		dummyService, err := namespace.NewService(repositoryMock, transformerMock)
		require.NoError(t, err)

		badNamespace := &namespace.Namespace{
			Id:       2,
			Provider: 1,
			Name:     "foo",
			Credentials: map[string]interface{}{
				"foo": make(chan int),
			},
			Labels:    labels,
			CreatedAt: timeNow,
			UpdatedAt: timeNow,
		}
		err = dummyService.UpdateNamespace(badNamespace)
		assert.True(t, strings.Contains(err.Error(), `json.Marshal: json: unsupported type: chan int`))
	})
}

func TestService_DeleteNamespace(t *testing.T) {
	credentials := make(map[string]interface{})
	credentials["foo"] = "bar"

	t.Run("should call repository Delete method and return nil if no error", func(t *testing.T) {
		repositoryMock := &mocks.NamespaceRepository{}
		dummyService, err := namespace.NewService(repositoryMock, nil)
		require.NoError(t, err)

		repositoryMock.On("Delete", uint64(1)).Return(nil).Once()
		err = dummyService.DeleteNamespace(1)
		assert.Nil(t, err)
		repositoryMock.AssertCalled(t, "Delete", uint64(1))
	})

	t.Run("should call repository Delete method and return error if any", func(t *testing.T) {
		repositoryMock := &mocks.NamespaceRepository{}
		dummyService, err := namespace.NewService(repositoryMock, nil)
		require.NoError(t, err)

		repositoryMock.On("Delete", uint64(1)).Return(errors.New("random error")).Once()
		err = dummyService.DeleteNamespace(1)
		assert.EqualError(t, err, "random error")
		repositoryMock.AssertCalled(t, "Delete", uint64(1))
	})
}

func TestService_Migrate(t *testing.T) {
	t.Run("should call repository Migrate method and return result", func(t *testing.T) {
		repositoryMock := &mocks.NamespaceRepository{}
		dummyService, err := namespace.NewService(repositoryMock, nil)
		require.NoError(t, err)

		repositoryMock.On("Migrate").Return(nil).Once()
		err = dummyService.Migrate()
		assert.Nil(t, err)
		repositoryMock.AssertCalled(t, "Migrate")
	})
}

func TestTransformer_Encrypt(t *testing.T) {
	t.Run("should ", func(t *testing.T) {

	})
}
