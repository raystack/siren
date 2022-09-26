package namespace_test

import (
	"context"
	testing "testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/odpf/siren/core/namespace"
	"github.com/odpf/siren/core/namespace/mocks"
	"github.com/odpf/siren/core/provider"
	"github.com/odpf/siren/pkg/errors"
	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
)

func TestService_ListNamespaces(t *testing.T) {

	type testCase struct {
		Description        string
		ExpectedNamespaces []namespace.Namespace
		Setup              func(*mocks.NamespaceRepository, *mocks.Encryptor, testCase)
		Err                error
	}
	var (
		ctx       = context.TODO()
		timeNow   = time.Now()
		testCases = []testCase{

			{
				Description: "should return error if List repository error",
				Setup: func(rr *mocks.NamespaceRepository, e *mocks.Encryptor, tc testCase) {
					rr.EXPECT().List(mock.AnythingOfType("*context.emptyCtx")).Return(nil, errors.New("some error"))
				},
				Err: errors.New("some error"),
			},
			{
				Description: "should return error if List repository success and decrypt error",
				Setup: func(rr *mocks.NamespaceRepository, e *mocks.Encryptor, tc testCase) {
					rr.EXPECT().List(mock.AnythingOfType("*context.emptyCtx")).Return([]namespace.EncryptedNamespace{
						{
							Namespace: &namespace.Namespace{
								ID: 1,
								Provider: provider.Provider{
									ID: 1,
								},
								Name:      "foo",
								Labels:    map[string]string{"foo": "bar"},
								CreatedAt: timeNow,
								UpdatedAt: timeNow,
							},
							Credentials: `encrypted-text-1`,
						},
						{
							Namespace: &namespace.Namespace{
								ID: 2,
								Provider: provider.Provider{
									ID: 1,
								},
								Name:      "foo",
								Labels:    map[string]string{"foo": "bar"},
								CreatedAt: timeNow,
								UpdatedAt: timeNow,
							},
							Credentials: `encrypted-text-2`,
						},
					}, nil)
					e.EXPECT().Decrypt(mock.AnythingOfType("string")).Return("", errors.New("decrypt error"))
				},
				Err: errors.New("decrypt error"),
			},
			{
				Description: "should return error if list repository success and decrypted object is not json",
				Setup: func(rr *mocks.NamespaceRepository, e *mocks.Encryptor, tc testCase) {
					rr.EXPECT().List(mock.AnythingOfType("*context.emptyCtx")).Return([]namespace.EncryptedNamespace{
						{
							Namespace: &namespace.Namespace{
								ID: 1,
								Provider: provider.Provider{
									ID: 1,
								},
								Name:      "foo",
								Labels:    map[string]string{"foo": "bar"},
								CreatedAt: timeNow,
								UpdatedAt: timeNow,
							},
							Credentials: `encrypted-text-1`,
						},
						{
							Namespace: &namespace.Namespace{
								ID: 2,
								Provider: provider.Provider{
									ID: 1,
								},
								Name:      "foo",
								Labels:    map[string]string{"foo": "bar"},
								CreatedAt: timeNow,
								UpdatedAt: timeNow,
							},
							Credentials: `encrypted-text-2`,
						},
					}, nil)
					e.EXPECT().Decrypt(mock.AnythingOfType("string")).Return("", nil)
				},
				Err: errors.New("unexpected end of JSON input"),
			},
			{
				Description: "should success if list repository and decrypt success",
				Setup: func(rr *mocks.NamespaceRepository, e *mocks.Encryptor, tc testCase) {
					rr.EXPECT().List(mock.AnythingOfType("*context.emptyCtx")).Return([]namespace.EncryptedNamespace{
						{
							Namespace: &namespace.Namespace{
								ID: 1,
								Provider: provider.Provider{
									ID: 1,
								},
								Name:      "foo",
								Labels:    map[string]string{"foo": "bar"},
								CreatedAt: timeNow,
								UpdatedAt: timeNow,
							},
							Credentials: `encrypted-text-1`,
						},
						{
							Namespace: &namespace.Namespace{
								ID: 2,
								Provider: provider.Provider{
									ID: 1,
								},
								Name:      "foo",
								Labels:    map[string]string{"foo": "bar"},
								CreatedAt: timeNow,
								UpdatedAt: timeNow,
							},
							Credentials: `encrypted-text-2`,
						},
					}, nil)
					e.EXPECT().Decrypt(mock.AnythingOfType("string")).Return("{\"name\": \"a\"}", nil)
				},
				ExpectedNamespaces: []namespace.Namespace{
					{
						ID: 1,
						Provider: provider.Provider{
							ID: 1,
						},
						Name:   "foo",
						Labels: map[string]string{"foo": "bar"},
						Credentials: map[string]interface{}{
							"name": "a",
						},
						CreatedAt: timeNow,
						UpdatedAt: timeNow,
					},
					{
						ID: 2,
						Provider: provider.Provider{
							ID: 1,
						},
						Name:   "foo",
						Labels: map[string]string{"foo": "bar"},
						Credentials: map[string]interface{}{
							"name": "a",
						},
						CreatedAt: timeNow,
						UpdatedAt: timeNow,
					},
				},
				Err: nil,
			},
		}
	)

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			var (
				repositoryMock = new(mocks.NamespaceRepository)
				encryptorMock  = new(mocks.Encryptor)
			)
			svc := namespace.NewService(encryptorMock, repositoryMock)

			tc.Setup(repositoryMock, encryptorMock, tc)

			got, err := svc.List(ctx)
			if tc.Err != err {
				if tc.Err.Error() != err.Error() {
					t.Fatalf("got error %s, expected was %s", err.Error(), tc.Err.Error())
				}
			}
			if !cmp.Equal(got, tc.ExpectedNamespaces) {
				t.Fatalf("got result %+v, expected was %+v", got, tc.ExpectedNamespaces)
			}
			repositoryMock.AssertExpectations(t)
			encryptorMock.AssertExpectations(t)
		})
	}
}

