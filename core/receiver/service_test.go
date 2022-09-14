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
		Setup       func(*mocks.ReceiverRepository, *mocks.ReceiverPlugin)
		Err         error
	}

	var (
		ctx       = context.TODO()
		timeNow   = time.Now()
		testCases = []testCase{
			{
				Description: "should return error if List repository error",
				Setup: func(rr *mocks.ReceiverRepository, ss *mocks.ReceiverPlugin) {
					rr.EXPECT().List(mock.AnythingOfType("*context.emptyCtx"), receiver.Filter{}).Return(nil, errors.New("some error"))
				},
				Err: errors.New("some error"),
			},
			{
				Description: "should return error if List repository success and decrypt error",
				Setup: func(rr *mocks.ReceiverRepository, ss *mocks.ReceiverPlugin) {
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
					ss.EXPECT().PostHookTransformConfigs(mock.AnythingOfType("*context.emptyCtx"), receiver.Configurations{
						"token": "key",
					}).Return(nil, errors.New("decrypt error"))
				},
				Err: errors.New("decrypt error"),
			},
			{
				Description: "should return error if type unknown",
				Setup: func(rr *mocks.ReceiverRepository, ss *mocks.ReceiverPlugin) {
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
				Setup: func(rr *mocks.ReceiverRepository, ss *mocks.ReceiverPlugin) {
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
					ss.EXPECT().PostHookTransformConfigs(mock.AnythingOfType("*context.emptyCtx"), receiver.Configurations{
						"token": "key",
					}).Return(receiver.Configurations{
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
				repositoryMock     = new(mocks.ReceiverRepository)
				receiverPluginMock = new(mocks.ReceiverPlugin)
			)

			registry := map[string]receiver.ReceiverPlugin{
				receiver.TypeSlack: receiverPluginMock,
			}

			svc := receiver.NewService(repositoryMock, registry)

			tc.Setup(repositoryMock, receiverPluginMock)

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
			receiverPluginMock.AssertExpectations(t)
		})
	}
}

func TestService_CreateReceiver(t *testing.T) {
	type testCase struct {
		Description string
		Setup       func(*mocks.ReceiverRepository, *mocks.ReceiverPlugin)
		Rcv         *receiver.Receiver
		Err         error
	}
	var (
		ctx       = context.TODO()
		testCases = []testCase{
			{
				Description: "should return error if configuration is not valid",
				Setup: func(rr *mocks.ReceiverRepository, ss *mocks.ReceiverPlugin) {
					ss.EXPECT().ValidateConfigurations(receiver.Configurations{
						"token": "key",
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
				Description: "should return error if encrypt return error",
				Setup: func(rr *mocks.ReceiverRepository, ss *mocks.ReceiverPlugin) {
					ss.EXPECT().ValidateConfigurations(receiver.Configurations{
						"token": "key",
					}).Return(nil)
					ss.EXPECT().PreHookTransformConfigs(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("receiver.Configurations")).Return(nil, errors.New("some error"))
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
				Description: "should return error if type unknown",
				Setup:       func(rr *mocks.ReceiverRepository, ss *mocks.ReceiverPlugin) {},
				Rcv: &receiver.Receiver{
					Type: "random",
				},
				Err: errors.New("unsupported receiver type: \"random\""),
			},
			{
				Description: "should return error if Create repository return error",
				Setup: func(rr *mocks.ReceiverRepository, ss *mocks.ReceiverPlugin) {
					ss.EXPECT().ValidateConfigurations(receiver.Configurations{
						"token": "key",
					}).Return(nil)
					ss.EXPECT().PreHookTransformConfigs(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("receiver.Configurations")).Return(receiver.Configurations{
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
				Setup: func(rr *mocks.ReceiverRepository, ss *mocks.ReceiverPlugin) {
					ss.EXPECT().ValidateConfigurations(receiver.Configurations{
						"token": "key",
					}).Return(nil)
					ss.EXPECT().PreHookTransformConfigs(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("receiver.Configurations")).Return(receiver.Configurations{
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
				repositoryMock     = new(mocks.ReceiverRepository)
				ReceiverPluginMock = new(mocks.ReceiverPlugin)
			)

			registry := map[string]receiver.ReceiverPlugin{
				receiver.TypeSlack: ReceiverPluginMock,
			}

			svc := receiver.NewService(repositoryMock, registry)

			tc.Setup(repositoryMock, ReceiverPluginMock)

			err := svc.Create(ctx, tc.Rcv)
			if tc.Err != err {
				if tc.Err.Error() != err.Error() {
					t.Fatalf("got error %s, expected was %s", err.Error(), tc.Err.Error())
				}
			}
			repositoryMock.AssertExpectations(t)
			ReceiverPluginMock.AssertExpectations(t)
		})
	}
}

func TestService_GetReceiver(t *testing.T) {
	type testCase struct {
		Description      string
		ExpectedReceiver *receiver.Receiver
		Setup            func(*mocks.ReceiverRepository, *mocks.ReceiverPlugin)
		Err              error
	}

	var (
		ctx       = context.TODO()
		timeNow   = time.Now()
		testID    = uint64(10)
		testCases = []testCase{
			{
				Description: "should return error if Get repository error",
				Setup: func(rr *mocks.ReceiverRepository, ss *mocks.ReceiverPlugin) {
					rr.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), testID).Return(nil, errors.New("some error"))
				},
				Err: errors.New("some error"),
			},
			{
				Description: "should return error if type unknown",
				Setup: func(rr *mocks.ReceiverRepository, ss *mocks.ReceiverPlugin) {
					rr.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), testID).Return(&receiver.Receiver{
						Type: "random",
					}, nil)
				},
				Err: errors.New("unsupported receiver type: \"random\""),
			},
			{
				Description: "should return error not found if Get repository return not found error",
				Setup: func(rr *mocks.ReceiverRepository, ss *mocks.ReceiverPlugin) {
					rr.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), testID).Return(nil, receiver.NotFoundError{})
				},
				Err: errors.New("receiver not found"),
			},
			{
				Description: "should return error if Get repository success and decrypt error",
				Setup: func(rr *mocks.ReceiverRepository, ss *mocks.ReceiverPlugin) {
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
					ss.EXPECT().PostHookTransformConfigs(mock.AnythingOfType("*context.emptyCtx"), receiver.Configurations{
						"token": "key",
					}).Return(nil, errors.New("decrypt error"))
				},
				Err: errors.New("decrypt error"),
			},
			{
				Description: "should success if Get repository and decrypt success",
				Setup: func(rr *mocks.ReceiverRepository, ss *mocks.ReceiverPlugin) {
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
					ss.EXPECT().PostHookTransformConfigs(mock.AnythingOfType("*context.emptyCtx"), receiver.Configurations{
						"token": "key",
					}).Return(receiver.Configurations{
						"token": "decrypted_key",
					}, nil)
					ss.EXPECT().PopulateDataFromConfigs(mock.AnythingOfType("*context.emptyCtx"), receiver.Configurations{
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
				repositoryMock     = new(mocks.ReceiverRepository)
				ReceiverPluginMock = new(mocks.ReceiverPlugin)
			)

			registry := map[string]receiver.ReceiverPlugin{
				receiver.TypeSlack: ReceiverPluginMock,
			}

			svc := receiver.NewService(repositoryMock, registry)

			tc.Setup(repositoryMock, ReceiverPluginMock)

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
			ReceiverPluginMock.AssertExpectations(t)
		})
	}
}

func TestService_UpdateReceiver(t *testing.T) {
	type testCase struct {
		Description string
		Setup       func(*mocks.ReceiverRepository, *mocks.ReceiverPlugin)
		Rcv         *receiver.Receiver
		Err         error
	}
	var (
		ctx       = context.TODO()
		testCases = []testCase{
			{
				Description: "should return error if validate configuration return error",
				Setup: func(rr *mocks.ReceiverRepository, ss *mocks.ReceiverPlugin) {
					ss.EXPECT().ValidateConfigurations(receiver.Configurations{
						"token": "key",
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
				Description: "should return error if encrypt return error",
				Setup: func(rr *mocks.ReceiverRepository, ss *mocks.ReceiverPlugin) {
					ss.EXPECT().ValidateConfigurations(receiver.Configurations{
						"token": "key",
					}).Return(nil)
					ss.EXPECT().PreHookTransformConfigs(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("receiver.Configurations")).Return(nil, errors.New("some error"))
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
				Description: "should return error if type unknown",
				Setup:       func(rr *mocks.ReceiverRepository, ss *mocks.ReceiverPlugin) {},
				Rcv: &receiver.Receiver{
					Type: "random",
				},
				Err: errors.New("unsupported receiver type: \"random\""),
			},
			{
				Description: "should return error if Update repository return error",
				Setup: func(rr *mocks.ReceiverRepository, ss *mocks.ReceiverPlugin) {
					ss.EXPECT().ValidateConfigurations(receiver.Configurations{
						"token": "key",
					}).Return(nil)
					ss.EXPECT().PreHookTransformConfigs(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("receiver.Configurations")).Return(receiver.Configurations{
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
				Setup: func(rr *mocks.ReceiverRepository, ss *mocks.ReceiverPlugin) {
					ss.EXPECT().ValidateConfigurations(receiver.Configurations{
						"token": "key",
					}).Return(nil)
					ss.EXPECT().PreHookTransformConfigs(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("receiver.Configurations")).Return(receiver.Configurations{
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
				Setup: func(rr *mocks.ReceiverRepository, ss *mocks.ReceiverPlugin) {
					ss.EXPECT().ValidateConfigurations(receiver.Configurations{
						"token": "key",
					}).Return(nil)
					ss.EXPECT().PreHookTransformConfigs(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("receiver.Configurations")).Return(receiver.Configurations{
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
				repositoryMock     = new(mocks.ReceiverRepository)
				ReceiverPluginMock = new(mocks.ReceiverPlugin)
			)

			registry := map[string]receiver.ReceiverPlugin{
				receiver.TypeSlack: ReceiverPluginMock,
			}

			svc := receiver.NewService(repositoryMock, registry)

			tc.Setup(repositoryMock, ReceiverPluginMock)

			err := svc.Update(ctx, tc.Rcv)
			if tc.Err != err {
				if tc.Err.Error() != err.Error() {
					t.Fatalf("got error %s, expected was %s", err.Error(), tc.Err.Error())
				}
			}
			repositoryMock.AssertExpectations(t)
			ReceiverPluginMock.AssertExpectations(t)
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

func TestService_Notify(t *testing.T) {
	type testCase struct {
		Description string
		Setup       func(*mocks.ReceiverRepository, *mocks.ReceiverPlugin)
		Receiver    *receiver.Receiver
		ErrString   string
	}
	var (
		receiverID = uint64(50)
		testCases  = []testCase{
			{
				Description: "should return error if get repository return error",
				Setup: func(rr *mocks.ReceiverRepository, ts *mocks.ReceiverPlugin) {
					rr.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), receiverID).Return(nil, errors.New("some error"))
				},
				ErrString: "error getting receiver with id 50",
			},
			{
				Description: "should return error if type from get repository unknown",
				Setup: func(rr *mocks.ReceiverRepository, ts *mocks.ReceiverPlugin) {
					rr.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), receiverID).Return(&receiver.Receiver{
						Type: "random",
					}, nil)
				},
				ErrString: "error getting receiver with id 50",
			},
			{
				Description: "should return error if type from get service unknown",
				Setup: func(rr *mocks.ReceiverRepository, ts *mocks.ReceiverPlugin) {
					rr.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), receiverID).Return(&receiver.Receiver{
						Type: "random",
					}, nil)
				},
				ErrString: "error getting receiver with id 50",
			},
			{
				Description: "should return error if type service notify return error",
				Setup: func(rr *mocks.ReceiverRepository, ts *mocks.ReceiverPlugin) {
					rr.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), receiverID).Return(&receiver.Receiver{
						Type: receiver.TypeSlack,
					}, nil)
					ts.EXPECT().PostHookTransformConfigs(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("receiver.Configurations")).Return(receiver.Configurations{
						"token": "decrypted_key",
					}, nil)
					ts.EXPECT().PopulateDataFromConfigs(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("receiver.Configurations")).Return(map[string]interface{}{}, nil)
					ts.EXPECT().Notify(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("receiver.Configurations"), mock.AnythingOfType("map[string]interface {}")).Return(errors.New("some error"))
				},
				ErrString: "some error",
			},
			{
				Description: "should return no error if type service notify return nil error",
				Setup: func(rr *mocks.ReceiverRepository, ts *mocks.ReceiverPlugin) {
					rr.EXPECT().Get(mock.AnythingOfType("*context.emptyCtx"), receiverID).Return(&receiver.Receiver{
						Type: receiver.TypeSlack,
					}, nil)
					ts.EXPECT().PostHookTransformConfigs(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("receiver.Configurations")).Return(receiver.Configurations{
						"token": "decrypted_key",
					}, nil)
					ts.EXPECT().PopulateDataFromConfigs(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("receiver.Configurations")).Return(map[string]interface{}{}, nil)
					ts.EXPECT().Notify(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("receiver.Configurations"), mock.AnythingOfType("map[string]interface {}")).Return(nil)
				},
			},
		}
	)

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			var (
				repositoryMock     = new(mocks.ReceiverRepository)
				ReceiverPluginMock = new(mocks.ReceiverPlugin)
			)

			registry := map[string]receiver.ReceiverPlugin{
				receiver.TypeSlack: ReceiverPluginMock,
			}

			svc := receiver.NewService(repositoryMock, registry)

			tc.Setup(repositoryMock, ReceiverPluginMock)

			err := svc.Notify(context.TODO(), receiverID, map[string]interface{}{})
			if tc.ErrString != "" {
				if tc.ErrString != err.Error() {
					t.Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}

			repositoryMock.AssertExpectations(t)
			ReceiverPluginMock.AssertExpectations(t)
		})
	}
}

// func TestService_GetSubscriptionConfig(t *testing.T) {
// 	type testCase struct {
// 		Description string
// 		Setup       func(*mocks.ReceiverPlugin)
// 		Receiver    *receiver.Receiver
// 		ErrString   string
// 	}
// 	var (
// 		testCases = []testCase{
// 			{
// 				Description: "should return error if receiver is nil",
// 				Setup:       func(ts *mocks.ReceiverPlugin) {},
// 				ErrString:   "request is not valid",
// 			},
// 			{
// 				Description: "should return error if type unknown",
// 				Setup:       func(ts *mocks.ReceiverPlugin) {},
// 				Receiver: &receiver.Receiver{
// 					Type: "random",
// 				},
// 				ErrString: "unsupported receiver type: \"random\"",
// 			},
// 			{
// 				Description: "should return error if type receiver get subscription config is returning error",
// 				Setup: func(ts *mocks.ReceiverPlugin) {
// 					ts.EXPECT().GetSubscriptionConfig(mock.AnythingOfType("map[string]string"), mock.AnythingOfType("map[string]interface{}")).Return(nil, errors.New("some error"))
// 				},
// 				Receiver: &receiver.Receiver{
// 					Type: receiver.TypeSlack,
// 				},
// 				ErrString: "some error",
// 			},
// 			{
// 				Description: "should return no error if type receiver get subscription config is returning no error",
// 				Setup: func(ts *mocks.ReceiverPlugin) {
// 					ts.EXPECT().GetSubscriptionConfig(mock.AnythingOfType("map[string]string"), mock.AnythingOfType("map[string]interface{}")).Return(map[string]string{}, nil)
// 				},
// 				Receiver: &receiver.Receiver{
// 					Type: receiver.TypeSlack,
// 				},
// 			},
// 		}
// 	)

// 	for _, tc := range testCases {
// 		t.Run(tc.Description, func(t *testing.T) {
// 			var (
// 				ReceiverPluginMock = new(mocks.ReceiverPlugin)
// 			)

// 			registry := map[string]receiver.ReceiverPlugin{
// 				receiver.TypeSlack: ReceiverPluginMock,
// 			}

// 			svc := receiver.NewService(nil, registry)

// 			tc.Setup(ReceiverPluginMock)

// 			_, err := svc.GetSubscriptionConfig(map[string]string{}, tc.Receiver)
// 			if tc.ErrString != "" {
// 				if tc.ErrString != err.Error() {
// 					t.Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
// 				}
// 			}

// 			ReceiverPluginMock.AssertExpectations(t)
// 		})
// 	}
// }
