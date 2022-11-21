package receiver_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/odpf/siren/core/receiver"
	"github.com/odpf/siren/core/receiver/mocks"
	"github.com/odpf/siren/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_ListReceivers(t *testing.T) {
	type testCase struct {
		Description string
		Receivers   []receiver.Receiver
		Setup       func(*mocks.ReceiverRepository, *mocks.ConfigResolver)
		Err         error
	}

	var (
		ctx       = context.TODO()
		timeNow   = time.Now()
		testCases = []testCase{
			{
				Description: "should return error if List repository error",
				Setup: func(rr *mocks.ReceiverRepository, ss *mocks.ConfigResolver) {
					rr.EXPECT().List(mock.AnythingOfType("*context.emptyCtx"), receiver.Filter{}).Return(nil, errors.New("some error"))
				},
				Err: errors.New("some error"),
			},
			{
				Description: "should return error if List repository success and decrypt error",
				Setup: func(rr *mocks.ReceiverRepository, ss *mocks.ConfigResolver) {
					rr.EXPECT().List(mock.AnythingOfType("*context.emptyCtx"), receiver.Filter{}).Return([]receiver.Receiver{
						{
							ID:   10,
							Name: "foo",
							Type: receiver.TypeSlack,
							Labels: map[string]string{
								"foo": "bar",
							},
							Configurations: map[string]interface{}{
								"token": "key",
							},
							CreatedAt: timeNow,
							UpdatedAt: timeNow,
						},
					}, nil)
					ss.EXPECT().PostHookDBTransformConfigs(mock.AnythingOfType("*context.emptyCtx"), map[string]interface{}{
						"token": "key",
					}).Return(nil, errors.New("decrypt error"))
				},
				Err: errors.New("decrypt error"),
			},
			{
				Description: "should return error if type unknown",
				Setup: func(rr *mocks.ReceiverRepository, ss *mocks.ConfigResolver) {
					rr.EXPECT().List(mock.AnythingOfType("*context.emptyCtx"), receiver.Filter{}).Return([]receiver.Receiver{
						{
							Type: "random",
						},
					}, nil)
				},
				Err: errors.New("unsupported receiver type: \"random\""),
			},
			{
				Description: "should success if list repository and decrypt success",
				Setup: func(rr *mocks.ReceiverRepository, ss *mocks.ConfigResolver) {
					rr.EXPECT().List(mock.AnythingOfType("*context.emptyCtx"), receiver.Filter{}).Return([]receiver.Receiver{
						{
							ID:   10,
							Name: "foo",
							Type: receiver.TypeSlack,
							Labels: map[string]string{
								"foo": "bar",
							},
							Configurations: map[string]interface{}{
								"token": "key",
							},
							CreatedAt: timeNow,
							UpdatedAt: timeNow,
						},
					}, nil)
					ss.EXPECT().PostHookDBTransformConfigs(mock.AnythingOfType("*context.emptyCtx"), map[string]interface{}{
						"token": "key",
					}).Return(map[string]interface{}{
						"token": "decrypted_key",
					}, nil)
				},
				Receivers: []receiver.Receiver{
					{
						ID:   10,
						Name: "foo",
						Type: receiver.TypeSlack,
						Labels: map[string]string{
							"foo": "bar",
						},
						Configurations: map[string]interface{}{
							"token": "decrypted_key",
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
				repositoryMock = new(mocks.ReceiverRepository)
				resolverMock   = new(mocks.ConfigResolver)
			)

			registry := map[string]receiver.ConfigResolver{
				receiver.TypeSlack: resolverMock,
			}

			svc := receiver.NewService(repositoryMock, registry)

			tc.Setup(repositoryMock, resolverMock)

			got, err := svc.List(ctx, receiver.Filter{})
			if tc.Err != err {
				if tc.Err.Error() != err.Error() {
					t.Fatalf("got error %s, expected was %s", err.Error(), tc.Err.Error())
				}
			}
			if !cmp.Equal(got, tc.Receivers) {
				t.Fatalf("got result %+v, expected was %+v", got, tc.Receivers)
			}
			repositoryMock.AssertExpectations(t)
			resolverMock.AssertExpectations(t)
		})
	}
}

func TestService_CreateReceiver(t *testing.T) {
	type testCase struct {
		Description string
		Setup       func(*mocks.ReceiverRepository, *mocks.ConfigResolver)
		Rcv         *receiver.Receiver
		Err         error
	}
	var (
		ctx       = context.TODO()
		testCases = []testCase{
			{
				Description: "should return error if type unknown",
				Setup:       func(rr *mocks.ReceiverRepository, ss *mocks.ConfigResolver) {},
				Rcv: &receiver.Receiver{
					Type: "random",
				},
				Err: errors.New("unsupported receiver type: \"random\""),
			},
			{
				Description: "should return error if PreHookDBTransformConfigs return error",
				Setup: func(rr *mocks.ReceiverRepository, ss *mocks.ConfigResolver) {
					ss.EXPECT().PreHookDBTransformConfigs(mock.AnythingOfType("*context.emptyCtx"), map[string]interface{}{"token": "key"}).Return(nil, errors.New("some error"))

				},
				Rcv: &receiver.Receiver{
					ID:   123,
					Type: receiver.TypeSlack,
					Configurations: map[string]interface{}{
						"token": "key",
					},
				},
				Err: errors.New("some error"),
			},
			{
				Description: "should return error if Create repository return error",
				Setup: func(rr *mocks.ReceiverRepository, ss *mocks.ConfigResolver) {
					ss.EXPECT().PreHookDBTransformConfigs(mock.AnythingOfType("*context.emptyCtx"), map[string]interface{}{"token": "key"}).Return(map[string]interface{}{
						"token": "encrypted_key",
					}, nil)
					rr.EXPECT().Create(mock.AnythingOfType("*context.emptyCtx"), &receiver.Receiver{
						ID:   123,
						Type: receiver.TypeSlack,
						Configurations: map[string]interface{}{
							"token": "encrypted_key",
						},
					}).Return(errors.New("some error"))
				},
				Rcv: &receiver.Receiver{
					ID:   123,
					Type: receiver.TypeSlack,
					Configurations: map[string]interface{}{
						"token": "key",
					},
				},
				Err: errors.New("some error"),
			},
			{
				Description: "should return nil error if no error returned",
				Setup: func(rr *mocks.ReceiverRepository, ss *mocks.ConfigResolver) {
					ss.EXPECT().PreHookDBTransformConfigs(mock.AnythingOfType("*context.emptyCtx"), map[string]interface{}{"token": "key"}).Return(map[string]interface{}{
						"token": "encrypted_key",
					}, nil)
					rr.EXPECT().Create(mock.AnythingOfType("*context.emptyCtx"), &receiver.Receiver{
						ID:   123,
						Type: receiver.TypeSlack,
						Configurations: map[string]interface{}{
							"token": "encrypted_key",
						},
					}).Return(nil)
				},
				Rcv: &receiver.Receiver{
					ID:   123,
					Type: receiver.TypeSlack,
					Configurations: map[string]interface{}{
						"token": "key",
					},
				},
				Err: nil,
			},
		}
	)

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			var (
				repositoryMock = new(mocks.ReceiverRepository)
				resolverMock   = new(mocks.ConfigResolver)
			)

			registry := map[string]receiver.ConfigResolver{
				receiver.TypeSlack: resolverMock,
			}

			svc := receiver.NewService(repositoryMock, registry)

			tc.Setup(repositoryMock, resolverMock)

			err := svc.Create(ctx, tc.Rcv)
			if tc.Err != err {
				if tc.Err.Error() != err.Error() {
					t.Fatalf("got error %s, expected was %s", err.Error(), tc.Err.Error())
				}
			}
			repositoryMock.AssertExpectations(t)
			resolverMock.AssertExpectations(t)
		})
	}
}