func TestService_CreateNamespace(t *testing.T) {

	type testCase struct {
		Description string
		NSpace      *namespace.Namespace
		Setup       func(*mocks.NamespaceRepository, *mocks.Encryptor, testCase)
		Err         error
	}
	var (
		ctx       = context.TODO()
		testCases = []testCase{
			{
				Description: "should return error if namespace is nil",
				Setup:       func(rr *mocks.NamespaceRepository, e *mocks.Encryptor, tc testCase) {},
				Err:         errors.New("incoming namespace is empty"),
			},
			{
				Description: "should return error if encrypt return error caused credential is not in json",
				Setup:       func(rr *mocks.NamespaceRepository, e *mocks.Encryptor, tc testCase) {},
				NSpace: &namespace.Namespace{
					Credentials: map[string]interface{}{
						"invalid": make(chan int),
					},
				},
				Err: errors.New("json: unsupported type: chan int"),
			},
			{
				Description: "should return error if encrypt return error",
				Setup: func(rr *mocks.NamespaceRepository, e *mocks.Encryptor, tc testCase) {
					e.EXPECT().Encrypt(mock.AnythingOfType("string")).Return("", errors.New("some error"))
				},
				NSpace: &namespace.Namespace{
					Credentials: map[string]interface{}{
						"credential": "value",
					},
				},
				Err: errors.New("some error"),
			},
			{
				Description: "should return error if encrypt success and create repository error",
				Setup: func(rr *mocks.NamespaceRepository, e *mocks.Encryptor, tc testCase) {
					e.EXPECT().Encrypt(mock.AnythingOfType("string")).Return("some-ciphertext", nil)
					rr.EXPECT().Create(mock.AnythingOfType("*context.emptyCtx"), &namespace.EncryptedNamespace{
						Namespace:   tc.NSpace,
						Credentials: "some-ciphertext",
					}).Return(errors.New("some error"))
				},
				NSpace: &namespace.Namespace{
					Credentials: map[string]interface{}{
						"credential": "value",
					},
				},
				Err: errors.New("some error"),
			},
			{
				Description: "should return error conflict if encrypt success and create repository return duplicate error",
				Setup: func(rr *mocks.NamespaceRepository, e *mocks.Encryptor, tc testCase) {
					e.EXPECT().Encrypt(mock.AnythingOfType("string")).Return("some-ciphertext", nil)
					rr.EXPECT().Create(mock.AnythingOfType("*context.emptyCtx"), &namespace.EncryptedNamespace{
						Namespace:   tc.NSpace,
						Credentials: "some-ciphertext",
					}).Return(namespace.ErrDuplicate)
				},
				NSpace: &namespace.Namespace{
					Credentials: map[string]interface{}{
						"credential": "value",
					},
				},
				Err: errors.New("urn and provider pair already exist"),
			},
			{
				Description: "should return error not found if encrypt success and create repository return relation error",
				Setup: func(rr *mocks.NamespaceRepository, e *mocks.Encryptor, tc testCase) {
					e.EXPECT().Encrypt(mock.AnythingOfType("string")).Return("some-ciphertext", nil)
					rr.EXPECT().Create(mock.AnythingOfType("*context.emptyCtx"), &namespace.EncryptedNamespace{
						Namespace:   tc.NSpace,
						Credentials: "some-ciphertext",
					}).Return(namespace.ErrRelation)
				},
				NSpace: &namespace.Namespace{
					Credentials: map[string]interface{}{
						"credential": "value",
					},
				},
				Err: errors.New("provider id does not exist"),
			},
			{
				Description: "should return nil error if encrypt success and create repository success",
				Setup: func(rr *mocks.NamespaceRepository, e *mocks.Encryptor, tc testCase) {
					e.EXPECT().Encrypt(mock.AnythingOfType("string")).Return("some-ciphertext", nil)
					rr.EXPECT().Create(mock.AnythingOfType("*context.emptyCtx"), &namespace.EncryptedNamespace{
						Namespace:   tc.NSpace,
						Credentials: "some-ciphertext",
					}).Return(nil)
				},
				NSpace: &namespace.Namespace{
					Credentials: map[string]interface{}{
						"credential": "value",
					},
				},
				Err: nil,
			},
		}
	)
	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			var (
				repositoryMock = new(mocks.NamespaceRepository)
				encryptorMock  = new(mocks.Encryptor)
			)
			svc := namespace.NewService(encryptorMock, repositoryMock)

			tc.Setup(repositoryMock, encryptorMock, tc)

			err := svc.Create(ctx, tc.NSpace)
			if tc.Err != err {
				if tc.Err.Error() != err.Error() {
					t.Fatalf("got error %s, expected was %s", err.Error(), tc.Err.Error())
				}
			}

			repositoryMock.AssertExpectations(t)
			encryptorMock.AssertExpectations(t)
		})
	}
}

