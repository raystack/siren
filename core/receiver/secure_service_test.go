package receiver_test

import (
	"errors"
	testing "testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/odpf/siren/core/receiver"
	"github.com/odpf/siren/core/receiver/mocks"
	mock "github.com/stretchr/testify/mock"
)

func TestSecureService_ListReceivers_Slack(t *testing.T) {
	type testCase struct {
		Description string
		Receivers   []*receiver.Receiver
		Setup       func(*mocks.ReceiverRepository, *mocks.Encryptor, testCase)
		Err         error
	}
	var testCases = []testCase{

		{
			Description: "should return error if List repository error",
			Setup: func(rr *mocks.ReceiverRepository, e *mocks.Encryptor, tc testCase) {
				rr.EXPECT().List().Return(nil, errors.New("some error"))
			},
			Err: errors.New("some error"),
		},
		{
			Description: "should return error if List repository success and no token in configurations field in decrypt error",
			Setup: func(rr *mocks.ReceiverRepository, e *mocks.Encryptor, tc testCase) {
				rr.EXPECT().List().Return([]*receiver.Receiver{
					{
						Id:   10,
						Name: "foo",
						Type: "slack",
						Labels: map[string]string{
							"foo": "bar",
						},
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
				}, nil)
			},
			Err: errors.New("no token field found"),
		},
		{
			Description: "should return error if List repository success and decrypt error",
			Setup: func(rr *mocks.ReceiverRepository, e *mocks.Encryptor, tc testCase) {
				rr.EXPECT().List().Return([]*receiver.Receiver{
					{
						Id:   10,
						Name: "foo",
						Type: "slack",
						Labels: map[string]string{
							"foo": "bar",
						},
						Configurations: map[string]interface{}{
							"token": "key",
						},
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
				}, nil)
				e.EXPECT().Decrypt(mock.AnythingOfType("string")).Return("", errors.New("decrypt error"))
			},
			Err: errors.New("post transform decrypt failed: decrypt error"),
		},
		{
			Description: "should success if list repository and decrypt success",
			Setup: func(rr *mocks.ReceiverRepository, e *mocks.Encryptor, tc testCase) {
				rr.EXPECT().List().Return(tc.Receivers, nil)
				e.EXPECT().Decrypt(mock.AnythingOfType("string")).Return("", nil)
			},
			Receivers: []*receiver.Receiver{
				{
					Id:   10,
					Name: "foo",
					Type: "slack",
					Labels: map[string]string{
						"foo": "bar",
					},
					Configurations: map[string]interface{}{
						"token": "key",
					},
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			},
			Err: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			var (
				repositoryMock = new(mocks.ReceiverRepository)
				encryptorMock  = new(mocks.Encryptor)
			)
			svc := receiver.NewSecureService(encryptorMock, repositoryMock)

			tc.Setup(repositoryMock, encryptorMock, tc)

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
			encryptorMock.AssertExpectations(t)
		})
	}
}

func TestSecureService_CreateReceiver_Slack(t *testing.T) {
	type testCase struct {
		Description string
		Setup       func(*mocks.ReceiverRepository, *mocks.Encryptor, testCase)
		Rcv         *receiver.Receiver
		Err         error
	}
	var testCases = []testCase{
		{
			Description: "should return error if no token in configuration field in encrypt return error",
			Setup:       func(rr *mocks.ReceiverRepository, e *mocks.Encryptor, tc testCase) {},
			Rcv: &receiver.Receiver{
				Id:   123,
				Type: "slack",
			},
			Err: errors.New("no token field found"),
		},
		{
			Description: "should return error if encrypt return error",
			Setup: func(rr *mocks.ReceiverRepository, e *mocks.Encryptor, tc testCase) {
				e.EXPECT().Encrypt(mock.AnythingOfType("string")).Return("", errors.New("some error"))
			},
			Rcv: &receiver.Receiver{
				Id:   123,
				Type: "slack",
				Configurations: map[string]interface{}{
					"token": "key",
				},
			},
			Err: errors.New("pre transform encrypt failed: some error"),
		},
		{
			Description: "should return error if Create repository return error",
			Setup: func(rr *mocks.ReceiverRepository, e *mocks.Encryptor, tc testCase) {
				e.EXPECT().Encrypt(mock.AnythingOfType("string")).Return("", nil)
				rr.EXPECT().Create(tc.Rcv).Return(errors.New("some error"))
			},
			Rcv: &receiver.Receiver{
				Id:   123,
				Type: "slack",
				Configurations: map[string]interface{}{
					"token": "key",
				},
			},
			Err: errors.New("some error"),
		},
		{
			Description: "should return error if decrypt return error",
			Setup: func(rr *mocks.ReceiverRepository, e *mocks.Encryptor, tc testCase) {
				e.EXPECT().Encrypt(mock.AnythingOfType("string")).Return("", nil)
				rr.EXPECT().Create(tc.Rcv).Return(nil)
				e.EXPECT().Decrypt(mock.AnythingOfType("string")).Return("", errors.New("some error"))
			},
			Rcv: &receiver.Receiver{
				Id:   123,
				Type: "slack",
				Configurations: map[string]interface{}{
					"token": "key",
				},
			},
			Err: errors.New("post transform decrypt failed: some error"),
		},
		{
			Description: "should return nil error if no error returned",
			Setup: func(rr *mocks.ReceiverRepository, e *mocks.Encryptor, tc testCase) {
				e.EXPECT().Encrypt(mock.AnythingOfType("string")).Return("", nil)
				rr.EXPECT().Create(tc.Rcv).Return(nil)
				e.EXPECT().Decrypt(mock.AnythingOfType("string")).Return("", nil)
			},
			Rcv: &receiver.Receiver{
				Id:   123,
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
				repositoryMock = new(mocks.ReceiverRepository)
				encryptorMock  = new(mocks.Encryptor)
			)
			svc := receiver.NewSecureService(encryptorMock, repositoryMock)

			tc.Setup(repositoryMock, encryptorMock, tc)

			err := svc.CreateReceiver(tc.Rcv)
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

func TestSecureService_GetReceiver_Slack(t *testing.T) {
	type testCase struct {
		Description string
		Rcv         *receiver.Receiver
		Setup       func(*mocks.ReceiverRepository, *mocks.Encryptor, testCase)
		Err         error
	}
	var (
		testID    = uint64(10)
		testCases = []testCase{
			{
				Description: "should return error if Get repository error",
				Setup: func(rr *mocks.ReceiverRepository, e *mocks.Encryptor, tc testCase) {
					rr.EXPECT().Get(testID).Return(nil, errors.New("some error"))
				},
				Err: errors.New("some error"),
			},
			{
				Description: "should return error if Get repository success and no token field in configurations in decrypt error",
				Setup: func(rr *mocks.ReceiverRepository, e *mocks.Encryptor, tc testCase) {
					rr.EXPECT().Get(testID).Return(&receiver.Receiver{
						Id:   testID,
						Name: "foo",
						Type: "slack",
						Labels: map[string]string{
							"foo": "bar",
						},
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					}, nil)
				},
				Err: errors.New("no token field found"),
			},
			{
				Description: "should return error if Get repository success and decrypt error",
				Setup: func(rr *mocks.ReceiverRepository, e *mocks.Encryptor, tc testCase) {
					rr.EXPECT().Get(testID).Return(&receiver.Receiver{
						Id:   testID,
						Name: "foo",
						Type: "slack",
						Labels: map[string]string{
							"foo": "bar",
						},
						Configurations: map[string]interface{}{
							"token": "key",
						},
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					}, nil)
					e.EXPECT().Decrypt(mock.AnythingOfType("string")).Return("", errors.New("decrypt error"))
				},
				Err: errors.New("post transform decrypt failed: decrypt error"),
			},
			{
				Description: "should success if Get repository and decrypt success",
				Setup: func(rr *mocks.ReceiverRepository, e *mocks.Encryptor, tc testCase) {
					rr.EXPECT().Get(tc.Rcv.Id).Return(tc.Rcv, nil)
					e.EXPECT().Decrypt(mock.AnythingOfType("string")).Return("", nil)
				},
				Rcv: &receiver.Receiver{
					Id:   testID,
					Name: "foo",
					Type: "slack",
					Labels: map[string]string{
						"foo": "bar",
					},
					Configurations: map[string]interface{}{
						"token": "key",
					},
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Err: nil,
			},
		}
	)

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			var (
				repositoryMock = new(mocks.ReceiverRepository)
				encryptorMock  = new(mocks.Encryptor)
			)
			svc := receiver.NewSecureService(encryptorMock, repositoryMock)

			tc.Setup(repositoryMock, encryptorMock, tc)

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
			encryptorMock.AssertExpectations(t)
		})
	}
}

func TestSecureService_UpdateReceiver_Slack(t *testing.T) {
	type testCase struct {
		Description string
		Setup       func(*mocks.ReceiverRepository, *mocks.Encryptor, testCase)
		Rcv         *receiver.Receiver
		Err         error
	}
	var testCases = []testCase{
		{
			Description: "should return error if no token in configurations field in encrypt return error",
			Setup:       func(rr *mocks.ReceiverRepository, e *mocks.Encryptor, tc testCase) {},
			Rcv: &receiver.Receiver{
				Id:   123,
				Type: "slack",
			},
			Err: errors.New("no token field found"),
		},
		{
			Description: "should return error if encrypt return error",
			Setup: func(rr *mocks.ReceiverRepository, e *mocks.Encryptor, tc testCase) {
				e.EXPECT().Encrypt(mock.AnythingOfType("string")).Return("", errors.New("some error"))
			},
			Rcv: &receiver.Receiver{
				Id:   123,
				Type: "slack",
				Configurations: map[string]interface{}{
					"token": "key",
				},
			},
			Err: errors.New("pre transform encrypt failed: some error"),
		},
		{
			Description: "should return error if Update repository return error",
			Setup: func(rr *mocks.ReceiverRepository, e *mocks.Encryptor, tc testCase) {
				e.EXPECT().Encrypt(mock.AnythingOfType("string")).Return("", nil)
				rr.EXPECT().Update(tc.Rcv).Return(errors.New("some error"))
			},
			Rcv: &receiver.Receiver{
				Id:   123,
				Type: "slack",
				Configurations: map[string]interface{}{
					"token": "key",
				},
			},
			Err: errors.New("some error"),
		},
		{
			Description: "should return nil error if no error returned",
			Setup: func(rr *mocks.ReceiverRepository, e *mocks.Encryptor, tc testCase) {
				e.EXPECT().Encrypt(mock.AnythingOfType("string")).Return("", nil)
				rr.EXPECT().Update(tc.Rcv).Return(nil)
			},
			Rcv: &receiver.Receiver{
				Id:   123,
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
				repositoryMock = new(mocks.ReceiverRepository)
				encryptorMock  = new(mocks.Encryptor)
			)
			svc := receiver.NewSecureService(encryptorMock, repositoryMock)

			tc.Setup(repositoryMock, encryptorMock, tc)

			err := svc.UpdateReceiver(tc.Rcv)
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