func TestService_GetReceiver(t *testing.T) {
	type testCase struct {
		Description      string
		ExpectedReceiver *receiver.Receiver
		Setup            func(*mocks.ReceiverRepository, *mocks.ConfigResolver)
		Err              error
	}

	var (
		ctx       = context.TODO()
		timeNow   = time.Now()
		testID    = uint64(10)
		testCases = []testCase{
			{
				Description: "should return error if Get repository error",
				Setup: func(rr *mocks.ReceiverRepository, ss *mocks.ConfigResolver) {
					rr.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), testID).Return(nil, errors.New("some error"))
				},
				Err: errors.New("some error"),
			},
			{
				Description: "should return error if type unknown",
				Setup: func(rr *mocks.ReceiverRepository, ss *mocks.ConfigResolver) {
					rr.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), testID).Return(&receiver.Receiver{
						Type: "random",
					}, nil)
				},
				Err: errors.New("unsupported receiver type: \"random\""),
			},
			{
				Description: "should return error not found if Get repository return not found error",
				Setup: func(rr *mocks.ReceiverRepository, ss *mocks.ConfigResolver) {
					rr.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), testID).Return(nil, receiver.NotFoundError{})
				},
				Err: errors.New("receiver not found"),
			},
			{
				Description: "should return error if Get repository success and decrypt error",
				Setup: func(rr *mocks.ReceiverRepository, ss *mocks.ConfigResolver) {
					rr.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), testID).Return(&receiver.Receiver{
						ID:   10,
						Name: "foo",
						Type: receiver.TypeSlack,
						Labels: map[string]string{
							"foo": "bar",
						},
						Configurations: map[string]interface{}{
							"token": "key",
						},
						CreatedAt: timeNow,
						UpdatedAt: timeNow,
					}, nil)
					ss.EXPECT().PostHookDBTransformConfigs(mock.AnythingOfType("*context.emptyCtx"), map[string]interface{}{
						"token": "key",
					}).Return(nil, errors.New("decrypt error"))
				},
				Err: errors.New("decrypt error"),
			},
			{
				Description: "should success if Get repository and decrypt success",
				Setup: func(rr *mocks.ReceiverRepository, ss *mocks.ConfigResolver) {
					rr.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), testID).Return(&receiver.Receiver{
						ID:   10,
						Name: "foo",
						Type: receiver.TypeSlack,
						Labels: map[string]string{
							"foo": "bar",
						},
						Configurations: map[string]interface{}{
							"token": "key",
						},
						CreatedAt: timeNow,
						UpdatedAt: timeNow,
					}, nil)
					ss.EXPECT().PostHookDBTransformConfigs(mock.AnythingOfType("*context.emptyCtx"), map[string]interface{}{
						"token": "key",
					}).Return(map[string]interface{}{
						"token": "decrypted_key",
					}, nil)
					ss.EXPECT().BuildData(mock.AnythingOfType("*context.emptyCtx"), map[string]interface{}{
						"token": "decrypted_key",
					}).Return(map[string]interface{}{
						"newdata": "populated",
					}, nil)
				},
				ExpectedReceiver: &receiver.Receiver{
					ID:   10,
					Name: "foo",
					Type: receiver.TypeSlack,
					Labels: map[string]string{
						"foo": "bar",
					},
					Data: map[string]interface{}{
						"newdata": "populated",
					},
					Configurations: map[string]interface{}{
						"token": "decrypted_key",
					},
					CreatedAt: timeNow,
					UpdatedAt: timeNow,
				},
				Err: nil,
			},
		}
	)

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			var (
				repositoryMock = new(mocks.ReceiverRepository)
				resolverMock   = new(mocks.ConfigResolver)
			)

			registry := map[string]receiver.ConfigResolver{
				receiver.TypeSlack: resolverMock,
			}

			svc := receiver.NewService(repositoryMock, registry)

			tc.Setup(repositoryMock, resolverMock)

			got, err := svc.Get(ctx, testID)
			if tc.Err != err {
				if tc.Err.Error() != err.Error() {
					t.Fatalf("got error %s, expected was %s", err.Error(), tc.Err.Error())
				}
			}
			if !cmp.Equal(got, tc.ExpectedReceiver) {
				t.Fatalf("got result %+v, expected was %+v", got, tc.ExpectedReceiver)
			}
			repositoryMock.AssertExpectations(t)
			resolverMock.AssertExpectations(t)
		})
	}
}