func TestService_GetNamespace(t *testing.T) {
	type testCase struct {
		Description string
		NSpace      *namespace.Namespace
		Setup       func(*mocks.NamespaceRepository, *mocks.Encryptor, testCase)
		Err         error
	}
	var (
		ctx       = context.TODO()
		testID    = uint64(10)
		testCases = []testCase{
			{
				Description: "should return error if Get repository error",
				Setup: func(rr *mocks.NamespaceRepository, e *mocks.Encryptor, tc testCase) {
					rr.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), testID).Return(nil, errors.New("some error"))
				},
				Err: errors.New("some error"),
			},
			{
				Description: "should return error not found if Get repository return not found error",
				Setup: func(rr *mocks.NamespaceRepository, e *mocks.Encryptor, tc testCase) {
					rr.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), testID).Return(nil, namespace.NotFoundError{})
				},
				Err: errors.New("namespace not found"),
			},
			{
				Description: "should return error if Get repository success and decrypt return error",
				Setup: func(rr *mocks.NamespaceRepository, e *mocks.Encryptor, tc testCase) {
					rr.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), testID).Return(&namespace.EncryptedNamespace{
						Namespace:   tc.NSpace,
						Credentials: "some-ciphertext",
					}, nil)
					e.EXPECT().Decrypt("some-ciphertext").Return("", errors.New("some error"))
				},
				Err: errors.New("some error"),
			},
			{
				Description: "should return error if Get repository success and decrypted credentials is not json marshallable",
				Setup: func(rr *mocks.NamespaceRepository, e *mocks.Encryptor, tc testCase) {
					rr.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), testID).Return(&namespace.EncryptedNamespace{
						Namespace:   tc.NSpace,
						Credentials: "some-ciphertext",
					}, nil)
					e.EXPECT().Decrypt("some-ciphertext").Return("", nil)
				},
				Err: errors.New("unexpected end of JSON input"),
			},
			{
				Description: "should return nil error if Get repository success and decrypt success",
				Setup: func(rr *mocks.NamespaceRepository, e *mocks.Encryptor, tc testCase) {
					rr.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), testID).Return(&namespace.EncryptedNamespace{
						Namespace:   tc.NSpace,
						Credentials: "some-ciphertext",
					}, nil)
					e.EXPECT().Decrypt("some-ciphertext").Return("{ \"key\": \"value\" }", nil)
				},
				NSpace: &namespace.Namespace{
					Credentials: map[string]interface{}{
						"key": "value",
					},
				},
				Err: nil,
			},
		}
	)

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			var (
				repositoryMock = new(mocks.NamespaceRepository)
				encryptorMock  = new(mocks.Encryptor)
			)
			svc := namespace.NewService(encryptorMock, repositoryMock)

			tc.Setup(repositoryMock, encryptorMock, tc)

			got, err := svc.Get(ctx, testID)
			if tc.Err != err {
				if tc.Err.Error() != err.Error() {
					t.Fatalf("got error %s, expected was %s", err.Error(), tc.Err.Error())
				}
			}
			if !cmp.Equal(got, tc.NSpace) {
				t.Fatalf("got result %+v, expected was %+v", got, tc.NSpace)
			}
			repositoryMock.AssertExpectations(t)
			encryptorMock.AssertExpectations(t)
		})
	}
}

