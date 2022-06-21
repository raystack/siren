package receiver_test

import (
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/odpf/siren/core/receiver"
	"github.com/odpf/siren/core/receiver/mocks"
	"github.com/stretchr/testify/mock"
)

func TestService_ListReceivers(t *testing.T) {
	type testCase struct {
		Description string
		Receivers   []*receiver.Receiver
		Setup       func(*mocks.ReceiverRepository, *mocks.TypeService)
		Err         error
	}

	var (
		timeNow   = time.Now()
		testCases = []testCase{
			{
				Description: "should return error if List repository error",
				Setup: func(rr *mocks.ReceiverRepository, ss *mocks.TypeService) {
					rr.EXPECT().List().Return(nil, errors.New("some error"))
				},
				Err: errors.New("secureService.repository.List: some error"),
			},
			{
				Description: "should return error if List repository success and decrypt error",
				Setup: func(rr *mocks.ReceiverRepository, ss *mocks.TypeService) {
					rr.EXPECT().List().Return([]*receiver.Receiver{
						{
							ID:   10,
							Name: "foo",
							Type: "slack",
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
					ss.EXPECT().Decrypt(&receiver.Receiver{
						ID:   10,
						Name: "foo",
						Type: "slack",
						Labels: map[string]string{
							"foo": "bar",
						},
						Configurations: map[string]interface{}{
							"token": "key",
						},
						CreatedAt: timeNow,
						UpdatedAt: timeNow,
					}).Return(errors.New("decrypt error"))
				},
				Err: errors.New("decrypt error"),
			},
			{
				Description: "should success if list repository and decrypt success",
				Setup: func(rr *mocks.ReceiverRepository, ss *mocks.TypeService) {
					rr.EXPECT().List().Return([]*receiver.Receiver{
						{
							ID:   10,
							Name: "foo",
							Type: "slack",
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
					ss.EXPECT().Decrypt(&receiver.Receiver{
						ID:   10,
						Name: "foo",
						Type: "slack",
						Labels: map[string]string{
							"foo": "bar",
						},
						Configurations: map[string]interface{}{
							"token": "key",
						},
						CreatedAt: timeNow,
						UpdatedAt: timeNow,
					}).Return(nil)
				},
				Receivers: []*receiver.Receiver{
					{
						ID:   10,
						Name: "foo",
						Type: "slack",
						Labels: map[string]string{
							"foo": "bar",
						},
						Configurations: map[string]interface{}{
							"token": "key",
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
				repositoryMock      = new(mocks.ReceiverRepository)
				strategyServiceMock = new(mocks.TypeService)
			)

			registry := map[string]receiver.TypeService{
				receiver.TypeSlack: strategyServiceMock,
			}

			svc := receiver.NewService(repositoryMock, registry)

			tc.Setup(repositoryMock, strategyServiceMock)

			got, err := svc.ListReceivers()
			if tc.Err != err {
				if tc.Err.Error() != err.Error() {
					t.Fatalf("got error %s, expected was %s", err.Error(), tc.Err.Error())
				}
			}
			if !cmp.Equal(got, tc.Receivers) {
				t.Fatalf("got result %+v, expected was %+v", got, tc.Receivers)
			}
			repositoryMock.AssertExpectations(t)
			strategyServiceMock.AssertExpectations(t)
		})
	}
}

func TestService_CreateReceiver(t *testing.T) {
	type testCase struct {
		Description string
		Setup       func(*mocks.ReceiverRepository, *mocks.TypeService)
		Rcv         *receiver.Receiver
		Err         error
	}
	var testCases = []testCase{
		{
			Description: "should return error if configuration is not valid",
			Setup: func(rr *mocks.ReceiverRepository, ss *mocks.TypeService) {
				ss.EXPECT().ValidateConfiguration(receiver.Configurations{
					"token": "key",
				}).Return(errors.New("some error"))
			},
			Rcv: &receiver.Receiver{
				ID:   123,
				Type: "slack",
				Configurations: map[string]interface{}{
					"token": "key",
				},
			},
			Err: errors.New("bad_request: invalid receiver configurations"),
		},
		{
			Description: "should return error if encrypt return error",
			Setup: func(rr *mocks.ReceiverRepository, ss *mocks.TypeService) {
				ss.EXPECT().ValidateConfiguration(receiver.Configurations{
					"token": "key",
				}).Return(nil)
				ss.EXPECT().Encrypt(mock.AnythingOfType("*receiver.Receiver")).Return(errors.New("some error"))
			},
			Rcv: &receiver.Receiver{
				ID:   123,
				Type: "slack",
				Configurations: map[string]interface{}{
					"token": "key",
				},
			},
			Err: errors.New("some error"),
		},
		{
			Description: "should return error if type unknown",
			Setup:       func(rr *mocks.ReceiverRepository, ss *mocks.TypeService) {},
			Rcv: &receiver.Receiver{
				Type: "random",
			},
			Err: errors.New("bad_request: unsupported receiver type"),
		},
		{
			Description: "should return error if Create repository return error",
			Setup: func(rr *mocks.ReceiverRepository, ss *mocks.TypeService) {
				ss.EXPECT().ValidateConfiguration(receiver.Configurations{
					"token": "key",
				}).Return(nil)
				ss.EXPECT().Encrypt(mock.AnythingOfType("*receiver.Receiver")).Return(nil)
				rr.EXPECT().Create(&receiver.Receiver{
					ID:   123,
					Type: "slack",
					Configurations: map[string]interface{}{
						"token": "key",
					},
				}).Return(errors.New("some error"))
			},
			Rcv: &receiver.Receiver{
				ID:   123,
				Type: "slack",
				Configurations: map[string]interface{}{
					"token": "key",
				},
			},
			Err: errors.New("secureService.repository.Create: some error"),
		},
		{
			Description: "should return error if decrypt return error",
			Setup: func(rr *mocks.ReceiverRepository, ss *mocks.TypeService) {
				ss.EXPECT().ValidateConfiguration(receiver.Configurations{
					"token": "key",
				}).Return(nil)
				ss.EXPECT().Encrypt(mock.AnythingOfType("*receiver.Receiver")).Return(nil)
				rr.EXPECT().Create(&receiver.Receiver{
					ID:   123,
					Type: "slack",
					Configurations: map[string]interface{}{
						"token": "key",
					},
				}).Return(nil)
				ss.EXPECT().Decrypt(mock.AnythingOfType("*receiver.Receiver")).Return(errors.New("some error"))
			},
			Rcv: &receiver.Receiver{
				ID:   123,
				Type: "slack",
				Configurations: map[string]interface{}{
					"token": "key",
				},
			},
			Err: errors.New("some error"),
		},
		{
			Description: "should return nil error if no error returned",
			Setup: func(rr *mocks.ReceiverRepository, ss *mocks.TypeService) {
				ss.EXPECT().ValidateConfiguration(receiver.Configurations{
					"token": "key",
				}).Return(nil)
				ss.EXPECT().Encrypt(mock.AnythingOfType("*receiver.Receiver")).Return(nil)
				rr.EXPECT().Create(&receiver.Receiver{
					ID:   123,
					Type: "slack",
					Configurations: map[string]interface{}{
						"token": "key",
					},
				}).Return(nil)
				ss.EXPECT().Decrypt(mock.AnythingOfType("*receiver.Receiver")).Return(nil)
			},
			Rcv: &receiver.Receiver{
				ID:   123,
				Type: "slack",
				Configurations: map[string]interface{}{
					"token": "key",
				},
			},
			Err: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			var (
				repositoryMock      = new(mocks.ReceiverRepository)
				strategyServiceMock = new(mocks.TypeService)
			)

			registry := map[string]receiver.TypeService{
				receiver.TypeSlack: strategyServiceMock,
			}

			svc := receiver.NewService(repositoryMock, registry)

			tc.Setup(repositoryMock, strategyServiceMock)

			err := svc.CreateReceiver(tc.Rcv)
			if tc.Err != err {
				if tc.Err.Error() != err.Error() {
					t.Fatalf("got error %s, expected was %s", err.Error(), tc.Err.Error())
				}
			}
			repositoryMock.AssertExpectations(t)
			strategyServiceMock.AssertExpectations(t)
		})
	}
}

func TestService_GetReceiver(t *testing.T) {
	type testCase struct {
		Description string
		Rcv         *receiver.Receiver
		Setup       func(*mocks.ReceiverRepository, *mocks.TypeService)
		Err         error
	}

	var (
		timeNow   = time.Now()
		testID    = uint64(10)
		testCases = []testCase{
			{
				Description: "should return error if Get repository error",
				Setup: func(rr *mocks.ReceiverRepository, ss *mocks.TypeService) {
					rr.EXPECT().Get(testID).Return(nil, errors.New("some error"))
				},
				Err: errors.New("secureService.repository.Get: some error"),
			},
			{
				Description: "should return error if type unknown",
				Setup: func(rr *mocks.ReceiverRepository, ss *mocks.TypeService) {
					rr.EXPECT().Get(testID).Return(&receiver.Receiver{
						Type: "random",
					}, nil)
				},
				Err: errors.New("bad_request: unsupported receiver type"),
			},
			{
				Description: "should return error if Get repository success and decrypt error",
				Setup: func(rr *mocks.ReceiverRepository, ss *mocks.TypeService) {
					rr.EXPECT().Get(testID).Return(&receiver.Receiver{
						ID:   10,
						Name: "foo",
						Type: "slack",
						Labels: map[string]string{
							"foo": "bar",
						},
						Configurations: map[string]interface{}{
							"token": "key",
						},
						CreatedAt: timeNow,
						UpdatedAt: timeNow,
					}, nil)
					ss.EXPECT().Decrypt(&receiver.Receiver{
						ID:   10,
						Name: "foo",
						Type: "slack",
						Labels: map[string]string{
							"foo": "bar",
						},
						Configurations: map[string]interface{}{
							"token": "key",
						},
						CreatedAt: timeNow,
						UpdatedAt: timeNow,
					}).Return(errors.New("decrypt error"))
				},
				Err: errors.New("decrypt error"),
			},
			{
				Description: "should success if Get repository and decrypt success",
				Setup: func(rr *mocks.ReceiverRepository, ss *mocks.TypeService) {
					rr.EXPECT().Get(testID).Return(&receiver.Receiver{
						ID:   10,
						Name: "foo",
						Type: "slack",
						Labels: map[string]string{
							"foo": "bar",
						},
						Configurations: map[string]interface{}{
							"token": "key",
						},
						CreatedAt: timeNow,
						UpdatedAt: timeNow,
					}, nil)
					ss.EXPECT().Decrypt(&receiver.Receiver{
						ID:   10,
						Name: "foo",
						Type: "slack",
						Labels: map[string]string{
							"foo": "bar",
						},
						Configurations: map[string]interface{}{
							"token": "key",
						},
						CreatedAt: timeNow,
						UpdatedAt: timeNow,
					}).Return(nil)
					ss.EXPECT().PopulateReceiver(&receiver.Receiver{
						ID:   10,
						Name: "foo",
						Type: "slack",
						Labels: map[string]string{
							"foo": "bar",
						},
						Configurations: map[string]interface{}{
							"token": "key",
						},
						CreatedAt: timeNow,
						UpdatedAt: timeNow,
					}).Return(&receiver.Receiver{
						ID:   10,
						Name: "foo",
						Type: "slack",
						Labels: map[string]string{
							"foo": "bar",
						},
						Data: map[string]interface{}{
							"newdata": "populated",
						},
						Configurations: map[string]interface{}{
							"token": "key",
						},
						CreatedAt: timeNow,
						UpdatedAt: timeNow,
					}, nil)
				},
				Rcv: &receiver.Receiver{
					ID:   10,
					Name: "foo",
					Type: "slack",
					Labels: map[string]string{
						"foo": "bar",
					},
					Data: map[string]interface{}{
						"newdata": "populated",
					},
					Configurations: map[string]interface{}{
						"token": "key",
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
				repositoryMock      = new(mocks.ReceiverRepository)
				strategyServiceMock = new(mocks.TypeService)
			)

			registry := map[string]receiver.TypeService{
				receiver.TypeSlack: strategyServiceMock,
			}

			svc := receiver.NewService(repositoryMock, registry)

			tc.Setup(repositoryMock, strategyServiceMock)

			got, err := svc.GetReceiver(testID)
			if tc.Err != err {
				if tc.Err.Error() != err.Error() {
					t.Fatalf("got error %s, expected was %s", err.Error(), tc.Err.Error())
				}
			}
			if !cmp.Equal(got, tc.Rcv) {
				t.Fatalf("got result %+v, expected was %+v", got, tc.Rcv)
			}
			repositoryMock.AssertExpectations(t)
			strategyServiceMock.AssertExpectations(t)
		})
	}
}

func TestService_UpdateReceiver(t *testing.T) {
	type testCase struct {
		Description string
		Setup       func(*mocks.ReceiverRepository, *mocks.TypeService)
		Rcv         *receiver.Receiver
		Err         error
	}
	var testCases = []testCase{
		{
			Description: "should return error if encrypt return error",
			Setup: func(rr *mocks.ReceiverRepository, ss *mocks.TypeService) {
				ss.EXPECT().ValidateConfiguration(receiver.Configurations{
					"token": "key",
				}).Return(errors.New("some error"))
			},
			Rcv: &receiver.Receiver{
				ID:   123,
				Type: "slack",
				Configurations: map[string]interface{}{
					"token": "key",
				},
			},
			Err: errors.New("bad_request: invalid receiver configurations"),
		},
		{
			Description: "should return error if type unknown",
			Setup:       func(rr *mocks.ReceiverRepository, ss *mocks.TypeService) {},
			Rcv: &receiver.Receiver{
				Type: "random",
			},
			Err: errors.New("bad_request: unsupported receiver type"),
		},
		{
			Description: "should return error if Update repository return error",
			Setup: func(rr *mocks.ReceiverRepository, ss *mocks.TypeService) {
				ss.EXPECT().ValidateConfiguration(receiver.Configurations{
					"token": "key",
				}).Return(nil)
				ss.EXPECT().Encrypt(mock.AnythingOfType("*receiver.Receiver")).Return(nil)
				rr.EXPECT().Update(&receiver.Receiver{
					ID:   123,
					Type: "slack",
					Configurations: map[string]interface{}{
						"token": "key",
					},
				}).Return(errors.New("some error"))
			},
			Rcv: &receiver.Receiver{
				ID:   123,
				Type: "slack",
				Configurations: map[string]interface{}{
					"token": "key",
				},
			},
			Err: errors.New("secureService.repository.Update: some error"),
		},
		{
			Description: "should return nil error if no error returned",
			Setup: func(rr *mocks.ReceiverRepository, ss *mocks.TypeService) {
				ss.EXPECT().ValidateConfiguration(receiver.Configurations{
					"token": "key",
				}).Return(nil)
				ss.EXPECT().Encrypt(mock.AnythingOfType("*receiver.Receiver")).Return(nil)
				rr.EXPECT().Update(&receiver.Receiver{
					ID:   123,
					Type: "slack",
					Configurations: map[string]interface{}{
						"token": "key",
					},
				}).Return(nil)
			},
			Rcv: &receiver.Receiver{
				ID:   123,
				Type: "slack",
				Configurations: map[string]interface{}{
					"token": "key",
				},
			},
			Err: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			var (
				repositoryMock      = new(mocks.ReceiverRepository)
				strategyServiceMock = new(mocks.TypeService)
			)

			registry := map[string]receiver.TypeService{
				receiver.TypeSlack: strategyServiceMock,
			}

			svc := receiver.NewService(repositoryMock, registry)

			tc.Setup(repositoryMock, strategyServiceMock)

			err := svc.UpdateReceiver(tc.Rcv)
			if tc.Err != err {
				if tc.Err.Error() != err.Error() {
					t.Fatalf("got error %s, expected was %s", err.Error(), tc.Err.Error())
				}
			}
			repositoryMock.AssertExpectations(t)
			strategyServiceMock.AssertExpectations(t)
		})
	}
}
