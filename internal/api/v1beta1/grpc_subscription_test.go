package v1beta1_test

import (
	"context"
	"testing"
	"time"

	"github.com/odpf/salt/log"
	"github.com/odpf/siren/core/subscription"
	"github.com/odpf/siren/internal/api/v1beta1"
	"github.com/odpf/siren/internal/api/v1beta1/mocks"
	"github.com/odpf/siren/pkg/errors"
	sirenv1beta1 "github.com/odpf/siren/proto/odpf/siren/v1beta1"
	"github.com/stretchr/testify/assert"
)

func TestGRPCServer_ListSubscriptions(t *testing.T) {
	configuration := make(map[string]string)
	configuration["foo"] = "bar"
	match := make(map[string]string)
	match["foo"] = "bar"

	t.Run("should return list of all subscriptions", func(t *testing.T) {
		mockedSubscriptionService := &mocks.SubscriptionService{}

		dummyGRPCServer := v1beta1.NewGRPCServer(nil, log.NewNoop(), nil, nil, nil, nil, nil, nil, mockedSubscriptionService)

		dummyResult := []subscription.Subscription{
			{
				ID:        1,
				URN:       "foo",
				Namespace: 1,
				Receivers: []subscription.Receiver{{ID: 1, Configuration: configuration}},
				Match:     match,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}

		mockedSubscriptionService.EXPECT().List(context.Background(), subscription.Filter{}).Return(dummyResult, nil).Once()
		res, err := dummyGRPCServer.ListSubscriptions(context.Background(), &sirenv1beta1.ListSubscriptionsRequest{})
		assert.Nil(t, err)
		assert.Equal(t, 1, len(res.GetSubscriptions()))
		assert.Equal(t, uint64(1), res.GetSubscriptions()[0].GetId())
		assert.Equal(t, "bar", res.GetSubscriptions()[0].GetMatch()["foo"])
	})

	t.Run("should return error Internal if getting subscriptions fails", func(t *testing.T) {
		mockedSubscriptionService := &mocks.SubscriptionService{}

		dummyGRPCServer := v1beta1.NewGRPCServer(nil, log.NewNoop(), nil, nil, nil, nil, nil, nil, mockedSubscriptionService)

		mockedSubscriptionService.EXPECT().List(context.Background(), subscription.Filter{}).Return(nil, errors.New("random error")).Once()
		res, err := dummyGRPCServer.ListSubscriptions(context.Background(), &sirenv1beta1.ListSubscriptionsRequest{})
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = Internal desc = some unexpected error occurred")
	})
}

func TestGRPCServer_GetSubscription(t *testing.T) {
	configuration := make(map[string]string)
	configuration["foo"] = "bar"
	match := make(map[string]string)
	match["foo"] = "bar"

	t.Run("should return a subscription", func(t *testing.T) {
		mockedSubscriptionService := &mocks.SubscriptionService{}
		dummyGRPCServer := v1beta1.NewGRPCServer(nil, log.NewNoop(), nil, nil, nil, nil, nil, nil, mockedSubscriptionService)
		dummyResult := &subscription.Subscription{
			ID:        1,
			URN:       "foo",
			Namespace: 1,
			Receivers: []subscription.Receiver{{ID: 1, Configuration: configuration}},
			Match:     match,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockedSubscriptionService.EXPECT().Get(context.Background(), uint64(1)).Return(dummyResult, nil).Once()
		res, err := dummyGRPCServer.GetSubscription(context.Background(), &sirenv1beta1.GetSubscriptionRequest{Id: 1})
		assert.Nil(t, err)
		assert.Equal(t, uint64(1), res.GetSubscription().GetId())
		assert.Equal(t, "bar", res.GetSubscription().GetMatch()["foo"])
	})

	t.Run("should return error Not Found if subscriptions not found", func(t *testing.T) {
		mockedSubscriptionService := &mocks.SubscriptionService{}
		dummyGRPCServer := v1beta1.NewGRPCServer(nil, log.NewNoop(), nil, nil, nil, nil, nil, nil, mockedSubscriptionService)
		mockedSubscriptionService.EXPECT().Get(context.Background(), uint64(1)).Return(nil, errors.ErrNotFound).Once()
		res, err := dummyGRPCServer.GetSubscription(context.Background(), &sirenv1beta1.GetSubscriptionRequest{Id: 1})
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = NotFound desc = requested entity not found")
	})

	t.Run("should return error Internal if getting subscription fails", func(t *testing.T) {
		mockedSubscriptionService := &mocks.SubscriptionService{}
		dummyGRPCServer := v1beta1.NewGRPCServer(nil, log.NewNoop(), nil, nil, nil, nil, nil, nil, mockedSubscriptionService)
		mockedSubscriptionService.EXPECT().Get(context.Background(), uint64(1)).
			Return(nil, errors.New("random error")).Once()
		res, err := dummyGRPCServer.GetSubscription(context.Background(), &sirenv1beta1.GetSubscriptionRequest{Id: 1})
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = Internal desc = some unexpected error occurred")
	})
}

func TestGRPCServer_CreateSubscription(t *testing.T) {
	configuration := make(map[string]string)
	configuration["foo"] = "bar"
	match := make(map[string]string)
	match["foo"] = "baz"

	payload := &subscription.Subscription{
		Namespace: 1,
		URN:       "foo",
		Receivers: []subscription.Receiver{{ID: 1, Configuration: configuration}},
		Match:     match,
	}

	dummyResult := &subscription.Subscription{
		ID:        1,
		URN:       "foo",
		Namespace: 10,
		Receivers: []subscription.Receiver{{ID: 1, Configuration: configuration}},
		Match:     match,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	t.Run("should create a subscription", func(t *testing.T) {
		mockedSubscriptionService := &mocks.SubscriptionService{}
		dummyGRPCServer := v1beta1.NewGRPCServer(nil, log.NewNoop(), nil, nil, nil, nil, nil, nil, mockedSubscriptionService)

		mockedSubscriptionService.EXPECT().Create(context.Background(), payload).Run(func(_a0 context.Context, _a1 *subscription.Subscription) {
			_a1.ID = dummyResult.ID
		}).Return(nil).Once()
		res, err := dummyGRPCServer.CreateSubscription(context.Background(), &sirenv1beta1.CreateSubscriptionRequest{
			Namespace: 1,
			Urn:       "foo",
			Receivers: []*sirenv1beta1.ReceiverMetadata{{Id: 1, Configuration: configuration}},
			Match:     match,
		})
		assert.Nil(t, err)
		assert.Equal(t, dummyResult.ID, res.GetId())
	})

	t.Run("should return error InvalidArgument if creating subscriptions return err invalid", func(t *testing.T) {
		mockedSubscriptionService := &mocks.SubscriptionService{}
		dummyGRPCServer := v1beta1.NewGRPCServer(nil, log.NewNoop(), nil, nil, nil, nil, nil, nil, mockedSubscriptionService)

		mockedSubscriptionService.EXPECT().Create(context.Background(), payload).Run(func(_a0 context.Context, _a1 *subscription.Subscription) {
			_a1.ID = dummyResult.ID
		}).Return(errors.ErrInvalid).Once()
		res, err := dummyGRPCServer.CreateSubscription(context.Background(), &sirenv1beta1.CreateSubscriptionRequest{
			Namespace: 1,
			Urn:       "foo",
			Receivers: []*sirenv1beta1.ReceiverMetadata{{Id: 1, Configuration: configuration}},
			Match:     match,
		})
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = request is not valid")
	})

	t.Run("should return error AlreadyExists if creating subscriptions return err conflict", func(t *testing.T) {
		mockedSubscriptionService := &mocks.SubscriptionService{}
		dummyGRPCServer := v1beta1.NewGRPCServer(nil, log.NewNoop(), nil, nil, nil, nil, nil, nil, mockedSubscriptionService)

		mockedSubscriptionService.EXPECT().Create(context.Background(), payload).Run(func(_a0 context.Context, _a1 *subscription.Subscription) {
			_a1.ID = dummyResult.ID
		}).Return(errors.ErrConflict).Once()
		res, err := dummyGRPCServer.CreateSubscription(context.Background(), &sirenv1beta1.CreateSubscriptionRequest{
			Namespace: 1,
			Urn:       "foo",
			Receivers: []*sirenv1beta1.ReceiverMetadata{{Id: 1, Configuration: configuration}},
			Match:     match,
		})
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = AlreadyExists desc = an entity with conflicting identifier exists")
	})

	t.Run("should return error Internal if creating subscriptions fails", func(t *testing.T) {
		mockedSubscriptionService := &mocks.SubscriptionService{}
		dummyGRPCServer := v1beta1.NewGRPCServer(nil, log.NewNoop(), nil, nil, nil, nil, nil, nil, mockedSubscriptionService)

		mockedSubscriptionService.EXPECT().Create(context.Background(), payload).Run(func(_a0 context.Context, _a1 *subscription.Subscription) {
			_a1.ID = dummyResult.ID
		}).Return(errors.New("random error")).Once()
		res, err := dummyGRPCServer.CreateSubscription(context.Background(), &sirenv1beta1.CreateSubscriptionRequest{
			Namespace: 1,
			Urn:       "foo",
			Receivers: []*sirenv1beta1.ReceiverMetadata{{Id: 1, Configuration: configuration}},
			Match:     match,
		})
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = Internal desc = some unexpected error occurred")
	})
}

func TestGRPCServer_UpdateSubscription(t *testing.T) {
	configuration := make(map[string]string)
	configuration["foo"] = "bar"
	match := make(map[string]string)
	match["foo"] = "baz"
	payload := &subscription.Subscription{
		ID:        1,
		Namespace: 10,
		URN:       "foo",
		Receivers: []subscription.Receiver{{ID: 1, Configuration: configuration}},
		Match:     match,
	}

	t.Run("should update a subscription", func(t *testing.T) {
		mockedSubscriptionService := &mocks.SubscriptionService{}
		dummyGRPCServer := v1beta1.NewGRPCServer(nil, log.NewNoop(), nil, nil, nil, nil, nil, nil, mockedSubscriptionService)

		mockedSubscriptionService.EXPECT().Update(context.Background(), payload).Run(func(_a0 context.Context, _a1 *subscription.Subscription) {
			_a1.ID = uint64(1)
		}).Return(nil).Once()
		res, err := dummyGRPCServer.UpdateSubscription(context.Background(), &sirenv1beta1.UpdateSubscriptionRequest{
			Id:        1,
			Namespace: 10,
			Urn:       "foo",
			Receivers: []*sirenv1beta1.ReceiverMetadata{{Id: 1, Configuration: configuration}},
			Match:     match,
		})
		assert.Nil(t, err)
		assert.Equal(t, uint64(1), res.GetId())
	})

	t.Run("should return error Invalid Argument if updating subscriptions return err invalid", func(t *testing.T) {
		mockedSubscriptionService := &mocks.SubscriptionService{}
		dummyGRPCServer := v1beta1.NewGRPCServer(nil, log.NewNoop(), nil, nil, nil, nil, nil, nil, mockedSubscriptionService)
		mockedSubscriptionService.EXPECT().Update(context.Background(), payload).Return(errors.ErrInvalid).Once()
		res, err := dummyGRPCServer.UpdateSubscription(context.Background(), &sirenv1beta1.UpdateSubscriptionRequest{
			Id:        1,
			Namespace: 10,
			Urn:       "foo",
			Receivers: []*sirenv1beta1.ReceiverMetadata{{Id: 1, Configuration: configuration}},
			Match:     match,
		})
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = request is not valid")
	})

	t.Run("should return error AlreadyExist if updating subscriptions return err conflict", func(t *testing.T) {
		mockedSubscriptionService := &mocks.SubscriptionService{}
		dummyGRPCServer := v1beta1.NewGRPCServer(nil, log.NewNoop(), nil, nil, nil, nil, nil, nil, mockedSubscriptionService)
		mockedSubscriptionService.EXPECT().Update(context.Background(), payload).Return(errors.ErrConflict).Once()
		res, err := dummyGRPCServer.UpdateSubscription(context.Background(), &sirenv1beta1.UpdateSubscriptionRequest{
			Id:        1,
			Namespace: 10,
			Urn:       "foo",
			Receivers: []*sirenv1beta1.ReceiverMetadata{{Id: 1, Configuration: configuration}},
			Match:     match,
		})
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = AlreadyExists desc = an entity with conflicting identifier exists")
	})

	t.Run("should return error Internal if updating subscriptions fails", func(t *testing.T) {
		mockedSubscriptionService := &mocks.SubscriptionService{}
		dummyGRPCServer := v1beta1.NewGRPCServer(nil, log.NewNoop(), nil, nil, nil, nil, nil, nil, mockedSubscriptionService)
		mockedSubscriptionService.EXPECT().Update(context.Background(), payload).Return(errors.New("random error")).Once()
		res, err := dummyGRPCServer.UpdateSubscription(context.Background(), &sirenv1beta1.UpdateSubscriptionRequest{
			Id:        1,
			Namespace: 10,
			Urn:       "foo",
			Receivers: []*sirenv1beta1.ReceiverMetadata{{Id: 1, Configuration: configuration}},
			Match:     match,
		})
		assert.Nil(t, res)
		assert.EqualError(t, err, "rpc error: code = Internal desc = some unexpected error occurred")
	})

	t.Run("should return error Invalid for bad requests", func(t *testing.T) {
		mockedSubscriptionService := &mocks.SubscriptionService{}
		dummyGRPCServer := v1beta1.NewGRPCServer(nil, log.NewNoop(), nil, nil, nil, nil, nil, nil, mockedSubscriptionService)
		mockedSubscriptionService.EXPECT().Update(context.Background(), payload).Return(errors.ErrInvalid).Once()
		res, err := dummyGRPCServer.UpdateSubscription(context.Background(), &sirenv1beta1.UpdateSubscriptionRequest{
			Id:        1,
			Namespace: 10,
			Urn:       "foo",
			Receivers: []*sirenv1beta1.ReceiverMetadata{{Id: 1, Configuration: configuration}},
			Match:     match,
		})
		assert.Nil(t, res)
		assert.EqualError(t, err, `rpc error: code = InvalidArgument desc = request is not valid`)
	})
}

func TestGRPCServer_DeleteSubscription(t *testing.T) {
	t.Run("should delete a subscription", func(t *testing.T) {
		mockedSubscriptionService := &mocks.SubscriptionService{}
		dummyGRPCServer := v1beta1.NewGRPCServer(nil, log.NewNoop(), nil, nil, nil, nil, nil, nil, mockedSubscriptionService)

		mockedSubscriptionService.EXPECT().Delete(context.Background(), uint64(1)).Return(nil).Once()
		res, err := dummyGRPCServer.DeleteSubscription(context.Background(), &sirenv1beta1.DeleteSubscriptionRequest{Id: 1})
		assert.Nil(t, err)
		assert.Equal(t, &sirenv1beta1.DeleteSubscriptionResponse{}, res)
	})

	t.Run("should return error Internal if deleting subscription fails", func(t *testing.T) {
		mockedSubscriptionService := &mocks.SubscriptionService{}
		dummyGRPCServer := v1beta1.NewGRPCServer(nil, log.NewNoop(), nil, nil, nil, nil, nil, nil, mockedSubscriptionService)

		mockedSubscriptionService.EXPECT().Delete(context.Background(), uint64(1)).Return(errors.New("random error")).Once()
		res, err := dummyGRPCServer.DeleteSubscription(context.Background(), &sirenv1beta1.DeleteSubscriptionRequest{Id: 1})
		assert.EqualError(t, err, "rpc error: code = Internal desc = some unexpected error occurred")
		assert.Nil(t, res)
	})
}
