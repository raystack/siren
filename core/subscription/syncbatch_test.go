package subscription_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/odpf/siren/core/receiver"
	"github.com/odpf/siren/core/subscription"
	"github.com/odpf/siren/core/subscription/mocks"
	"github.com/odpf/siren/pkg/errors"
	"github.com/stretchr/testify/mock"
)

func TestCreateReceiversMap(t *testing.T) {

	type testCase struct {
		Description          string
		Subscriptions        []subscription.Subscription
		Setup                func(*mocks.ReceiverService)
		ExpectedReceiversMap map[uint64]*receiver.Receiver
		ErrString            string
	}

	var testCases = []testCase{
		{
			Description: "should return error if subscription does not have receivers",
			Setup:       func(rs *mocks.ReceiverService) {},
			ErrString:   "no receivers found in subscription",
		},
		{
			Description: "should return error if receiver service list return error",
			Subscriptions: []subscription.Subscription{
				{
					Receivers: []subscription.Receiver{
						{
							ID: 1,
						},
					},
				},
			},
			Setup: func(rs *mocks.ReceiverService) {
				rs.EXPECT().List(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("receiver.Filter")).Return(nil, errors.New("some error"))
			},
			ErrString: "some error",
		},
		{
			Description: "should return error empty receivers if at least one receiver ids in map don't exist",
			Subscriptions: []subscription.Subscription{
				{
					Receivers: []subscription.Receiver{
						{ID: 1}, {ID: 2}, {ID: 5}, {ID: 7}, {ID: 9}, {ID: 11}, {ID: 12},
					},
				},
			},
			Setup: func(rs *mocks.ReceiverService) {
				rs.EXPECT().List(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("receiver.Filter")).Return([]receiver.Receiver{
					{ID: 1}, {ID: 2}, {ID: 5}, {ID: 7}, {ID: 9}, {ID: 11},
				}, nil)
			},
			ErrString: "receiver id [12] don't exist",
		},
		{
			Description: "should return receivers map if all receiver ids exist",
			Subscriptions: []subscription.Subscription{
				{
					Receivers: []subscription.Receiver{
						{ID: 1}, {ID: 2}, {ID: 5}, {ID: 7}, {ID: 9}, {ID: 11}, {ID: 12},
					},
				},
			},
			ExpectedReceiversMap: map[uint64]*receiver.Receiver{
				1:  {ID: 1},
				2:  {ID: 2},
				5:  {ID: 5},
				7:  {ID: 7},
				9:  {ID: 9},
				11: {ID: 11},
				12: {ID: 12},
			},
			Setup: func(rs *mocks.ReceiverService) {
				rs.EXPECT().List(mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("receiver.Filter")).Return([]receiver.Receiver{
					{ID: 1}, {ID: 2}, {ID: 5}, {ID: 7}, {ID: 9}, {ID: 11}, {ID: 12},
				}, nil)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			var (
				receiverServiceMock = new(mocks.ReceiverService)
			)

			tc.Setup(receiverServiceMock)

			got, err := subscription.CreateReceiversMap(context.Background(), receiverServiceMock, tc.Subscriptions)
			if tc.ErrString != "" {
				if tc.ErrString != err.Error() {
					t.Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}
			if !cmp.Equal(got, tc.ExpectedReceiversMap) {
				t.Fatalf("got result %+v, expected was %+v", got, tc.ExpectedReceiversMap)
			}
		})
	}
}

func TestAssignReceivers(t *testing.T) {

	type testCase struct {
		Description           string
		Subscriptions         []subscription.Subscription
		ReceiversMap          map[uint64]*receiver.Receiver
		Setup                 func(*mocks.ReceiverService)
		ExpectedSubscriptions []subscription.Subscription
		ErrString             string
	}

	var inputSubscriptions = []subscription.Subscription{
		{
			Receivers: []subscription.Receiver{
				{
					ID: 1,
					Configuration: map[string]string{
						"token": "abcabc",
					},
				},
				{
					ID: 2,
					Configuration: map[string]string{
						"token": "abcabc",
					},
				},
			},
		},
		{
			Receivers: []subscription.Receiver{
				{
					ID: 3,
					Configuration: map[string]string{
						"token": "abcabc",
					},
				},
				{
					ID: 4,
					Configuration: map[string]string{
						"token": "abcabc",
					},
				},
			},
		},
	}

	var testCases = []testCase{
		{
			Description:   "should return error if can't found receiver id from map",
			Subscriptions: inputSubscriptions,
			ReceiversMap: map[uint64]*receiver.Receiver{
				1: nil,
				2: nil,
				3: nil,
				4: nil,
			},
			Setup: func(rs *mocks.ReceiverService) {
			},
			ErrString: "receiver id 1 not found",
		},
		{
			Description:   "should return error if get subscription config return error",
			Subscriptions: inputSubscriptions,
			ReceiversMap: map[uint64]*receiver.Receiver{
				1: {ID: 1, Type: receiver.TypeHTTP},
				2: {ID: 2, Type: receiver.TypePagerDuty},
				3: {ID: 3, Type: receiver.TypeSlack},
				4: {ID: 4, Type: receiver.TypeSlack},
			},
			Setup: func(rs *mocks.ReceiverService) {
				rs.EXPECT().EnrichSubscriptionConfig(mock.AnythingOfType("map[string]string"), mock.AnythingOfType("*receiver.Receiver")).Return(nil, errors.New("some error"))
			},
			ErrString: "some error",
		},
		{
			Description:   "should assign receivers to subscription if assigning receivers return no error",
			Subscriptions: inputSubscriptions,
			ReceiversMap: map[uint64]*receiver.Receiver{
				1: {ID: 1, Type: receiver.TypeHTTP},
				2: {ID: 2, Type: receiver.TypePagerDuty},
				3: {ID: 3, Type: receiver.TypeSlack},
				4: {ID: 4, Type: receiver.TypeSlack},
			},
			ExpectedSubscriptions: []subscription.Subscription{
				{
					Receivers: []subscription.Receiver{
						{
							ID:   1,
							Type: receiver.TypeHTTP,
							Configuration: map[string]string{
								"newkey": "newvalue",
								"token":  "abcabc",
							},
						},
						{
							ID:   2,
							Type: receiver.TypePagerDuty,
							Configuration: map[string]string{
								"newkey": "newvalue",
								"token":  "abcabc",
							},
						},
					},
				},
				{
					Receivers: []subscription.Receiver{
						{
							ID:   3,
							Type: receiver.TypeSlack,
							Configuration: map[string]string{
								"newkey": "newvalue",
								"token":  "abcabc",
							},
						},
						{
							ID:   4,
							Type: receiver.TypeSlack,
							Configuration: map[string]string{
								"newkey": "newvalue",
								"token":  "abcabc",
							},
						},
					},
				},
			},
			Setup: func(rs *mocks.ReceiverService) {
				rs.EXPECT().EnrichSubscriptionConfig(mock.AnythingOfType("map[string]string"), mock.AnythingOfType("*receiver.Receiver")).Return(map[string]string{
					"token":  "abcabc",
					"newkey": "newvalue",
				}, nil)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Description, func(t *testing.T) {
			var (
				receiverServiceMock = new(mocks.ReceiverService)
			)

			tc.Setup(receiverServiceMock)

			got, err := subscription.AssignReceivers(receiverServiceMock, tc.ReceiversMap, tc.Subscriptions)
			if tc.ErrString != "" {
				if tc.ErrString != err.Error() {
					t.Fatalf("got error %s, expected was %s", err.Error(), tc.ErrString)
				}
			}
			if !cmp.Equal(got, tc.ExpectedSubscriptions) {
				t.Fatalf("got result %+v, expected was %+v", got, tc.ExpectedSubscriptions)
			}
		})
	}
}
