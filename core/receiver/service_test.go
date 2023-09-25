package receiver_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/goto/siren/core/receiver"
	"github.com/goto/siren/core/receiver/mocks"
	"github.com/goto/siren/pkg/errors"
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
					rr.EXPECT().List(mock.AnythingOfType("context.todoCtx"), receiver.Filter{Expanded: true}).Return(nil, errors.New("some error"))
				},
				Err: errors.New("some error"),
			},
			{
				Description: "should return error if List repository success and decrypt error",
				Setup: func(rr *mocks.ReceiverRepository, ss *mocks.ConfigResolver) {
					rr.EXPECT().List(mock.AnythingOfType("context.todoCtx"), receiver.Filter{
						Expanded: true,
					}).Return([]receiver.Receiver{
						{
							ID:   10,
							Name: "foo",
							Type: receiver.TypeSlack,
							Labels: map[string]string{
								"foo": "bar",
							},
							Configurations: map[string]any{
								"token": "key",
							},
							CreatedAt: timeNow,
							UpdatedAt: timeNow,
						},
					}, nil)
					ss.EXPECT().PostHookDBTransformConfigs(mock.AnythingOfType("context.todoCtx"), map[string]any{
						"token": "key",
					}).Return(nil, errors.New("decrypt error"))
				},
				Err: errors.New("decrypt error"),
			},
			{
				Description: "should return error if type unknown",
				Setup: func(rr *mocks.ReceiverRepository, ss *mocks.ConfigResolver) {
					rr.EXPECT().List(mock.AnythingOfType("context.todoCtx"), receiver.Filter{Expanded: true}).Return([]receiver.Receiver{
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
					rr.EXPECT().List(mock.AnythingOfType("context.todoCtx"), receiver.Filter{Expanded: true}).Return([]receiver.Receiver{
						{
							ID:   10,
							Name: "foo",
							Type: receiver.TypeSlack,
							Labels: map[string]string{
								"foo": "bar",
							},
							Configurations: map[string]any{
								"token": "key",
							},
							CreatedAt: timeNow,
							UpdatedAt: timeNow,
						},
					}, nil)
					ss.EXPECT().PostHookDBTransformConfigs(mock.AnythingOfType("context.todoCtx"), map[string]any{
						"token": "key",
					}).Return(map[string]any{
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
						Configurations: map[string]any{
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

			got, err := svc.List(ctx, receiver.Filter{Expanded: true})
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
				Description: "should return error if type child but wrong parent",
				Setup: func(rr *mocks.ReceiverRepository, ss *mocks.ConfigResolver) {
					rr.EXPECT().Get(mock.AnythingOfType("context.todoCtx"), uint64(2)).Return(&receiver.Receiver{Type: receiver.TypeFile}, nil)
				},
				Rcv: &receiver.Receiver{
					ID:   123,
					Type: receiver.TypeSlackChannel,
					Configurations: map[string]any{
						"token": "key",
					},
					ParentID: 2,
				},
				Err: errors.New("parent of slack_channel type should be slack but found file"),
			},
			{
				Description: "should return error if validateParent return error",
				Setup: func(rr *mocks.ReceiverRepository, ss *mocks.ConfigResolver) {
					rr.EXPECT().Get(mock.AnythingOfType("context.todoCtx"), uint64(2)).Return(nil, errors.New("some error"))
				},
				Rcv: &receiver.Receiver{
					ID:   123,
					Type: receiver.TypeSlackChannel,
					Configurations: map[string]any{
						"token": "key",
					},
					ParentID: 2,
				},
				Err: errors.New("failed to check parent id 2"),
			},
			{
				Description: "should return error if PreHookDBTransformConfigs return error",
				Setup: func(rr *mocks.ReceiverRepository, ss *mocks.ConfigResolver) {
					ss.EXPECT().PreHookDBTransformConfigs(mock.AnythingOfType("context.todoCtx"), map[string]any{"token": "key"}).Return(nil, errors.New("some error"))

				},
				Rcv: &receiver.Receiver{
					ID:   123,
					Type: receiver.TypeSlack,
					Configurations: map[string]any{
						"token": "key",
					},
				},
				Err: errors.New("some error"),
			},
			{
				Description: "should return error if Create repository return error",
				Setup: func(rr *mocks.ReceiverRepository, ss *mocks.ConfigResolver) {
					ss.EXPECT().PreHookDBTransformConfigs(mock.AnythingOfType("context.todoCtx"), map[string]any{"token": "key"}).Return(map[string]any{
						"token": "encrypted_key",
					}, nil)
					rr.EXPECT().Create(mock.AnythingOfType("context.todoCtx"), &receiver.Receiver{
						ID:   123,
						Type: receiver.TypeSlack,
						Configurations: map[string]any{
							"token": "encrypted_key",
						},
					}).Return(errors.New("some error"))
				},
				Rcv: &receiver.Receiver{
					ID:   123,
					Type: receiver.TypeSlack,
					Configurations: map[string]any{
						"token": "key",
					},
				},
				Err: errors.New("some error"),
			},
			{
				Description: "should return nil error if no error returned",
				Setup: func(rr *mocks.ReceiverRepository, ss *mocks.ConfigResolver) {
					ss.EXPECT().PreHookDBTransformConfigs(mock.AnythingOfType("context.todoCtx"), map[string]any{"token": "key"}).Return(map[string]any{
						"token": "encrypted_key",
					}, nil)
					rr.EXPECT().Create(mock.AnythingOfType("context.todoCtx"), &receiver.Receiver{
						ID:   123,
						Type: receiver.TypeSlack,
						Configurations: map[string]any{
							"token": "encrypted_key",
						},
					}).Return(nil)
				},
				Rcv: &receiver.Receiver{
					ID:   123,
					Type: receiver.TypeSlack,
					Configurations: map[string]any{
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
					rr.EXPECT().Get(mock.AnythingOfType("context.todoCtx"), testID).Return(nil, errors.New("some error"))
				},
				Err: errors.New("some error"),
			},
			{
				Description: "should return error if type unknown",
				Setup: func(rr *mocks.ReceiverRepository, ss *mocks.ConfigResolver) {
					rr.EXPECT().Get(mock.AnythingOfType("context.todoCtx"), testID).Return(&receiver.Receiver{
						Type: "random",
					}, nil)
				},
				Err: errors.New("unsupported receiver type: \"random\""),
			},
			{
				Description: "should return error not found if Get repository return not found error",
				Setup: func(rr *mocks.ReceiverRepository, ss *mocks.ConfigResolver) {
					rr.EXPECT().Get(mock.AnythingOfType("context.todoCtx"), testID).Return(nil, receiver.NotFoundError{})
				},
				Err: errors.New("receiver not found"),
			},
			{
				Description: "should return error if Get repository success and decrypt error",
				Setup: func(rr *mocks.ReceiverRepository, ss *mocks.ConfigResolver) {
					rr.EXPECT().Get(mock.AnythingOfType("context.todoCtx"), testID).Return(&receiver.Receiver{
						ID:   10,
						Name: "foo",
						Type: receiver.TypeSlack,
						Labels: map[string]string{
							"foo": "bar",
						},
						Configurations: map[string]any{
							"token": "key",
						},
						CreatedAt: timeNow,
						UpdatedAt: timeNow,
					}, nil)
					ss.EXPECT().PostHookDBTransformConfigs(mock.AnythingOfType("context.todoCtx"), map[string]any{
						"token": "key",
					}).Return(nil, errors.New("decrypt error"))
				},
				Err: errors.New("decrypt error"),
			},
			{
				Description: "should success if Get repository and decrypt success",
				Setup: func(rr *mocks.ReceiverRepository, ss *mocks.ConfigResolver) {
					rr.EXPECT().Get(mock.AnythingOfType("context.todoCtx"), testID).Return(&receiver.Receiver{
						ID:   10,
						Name: "foo",
						Type: receiver.TypeSlack,
						Labels: map[string]string{
							"foo": "bar",
						},
						Configurations: map[string]any{
							"token": "key",
						},
						CreatedAt: timeNow,
						UpdatedAt: timeNow,
					}, nil)
					ss.EXPECT().PostHookDBTransformConfigs(mock.AnythingOfType("context.todoCtx"), map[string]any{
						"token": "key",
					}).Return(map[string]any{
						"token": "decrypted_key",
					}, nil)
					ss.EXPECT().BuildData(mock.AnythingOfType("context.todoCtx"), map[string]any{
						"token": "decrypted_key",
					}).Return(map[string]any{
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
					Data: map[string]any{
						"newdata": "populated",
					},
					Configurations: map[string]any{
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

			got, err := svc.Get(ctx, testID, receiver.GetWithData(), receiver.GetWithExpand())
			if tc.Err != err {
				if tc.Err.Error() != err.Error() {
					t.Fatalf("got error %s, expected was %s", err.Error(), tc.Err.Error())
				}
			}
			if diff := cmp.Diff(got, tc.ExpectedReceiver); diff != "" {
				t.Fatalf("got diff %+v", diff)
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
					rr.EXPECT().Get(mock.AnythingOfType("context.todoCtx"), mock.AnythingOfType("uint64")).Return(nil, errors.New("some error"))
				},
				Rcv: &receiver.Receiver{
					ID: 123,
				},
				Err: errors.New("some error"),
			},
			{
				Description: "should return error if type unknown",
				Setup: func(rr *mocks.ReceiverRepository, ss *mocks.ConfigResolver) {
					rr.EXPECT().Get(mock.AnythingOfType("context.todoCtx"), mock.AnythingOfType("uint64")).Return(&receiver.Receiver{
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
					rr.EXPECT().Get(mock.AnythingOfType("context.todoCtx"), mock.AnythingOfType("uint64")).Return(&receiver.Receiver{
						ID:   123,
						Type: receiver.TypeSlack,
						Configurations: map[string]any{
							"token": "key",
						},
					}, nil)
					ss.EXPECT().PreHookDBTransformConfigs(mock.AnythingOfType("context.todoCtx"), map[string]any{"token": "key"}).Return(nil, errors.New("some error"))
				},
				Rcv: &receiver.Receiver{
					ID:   123,
					Type: receiver.TypeSlack,
					Configurations: map[string]any{
						"token": "key",
					},
				},
				Err: errors.New("some error"),
			},
			{
				Description: "should return error if Update repository return error",
				Setup: func(rr *mocks.ReceiverRepository, ss *mocks.ConfigResolver) {
					rr.EXPECT().Get(mock.AnythingOfType("context.todoCtx"), mock.AnythingOfType("uint64")).Return(&receiver.Receiver{
						ID:   123,
						Type: receiver.TypeSlack,
						Configurations: map[string]any{
							"token": "key",
						},
					}, nil)
					ss.EXPECT().PreHookDBTransformConfigs(mock.AnythingOfType("context.todoCtx"), map[string]any{"token": "key"}).Return(map[string]any{
						"token": "encrypted_key",
					}, nil)
					rr.EXPECT().Update(mock.AnythingOfType("context.todoCtx"), &receiver.Receiver{
						ID:   123,
						Type: receiver.TypeSlack,
						Configurations: map[string]any{
							"token": "encrypted_key",
						},
					}).Return(errors.New("some error"))
				},
				Rcv: &receiver.Receiver{
					ID:   123,
					Type: receiver.TypeSlack,
					Configurations: map[string]any{
						"token": "key",
					},
				},
				Err: errors.New("some error"),
			},
			{
				Description: "should return nil error if no error returned",
				Setup: func(rr *mocks.ReceiverRepository, ss *mocks.ConfigResolver) {
					rr.EXPECT().Get(mock.AnythingOfType("context.todoCtx"), mock.AnythingOfType("uint64")).Return(&receiver.Receiver{
						ID:   123,
						Type: receiver.TypeSlack,
						Configurations: map[string]any{
							"token": "old-key",
						},
					}, nil)
					ss.EXPECT().PreHookDBTransformConfigs(mock.AnythingOfType("context.todoCtx"), map[string]any{"token": "key"}).Return(map[string]any{
						"token": "encrypted_key",
					}, nil)
					rr.EXPECT().Update(mock.AnythingOfType("context.todoCtx"), &receiver.Receiver{
						ID:   123,
						Type: receiver.TypeSlack,
						Configurations: map[string]any{
							"token": "encrypted_key",
						},
					}).Return(nil)
				},
				Rcv: &receiver.Receiver{
					ID:   123,
					Type: receiver.TypeSlack,
					Configurations: map[string]any{
						"token": "key",
					},
				},
				Err: nil,
			}, {
				Description: "should return error not found if repository return not found error",
				Setup: func(rr *mocks.ReceiverRepository, ss *mocks.ConfigResolver) {
					rr.EXPECT().Get(mock.AnythingOfType("context.todoCtx"), mock.AnythingOfType("uint64")).Return(&receiver.Receiver{
						ID:   123,
						Type: receiver.TypeSlack,
						Configurations: map[string]any{
							"token": "key",
						},
					}, nil)
					ss.EXPECT().PreHookDBTransformConfigs(mock.AnythingOfType("context.todoCtx"), map[string]any{"token": "key"}).Return(map[string]any{
						"token": "encrypted_key",
					}, nil)
					rr.EXPECT().Update(mock.AnythingOfType("context.todoCtx"), &receiver.Receiver{
						ID:   123,
						Type: receiver.TypeSlack,
						Configurations: map[string]any{
							"token": "encrypted_key",
						},
					}).Return(receiver.NotFoundError{})
				},
				Rcv: &receiver.Receiver{
					ID:   123,
					Type: receiver.TypeSlack,
					Configurations: map[string]any{
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
		repositoryMock.EXPECT().Delete(mock.AnythingOfType("context.todoCtx"), receiverID).Return(nil).Once()
		err := dummyService.Delete(ctx, receiverID)
		assert.Nil(t, err)
		repositoryMock.AssertExpectations(t)
	})

	t.Run("should call repository Delete method and return error if any", func(t *testing.T) {
		repositoryMock := &mocks.ReceiverRepository{}
		dummyService := receiver.NewService(repositoryMock, nil)
		repositoryMock.EXPECT().Delete(mock.AnythingOfType("context.todoCtx"), receiverID).Return(errors.New("random error")).Once()
		err := dummyService.Delete(ctx, receiverID)
		assert.EqualError(t, err, "random error")
		repositoryMock.AssertExpectations(t)
	})
}

func TestService_ExpandParents(t *testing.T) {
	type testCase struct {
		Description       string
		Receivers         []receiver.Receiver
		ExpandedReceivers []receiver.Receiver
		Setup             func(*mocks.ReceiverRepository, *mocks.ConfigResolver, *mocks.ConfigResolver)
		Err               error
	}

	var (
		ctx             = context.TODO()
		sampleReceivers = []receiver.Receiver{
			{
				ID: 3,
				Configurations: map[string]any{
					"channel_name": "my-channel",
				},
				ParentID: 1,
			},
			{
				ID: 4,
				Configurations: map[string]any{
					"channel_name": "any-channel",
				},
				ParentID: 2,
			},
			{
				ID: 5,
				Configurations: map[string]any{
					"channel_name": "my-channel",
				},
				ParentID: 2,
			},
		}
		timeNow   = time.Now()
		testCases = []testCase{
			{
				Description: "should return empty expanded receivers if no receivers passed",
				Setup:       func(rr *mocks.ReceiverRepository, _ *mocks.ConfigResolver, _ *mocks.ConfigResolver) {},
				Err:         nil,
			},
			{
				Description: "should return as-is receivers if no child receivers found",
				Setup:       func(rr *mocks.ReceiverRepository, _ *mocks.ConfigResolver, _ *mocks.ConfigResolver) {},
				Receivers: []receiver.Receiver{
					{
						ID: 1,
					},
				},
				ExpandedReceivers: []receiver.Receiver{
					{
						ID: 1,
					},
				},
				Err: nil,
			},
			{
				Description: "should return expanded receivers if there is any child receivers found",
				Setup: func(rr *mocks.ReceiverRepository, ss *mocks.ConfigResolver, sc *mocks.ConfigResolver) {
					rr.EXPECT().List(mock.AnythingOfType("context.todoCtx"), receiver.Filter{
						ReceiverIDs: []uint64{1, 2},
						Expanded:    false,
					}).Return([]receiver.Receiver{
						{
							ID:   1,
							Name: "foo",
							Type: receiver.TypeSlack,
							Labels: map[string]string{
								"foo": "bar",
							},
							Configurations: map[string]any{
								"token":     "key1",
								"workspace": "company1",
							},
							CreatedAt: timeNow,
							UpdatedAt: timeNow,
						},
						{
							ID:   2,
							Name: "bar",
							Type: receiver.TypeSlack,
							Labels: map[string]string{
								"foo": "bar",
							},
							Configurations: map[string]any{
								"token":     "key2",
								"workspace": "company2",
							},
							CreatedAt: timeNow,
							UpdatedAt: timeNow,
						},
					}, nil)
				},
				Receivers: sampleReceivers,
				ExpandedReceivers: []receiver.Receiver{
					{
						ID: 3,
						Configurations: map[string]any{
							"channel_name": "my-channel",
							"token":        "key1",
							"workspace":    "company1",
						},
						ParentID: 1,
					},
					{
						ID: 4,
						Configurations: map[string]any{
							"channel_name": "any-channel",
							"token":        "key2",
							"workspace":    "company2",
						},
						ParentID: 2,
					},
					{
						ID: 5,
						Configurations: map[string]any{
							"channel_name": "my-channel",
							"token":        "key2",
							"workspace":    "company2",
						},
						ParentID: 2,
					},
				},
				Err: nil,
			},
		}
	)

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			var (
				repositoryMock           = new(mocks.ReceiverRepository)
				resolverSlackMock        = new(mocks.ConfigResolver)
				resolverSlackChannelMock = new(mocks.ConfigResolver)
			)

			registry := map[string]receiver.ConfigResolver{
				receiver.TypeSlackChannel: resolverSlackChannelMock,
				receiver.TypeSlack:        resolverSlackMock,
			}

			svc := receiver.NewService(repositoryMock, registry)

			tc.Setup(repositoryMock, resolverSlackMock, resolverSlackChannelMock)

			got, err := svc.ExpandParents(ctx, tc.Receivers)
			if tc.Err != err {
				if tc.Err.Error() != err.Error() {
					t.Fatalf("got error %s, expected was %s", err.Error(), tc.Err.Error())
				}
			}
			if diff := cmp.Diff(got, tc.ExpandedReceivers); diff != "" {
				t.Fatalf("got diff %+v", diff)
			}
			repositoryMock.AssertExpectations(t)
			resolverSlackMock.AssertExpectations(t)
			resolverSlackChannelMock.AssertExpectations(t)
		})
	}
}