func TestService_UpdateReceiver(t *testing.T) {
	type testCase struct {
		Description string
		Setup       func(*mocks.ReceiverRepository, *mocks.ConfigResolver)
		Rcv         *receiver.Receiver
		Err         error
	}
	var (
		ctx       = context.TODO()
		testCases = []testCase{
			{
				Description: "should return error if get receiver return error",
				Setup: func(rr *mocks.ReceiverRepository, ss *mocks.ConfigResolver) {
					rr.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(nil, errors.New("some error"))
				},
				Rcv: &receiver.Receiver{
					ID: 123,
				},
				Err: errors.New("some error"),
			},
			{
				Description: "should return error if type unknown",
				Setup: func(rr *mocks.ReceiverRepository, ss *mocks.ConfigResolver) {
					rr.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&receiver.Receiver{
						Type: "random",
					}, nil)
				},
				Rcv: &receiver.Receiver{
					ID: 123,
				},
				Err: errors.New("unsupported receiver type: \"random\""),
			},
			{
				Description: "should return error if PreHookDBTransformConfigs return error",
				Setup: func(rr *mocks.ReceiverRepository, ss *mocks.ConfigResolver) {
					rr.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&receiver.Receiver{
						ID:   123,
						Type: receiver.TypeSlack,
						Configurations: map[string]interface{}{
							"token": "key",
						},
					}, nil)
					ss.EXPECT().PreHookDBTransformConfigs(mock.AnythingOfType("*context.emptyCtx"), map[string]interface{}{"token": "key"}).Return(nil, errors.New("some error"))
				},
				Rcv: &receiver.Receiver{
					ID:   123,
					Type: receiver.TypeSlack,
					Configurations: map[string]interface{}{
						"token": "key",
					},
				},
				Err: errors.New("some error"),
			},
			{
				Description: "should return error if Update repository return error",
				Setup: func(rr *mocks.ReceiverRepository, ss *mocks.ConfigResolver) {
					rr.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&receiver.Receiver{
						ID:   123,
						Type: receiver.TypeSlack,
						Configurations: map[string]interface{}{
							"token": "key",
						},
					}, nil)
					ss.EXPECT().PreHookDBTransformConfigs(mock.AnythingOfType("*context.emptyCtx"), map[string]interface{}{"token": "key"}).Return(map[string]interface{}{
						"token": "encrypted_key",
					}, nil)
					rr.EXPECT().Update(mock.AnythingOfType("*context.emptyCtx"), &receiver.Receiver{
						ID:   123,
						Type: receiver.TypeSlack,
						Configurations: map[string]interface{}{
							"token": "encrypted_key",
						},
					}).Return(errors.New("some error"))
				},
				Rcv: &receiver.Receiver{
					ID:   123,
					Type: receiver.TypeSlack,
					Configurations: map[string]interface{}{
						"token": "key",
					},
				},
				Err: errors.New("some error"),
			},
			{
				Description: "should return nil error if no error returned",
				Setup: func(rr *mocks.ReceiverRepository, ss *mocks.ConfigResolver) {
					rr.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&receiver.Receiver{
						ID:   123,
						Type: receiver.TypeSlack,
						Configurations: map[string]interface{}{
							"token": "key",
						},
					}, nil)
					ss.EXPECT().PreHookDBTransformConfigs(mock.AnythingOfType("*context.emptyCtx"), map[string]interface{}{"token": "key"}).Return(map[string]interface{}{
						"token": "encrypted_key",
					}, nil)
					rr.EXPECT().Update(mock.AnythingOfType("*context.emptyCtx"), &receiver.Receiver{
						ID:   123,
						Type: receiver.TypeSlack,
						Configurations: map[string]interface{}{
							"token": "encrypted_key",
						},
					}).Return(nil)
				},
				Rcv: &receiver.Receiver{
					ID:   123,
					Type: receiver.TypeSlack,
					Configurations: map[string]interface{}{
						"token": "key",
					},
				},
				Err: nil,
			}, {
				Description: "should return error not found if repository return not found error",
				Setup: func(rr *mocks.ReceiverRepository, ss *mocks.ConfigResolver) {
					rr.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("uint64")).Return(&receiver.Receiver{
						ID:   123,
						Type: receiver.TypeSlack,
						Configurations: map[string]interface{}{
							"token": "key",
						},
					}, nil)
					ss.EXPECT().PreHookDBTransformConfigs(mock.AnythingOfType("*context.emptyCtx"), map[string]interface{}{"token": "key"}).Return(map[string]interface{}{
						"token": "encrypted_key",
					}, nil)
					rr.EXPECT().Update(mock.AnythingOfType("*context.emptyCtx"), &receiver.Receiver{
						ID:   123,
						Type: receiver.TypeSlack,
						Configurations: map[string]interface{}{
							"token": "encrypted_key",
						},
					}).Return(receiver.NotFoundError{})
				},
				Rcv: &receiver.Receiver{
					ID:   123,
					Type: receiver.TypeSlack,
					Configurations: map[string]interface{}{
						"token": "key",
					},
				},
				Err: errors.New("receiver not found"),
			},
		}
	)

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			var (
				repositoryMock = new(mocks.ReceiverRepository)
				resolverMock   = new(mocks.ConfigResolver)
			)

			registry := map[string]receiver.ConfigResolver{
				receiver.TypeSlack: resolverMock,
			}

			svc := receiver.NewService(repositoryMock, registry)

			tc.Setup(repositoryMock, resolverMock)

			err := svc.Update(ctx, tc.Rcv)
			if tc.Err != err {
				if tc.Err.Error() != err.Error() {
					t.Fatalf("got error %s, expected was %s", err.Error(), tc.Err.Error())
				}
			}
			repositoryMock.AssertExpectations(t)
			resolverMock.AssertExpectations(t)
		})
	}
}

func TestDeleteReceiver(t *testing.T) {
	ctx := context.TODO()
	receiverID := uint64(10)

	t.Run("should call repository Delete method and return nil if no error", func(t *testing.T) {
		repositoryMock := &mocks.ReceiverRepository{}
		dummyService := receiver.NewService(repositoryMock, nil)
		repositoryMock.EXPECT().Delete(mock.AnythingOfType("*context.emptyCtx"), receiverID).Return(nil).Once()
		err := dummyService.Delete(ctx, receiverID)
		assert.Nil(t, err)
		repositoryMock.AssertExpectations(t)
	})

	t.Run("should call repository Delete method and return error if any", func(t *testing.T) {
		repositoryMock := &mocks.ReceiverRepository{}
		dummyService := receiver.NewService(repositoryMock, nil)
		repositoryMock.EXPECT().Delete(mock.AnythingOfType("*context.emptyCtx"), receiverID).Return(errors.New("random error")).Once()
		err := dummyService.Delete(ctx, receiverID)
		assert.EqualError(t, err, "random error")
		repositoryMock.AssertExpectations(t)
	})
}