func TestService_UpdateNamespace(t *testing.T) {
	type testCase struct {
		Description string
		NSpace      *namespace.Namespace
		Setup       func(*mocks.NamespaceRepository, *mocks.Encryptor, testCase)
		Err         error
	}
	var (
		ctx       = context.TODO()
		testCases = []testCase{
			{
				Description: "should return error if encrypt return error caused credential is not in json",
				Setup:       func(rr *mocks.NamespaceRepository, e *mocks.Encryptor, tc testCase) {},
				NSpace: &namespace.Namespace{
					Credentials: map[string]interface{}{
						"invalid": make(chan int),
					},
				},
				Err: errors.New("json: unsupported type: chan int"),
			},
			{
				Description: "should return error if encrypt return error",
				Setup: func(rr *mocks.NamespaceRepository, e *mocks.Encryptor, tc testCase) {
					e.EXPECT().Encrypt(mock.AnythingOfType("string")).Return("", errors.New("some error"))
				},
				NSpace: &namespace.Namespace{
					Credentials: map[string]interface{}{
						"credential": "value",
					},
				},
				Err: errors.New("some error"),
			},
			{
				Description: "should return error if encrypt success and update repository error",
				Setup: func(rr *mocks.NamespaceRepository, e *mocks.Encryptor, tc testCase) {
					e.EXPECT().Encrypt(mock.AnythingOfType("string")).Return("some-ciphertext", nil)
					rr.EXPECT().Update(mock.AnythingOfType("*context.emptyCtx"), &namespace.EncryptedNamespace{
						Namespace:   tc.NSpace,
						Credentials: "some-ciphertext",
					}).Return(errors.New("some error"))
				},
				NSpace: &namespace.Namespace{
					Credentials: map[string]interface{}{
						"credential": "value",
					},
				},
				Err: errors.New("some error"),
			},
			{
				Description: "should return error not found if encrypt success and update repository return not found error",
				Setup: func(rr *mocks.NamespaceRepository, e *mocks.Encryptor, tc testCase) {
					e.EXPECT().Encrypt(mock.AnythingOfType("string")).Return("some-ciphertext", nil)
					rr.EXPECT().Update(mock.AnythingOfType("*context.emptyCtx"), &namespace.EncryptedNamespace{
						Namespace:   tc.NSpace,
						Credentials: "some-ciphertext",
					}).Return(namespace.NotFoundError{})
				},
				NSpace: &namespace.Namespace{
					Credentials: map[string]interface{}{
						"credential": "value",
					},
				},
				Err: errors.New("namespace not found"),
			},
			{
				Description: "should return error not found if encrypt success and update repository return relation error",
				Setup: func(rr *mocks.NamespaceRepository, e *mocks.Encryptor, tc testCase) {
					e.EXPECT().Encrypt(mock.AnythingOfType("string")).Return("some-ciphertext", nil)
					rr.EXPECT().Update(mock.AnythingOfType("*context.emptyCtx"), &namespace.EncryptedNamespace{
						Namespace:   tc.NSpace,
						Credentials: "some-ciphertext",
					}).Return(namespace.ErrRelation)
				},
				NSpace: &namespace.Namespace{
					Credentials: map[string]interface{}{
						"credential": "value",
					},
				},
				Err: errors.New("provider id does not exist"),
			},
			{
				Description: "should return error conflict if encrypt success and update repository return error duplicate",
				Setup: func(rr *mocks.NamespaceRepository, e *mocks.Encryptor, tc testCase) {
					e.EXPECT().Encrypt(mock.AnythingOfType("string")).Return("some-ciphertext", nil)
					rr.EXPECT().Update(mock.AnythingOfType("*context.emptyCtx"), &namespace.EncryptedNamespace{
						Namespace:   tc.NSpace,
						Credentials: "some-ciphertext",
					}).Return(namespace.ErrDuplicate)
				},
				NSpace: &namespace.Namespace{
					Credentials: map[string]interface{}{
						"credential": "value",
					},
				},
				Err: errors.New("urn and provider pair already exist"),
			},
			{
				Description: "should return nil error if encrypt success and update repository success",
				Setup: func(rr *mocks.NamespaceRepository, e *mocks.Encryptor, tc testCase) {
					e.EXPECT().Encrypt(mock.AnythingOfType("string")).Return("some-ciphertext", nil)
					rr.EXPECT().Update(mock.AnythingOfType("*context.emptyCtx"), &namespace.EncryptedNamespace{
						Namespace:   tc.NSpace,
						Credentials: "some-ciphertext",
					}).Return(nil)
				},
				NSpace: &namespace.Namespace{
					Credentials: map[string]interface{}{
						"credential": "value",
					},
				},
				Err: nil,
			},
		}
	)

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			var (
				repositoryMock = new(mocks.NamespaceRepository)
				encryptorMock  = new(mocks.Encryptor)
			)
			svc := namespace.NewService(encryptorMock, repositoryMock)

			tc.Setup(repositoryMock, encryptorMock, tc)

			err := svc.Update(ctx, tc.NSpace)
			if tc.Err != err {
				if tc.Err.Error() != err.Error() {
					t.Fatalf("got error %s, expected was %s", err.Error(), tc.Err.Error())
				}
			}

			repositoryMock.AssertExpectations(t)
			encryptorMock.AssertExpectations(t)
		})
	}
}

func TestDeleteNamespace(t *testing.T) {
	ctx := context.TODO()
	namespaceID := uint64(10)

	t.Run("should call repository Delete method and return nil if no error", func(t *testing.T) {
		repositoryMock := &mocks.NamespaceRepository{}
		dummyService := namespace.NewService(nil, repositoryMock)
		repositoryMock.EXPECT().Delete(mock.AnythingOfType("*context.emptyCtx"), namespaceID).Return(nil).Once()
		err := dummyService.Delete(ctx, namespaceID)
		assert.Nil(t, err)
		repositoryMock.AssertExpectations(t)
	})

	t.Run("should call repository Delete method and return error if any", func(t *testing.T) {
		repositoryMock := &mocks.NamespaceRepository{}
		dummyService := namespace.NewService(nil, repositoryMock)
		repositoryMock.EXPECT().Delete(mock.AnythingOfType("*context.emptyCtx"), namespaceID).Return(errors.New("random error")).Once()
		err := dummyService.Delete(ctx, namespaceID)
		assert.EqualError(t, err, "random error")
		repositoryMock.AssertExpectations(t)
	})
}
