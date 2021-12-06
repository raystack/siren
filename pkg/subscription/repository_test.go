package subscription

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/mocks"
	"github.com/odpf/siren/pkg/subscription/alertmanager"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"regexp"
	"testing"
	"time"
)

// AnyTime is used to expect arbitrary time value
type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

type RepositoryTestSuite struct {
	suite.Suite
	sqldb      *sql.DB
	dbmock     sqlmock.Sqlmock
	repository SubscriptionRepository
}

func (s *RepositoryTestSuite) SetupTest() {
	db, mock, _ := mocks.NewStore()
	s.sqldb, _ = db.DB()
	s.dbmock = mock
	s.repository = NewRepository(db)
}

func (s *RepositoryTestSuite) TearDownTest() {
	s.sqldb.Close()
}

func (s *RepositoryTestSuite) TestCreate() {
	match := make(StringStringMap)
	inputRandomConfig := make(StringStringMap)
	randomSlackReceiverConfig := make(map[string]interface{})
	randomPagerdutyReceiverConfig := make(map[string]interface{})
	randomHTTPReceiverConfig := make(map[string]interface{})
	match["foo"] = "bar"
	inputRandomConfig["channel_name"] = "test"
	randomSlackReceiverConfig["token"] = "xoxb"
	randomPagerdutyReceiverConfig["service_key"] = "abcd"
	randomHTTPReceiverConfig["url"] = "http://localhost:3000"
	receiver1 := ReceiverMetadata{Id: 1, Configuration: inputRandomConfig}
	receiver2 := ReceiverMetadata{Id: 2, Configuration: make(StringStringMap)}
	expectedSubscription := &Subscription{
		Id:          1,
		NamespaceId: 1,
		Urn:         "foo",
		Match:       match,
		Receiver:    []ReceiverMetadata{receiver2, receiver1},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	dummyNamespace := &domain.Namespace{Id: 1, Provider: 1, Urn: "dummy"}
	dummyProvider := &domain.Provider{Id: 1, Urn: "test", Type: "cortex", Host: "http://localhost:8080"}

	insertQuery := regexp.QuoteMeta(`INSERT INTO "subscriptions" ("namespace_id","urn","receiver","match","created_at","updated_at","id") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "id"`)
	fetchLastInsertedQuery := regexp.QuoteMeta(`SELECT * FROM "subscriptions" WHERE id = 1`)
	fetchSubscriptionsWithinNamespaceQuery := regexp.QuoteMeta(`SELECT * FROM "subscriptions" WHERE namespace_id = 1`)

	s.Run("should create a subscription", func() {
		providerMock := &mocks.ProviderService{}
		namespaceMock := &mocks.NamespaceService{}
		receiverMock := &mocks.ReceiverService{}
		amClientMock := &ClientMock{}
		dummySlackReceivers := &domain.Receiver{Id: 1, Type: "slack", Configurations: randomSlackReceiverConfig}
		dummyPagerdutyReceivers := &domain.Receiver{Id: 2, Type: "pagerduty", Configurations: randomPagerdutyReceiverConfig}
		dummyHTTPReceivers := &domain.Receiver{Id: 3, Type: "http", Configurations: randomHTTPReceiverConfig}
		providerMock.On("GetProvider", uint64(1)).Return(dummyProvider, nil).Once()
		namespaceMock.On("GetNamespace", uint64(1)).Return(dummyNamespace, nil).Once()
		receiverMock.On("ListReceivers").Return(
			[]*domain.Receiver{dummySlackReceivers, dummyPagerdutyReceivers, dummyHTTPReceivers}, nil).Once()
		amClientMock.On("SyncConfig", mock.AnythingOfType("alertmanager.AMConfig"), "dummy").Run(func(args mock.Arguments) {
			rarg := args.Get(0)
			s.Require().IsType(alertmanager.AMConfig{}, rarg)
			r := rarg.(alertmanager.AMConfig)
			s.Equal(3, len(r.Receivers))
			s.Equal("foo_receiverId_1_idx_0", r.Receivers[0].Receiver)
			s.Equal("bar_receiverId_2_idx_0", r.Receivers[1].Receiver)
			s.Equal("baz_receiverId_3_idx_0", r.Receivers[2].Receiver)
		}).Return(nil).Once()

		oldAMClientCreator := alertmanagerClientCreator
		defer func() { alertmanagerClientCreator = oldAMClientCreator }()
		alertmanagerClientCreator = func(c domain.CortexConfig) (alertmanager.Client, error) {
			return amClientMock, nil
		}

		s.dbmock.ExpectBegin()
		s.dbmock.ExpectQuery(insertQuery).WithArgs(expectedSubscription.NamespaceId, expectedSubscription.Urn,
			expectedSubscription.Receiver, expectedSubscription.Match, expectedSubscription.CreatedAt,
			expectedSubscription.UpdatedAt, expectedSubscription.Id).WillReturnRows(sqlmock.NewRows(nil))

		expectedRows := sqlmock.
			NewRows([]string{"urn", "namespace_id", "receiver", "match", "created_at", "updated_at", "id"}).
			AddRow(expectedSubscription.Urn, expectedSubscription.NamespaceId,
				json.RawMessage(`[{"id":1 ,"configuration": {"channel_name": "test"}}, {"id":2 ,"configuration": {}}]`),
				json.RawMessage(`{"foo": "bar"}`), expectedSubscription.CreatedAt, expectedSubscription.UpdatedAt, expectedSubscription.Id)

		expectedRowsInNamespace := sqlmock.
			NewRows([]string{"urn", "namespace_id", "receiver", "match", "created_at", "updated_at", "id"}).
			AddRow(expectedSubscription.Urn, expectedSubscription.NamespaceId,
				json.RawMessage(`[{"id":1 ,"configuration": {"channel_name": "test"}}]`),
				json.RawMessage(`{"foo": "bar"}`), expectedSubscription.CreatedAt, expectedSubscription.UpdatedAt, expectedSubscription.Id).
			AddRow("bar", expectedSubscription.NamespaceId,
				json.RawMessage(`[{"id":2 ,"configuration": {"service_key": "test"}}]`),
				json.RawMessage(`{"foo": "bar"}`), expectedSubscription.CreatedAt, expectedSubscription.UpdatedAt, 2).
			AddRow("baz", expectedSubscription.NamespaceId,
				json.RawMessage(`[{"id":3 ,"configuration": {"url": "http://localhost:3000"}}]`),
				json.RawMessage(`{"foo": "bar"}`), expectedSubscription.CreatedAt, expectedSubscription.UpdatedAt, 3)

		s.dbmock.ExpectQuery(fetchLastInsertedQuery).WillReturnRows(expectedRows)
		s.dbmock.ExpectQuery(fetchSubscriptionsWithinNamespaceQuery).WillReturnRows(expectedRowsInNamespace)
		s.dbmock.ExpectCommit()

		actualSubscription, err := s.repository.Create(expectedSubscription, namespaceMock, providerMock, receiverMock)
		s.Equal(expectedSubscription, actualSubscription)
		s.Nil(err)
	})

	s.Run("should return error in creating subscription", func() {
		providerMock := &mocks.ProviderService{}
		namespaceMock := &mocks.NamespaceService{}
		receiverMock := &mocks.ReceiverService{}

		s.dbmock.ExpectBegin()
		s.dbmock.ExpectQuery(insertQuery).WithArgs(expectedSubscription.NamespaceId, expectedSubscription.Urn,
			expectedSubscription.Receiver, expectedSubscription.Match, expectedSubscription.CreatedAt,
			expectedSubscription.UpdatedAt, expectedSubscription.Id).WillReturnError(errors.New("random error"))
		s.dbmock.ExpectRollback()

		actualSubscription, err := s.repository.Create(expectedSubscription, namespaceMock, providerMock, receiverMock)
		s.EqualError(err, "r.insertSubscriptionIntoDB: failed to insert subscription: random error")
		s.Nil(actualSubscription)
	})

	s.Run("should return error in fetching newly inserted subscription", func() {
		providerMock := &mocks.ProviderService{}
		namespaceMock := &mocks.NamespaceService{}
		receiverMock := &mocks.ReceiverService{}

		s.dbmock.ExpectBegin()
		s.dbmock.ExpectQuery(insertQuery).WithArgs(expectedSubscription.NamespaceId, expectedSubscription.Urn,
			expectedSubscription.Receiver, expectedSubscription.Match, expectedSubscription.CreatedAt,
			expectedSubscription.UpdatedAt, expectedSubscription.Id).WillReturnRows(sqlmock.NewRows(nil))
		s.dbmock.ExpectQuery(fetchLastInsertedQuery).WillReturnError(errors.New("random error"))
		s.dbmock.ExpectRollback()

		actualSubscription, err := s.repository.Create(expectedSubscription, namespaceMock, providerMock, receiverMock)
		s.EqualError(err, "r.insertSubscriptionIntoDB: failed to get newly inserted subscription: random error")
		s.Nil(actualSubscription)
	})

	s.Run("should return error in fetching all subscriptions within given namespace", func() {
		providerMock := &mocks.ProviderService{}
		namespaceMock := &mocks.NamespaceService{}
		receiverMock := &mocks.ReceiverService{}

		s.dbmock.ExpectBegin()
		s.dbmock.ExpectQuery(insertQuery).WithArgs(expectedSubscription.NamespaceId, expectedSubscription.Urn,
			expectedSubscription.Receiver, expectedSubscription.Match, expectedSubscription.CreatedAt,
			expectedSubscription.UpdatedAt, expectedSubscription.Id).WillReturnRows(sqlmock.NewRows(nil))
		expectedRows := sqlmock.
			NewRows([]string{"urn", "namespace_id", "receiver", "match", "created_at", "updated_at", "id"}).
			AddRow(expectedSubscription.Urn, expectedSubscription.NamespaceId,
				json.RawMessage(`[{"id":1 ,"configuration": {"channel_name": "test"}}, {"id":2 ,"configuration": {}}]`),
				json.RawMessage(`{"foo": "bar"}`), expectedSubscription.CreatedAt, expectedSubscription.UpdatedAt, expectedSubscription.Id)
		s.dbmock.ExpectQuery(fetchLastInsertedQuery).WillReturnRows(expectedRows)
		s.dbmock.ExpectQuery(fetchSubscriptionsWithinNamespaceQuery).WillReturnError(errors.New("random error"))
		s.dbmock.ExpectRollback()

		actualSubscription, err := s.repository.Create(expectedSubscription, namespaceMock, providerMock, receiverMock)
		s.EqualError(err, "r.getAllSubscriptionsWithinNamespace: random error")
		s.Nil(actualSubscription)
	})

	s.Run("should return error in fetching namespace details", func() {
		providerMock := &mocks.ProviderService{}
		namespaceMock := &mocks.NamespaceService{}
		receiverMock := &mocks.ReceiverService{}
		providerMock.On("GetProvider", uint64(1)).Return(dummyProvider, nil).Once()
		namespaceMock.On("GetNamespace", uint64(1)).Return(nil, errors.New("random error")).Once()

		s.dbmock.ExpectBegin()
		s.dbmock.ExpectQuery(insertQuery).WithArgs(expectedSubscription.NamespaceId, expectedSubscription.Urn,
			expectedSubscription.Receiver, expectedSubscription.Match, expectedSubscription.CreatedAt,
			expectedSubscription.UpdatedAt, expectedSubscription.Id).WillReturnRows(sqlmock.NewRows(nil))

		expectedRows := sqlmock.
			NewRows([]string{"urn", "namespace_id", "receiver", "match", "created_at", "updated_at", "id"}).
			AddRow(expectedSubscription.Urn, expectedSubscription.NamespaceId,
				json.RawMessage(`[{"id":1 ,"configuration": {"channel_name": "test"}}, {"id":2 ,"configuration": {}}]`),
				json.RawMessage(`{"foo": "bar"}`), expectedSubscription.CreatedAt, expectedSubscription.UpdatedAt, expectedSubscription.Id)

		expectedRowsInNamespace := sqlmock.
			NewRows([]string{"urn", "namespace_id", "receiver", "match", "created_at", "updated_at", "id"}).
			AddRow(expectedSubscription.Urn, expectedSubscription.NamespaceId,
				json.RawMessage(`[{"id":1 ,"configuration": {"channel_name": "test"}}]`),
				json.RawMessage(`{"foo": "bar"}`), expectedSubscription.CreatedAt, expectedSubscription.UpdatedAt, expectedSubscription.Id).
			AddRow("bar", expectedSubscription.NamespaceId,
				json.RawMessage(`[{"id":2 ,"configuration": {"service_key": "test"}}]`),
				json.RawMessage(`{"foo": "bar"}`), expectedSubscription.CreatedAt, expectedSubscription.UpdatedAt, 2).
			AddRow("baz", expectedSubscription.NamespaceId,
				json.RawMessage(`[{"id":3 ,"configuration": {"url": "http://localhost:3000"}}]`),
				json.RawMessage(`{"foo": "bar"}`), expectedSubscription.CreatedAt, expectedSubscription.UpdatedAt, 3)

		s.dbmock.ExpectQuery(fetchLastInsertedQuery).WillReturnRows(expectedRows)
		s.dbmock.ExpectQuery(fetchSubscriptionsWithinNamespaceQuery).WillReturnRows(expectedRowsInNamespace)
		s.dbmock.ExpectRollback()

		actualSubscription, err := s.repository.Create(expectedSubscription, namespaceMock, providerMock, receiverMock)
		s.EqualError(err, "r.getProviderAndNamespaceInfoFromNamespaceId: failed to get namespace details: random error")
		s.Nil(actualSubscription)
	})

	s.Run("should return error in fetching provider details", func() {
		providerMock := &mocks.ProviderService{}
		namespaceMock := &mocks.NamespaceService{}
		receiverMock := &mocks.ReceiverService{}

		providerMock.On("GetProvider", uint64(1)).Return(nil, errors.New("random error")).Once()
		namespaceMock.On("GetNamespace", uint64(1)).Return(dummyNamespace, nil).Once()

		s.dbmock.ExpectBegin()
		s.dbmock.ExpectQuery(insertQuery).WithArgs(expectedSubscription.NamespaceId, expectedSubscription.Urn,
			expectedSubscription.Receiver, expectedSubscription.Match, expectedSubscription.CreatedAt,
			expectedSubscription.UpdatedAt, expectedSubscription.Id).WillReturnRows(sqlmock.NewRows(nil))

		expectedRows := sqlmock.
			NewRows([]string{"urn", "namespace_id", "receiver", "match", "created_at", "updated_at", "id"}).
			AddRow(expectedSubscription.Urn, expectedSubscription.NamespaceId,
				json.RawMessage(`[{"id":1 ,"configuration": {"channel_name": "test"}}, {"id":2 ,"configuration": {}}]`),
				json.RawMessage(`{"foo": "bar"}`), expectedSubscription.CreatedAt, expectedSubscription.UpdatedAt, expectedSubscription.Id)

		expectedRowsInNamespace := sqlmock.
			NewRows([]string{"urn", "namespace_id", "receiver", "match", "created_at", "updated_at", "id"}).
			AddRow(expectedSubscription.Urn, expectedSubscription.NamespaceId,
				json.RawMessage(`[{"id":1 ,"configuration": {"channel_name": "test"}}]`),
				json.RawMessage(`{"foo": "bar"}`), expectedSubscription.CreatedAt, expectedSubscription.UpdatedAt, expectedSubscription.Id).
			AddRow("bar", expectedSubscription.NamespaceId,
				json.RawMessage(`[{"id":2 ,"configuration": {"service_key": "test"}}]`),
				json.RawMessage(`{"foo": "bar"}`), expectedSubscription.CreatedAt, expectedSubscription.UpdatedAt, 2).
			AddRow("baz", expectedSubscription.NamespaceId,
				json.RawMessage(`[{"id":3 ,"configuration": {"url": "http://localhost:3000"}}]`),
				json.RawMessage(`{"foo": "bar"}`), expectedSubscription.CreatedAt, expectedSubscription.UpdatedAt, 3)

		s.dbmock.ExpectQuery(fetchLastInsertedQuery).WillReturnRows(expectedRows)
		s.dbmock.ExpectQuery(fetchSubscriptionsWithinNamespaceQuery).WillReturnRows(expectedRowsInNamespace)
		s.dbmock.ExpectRollback()

		actualSubscription, err := s.repository.Create(expectedSubscription, namespaceMock, providerMock, receiverMock)
		s.EqualError(err, "r.getProviderAndNamespaceInfoFromNamespaceId: failed to get provider details: random error")
		s.Nil(actualSubscription)
	})

	s.Run("should return error in fetching all receivers", func() {
		providerMock := &mocks.ProviderService{}
		namespaceMock := &mocks.NamespaceService{}
		receiverMock := &mocks.ReceiverService{}

		providerMock.On("GetProvider", uint64(1)).Return(dummyProvider, nil).Once()
		namespaceMock.On("GetNamespace", uint64(1)).Return(dummyNamespace, nil).Once()
		receiverMock.On("ListReceivers").Return(nil, errors.New("random error")).Once()

		s.dbmock.ExpectBegin()
		s.dbmock.ExpectQuery(insertQuery).WithArgs(expectedSubscription.NamespaceId, expectedSubscription.Urn,
			expectedSubscription.Receiver, expectedSubscription.Match, expectedSubscription.CreatedAt,
			expectedSubscription.UpdatedAt, expectedSubscription.Id).WillReturnRows(sqlmock.NewRows(nil))

		expectedRows := sqlmock.
			NewRows([]string{"urn", "namespace_id", "receiver", "match", "created_at", "updated_at", "id"}).
			AddRow(expectedSubscription.Urn, expectedSubscription.NamespaceId,
				json.RawMessage(`[{"id":1 ,"configuration": {"channel_name": "test"}}, {"id":2 ,"configuration": {}}]`),
				json.RawMessage(`{"foo": "bar"}`), expectedSubscription.CreatedAt, expectedSubscription.UpdatedAt, expectedSubscription.Id)

		expectedRowsInNamespace := sqlmock.
			NewRows([]string{"urn", "namespace_id", "receiver", "match", "created_at", "updated_at", "id"}).
			AddRow(expectedSubscription.Urn, expectedSubscription.NamespaceId,
				json.RawMessage(`[{"id":1 ,"configuration": {"channel_name": "test"}}]`),
				json.RawMessage(`{"foo": "bar"}`), expectedSubscription.CreatedAt, expectedSubscription.UpdatedAt, expectedSubscription.Id).
			AddRow("bar", expectedSubscription.NamespaceId,
				json.RawMessage(`[{"id":2 ,"configuration": {"service_key": "test"}}]`),
				json.RawMessage(`{"foo": "bar"}`), expectedSubscription.CreatedAt, expectedSubscription.UpdatedAt, 2).
			AddRow("baz", expectedSubscription.NamespaceId,
				json.RawMessage(`[{"id":3 ,"configuration": {"url": "http://localhost:3000"}}]`),
				json.RawMessage(`{"foo": "bar"}`), expectedSubscription.CreatedAt, expectedSubscription.UpdatedAt, 3)

		s.dbmock.ExpectQuery(fetchLastInsertedQuery).WillReturnRows(expectedRows)
		s.dbmock.ExpectQuery(fetchSubscriptionsWithinNamespaceQuery).WillReturnRows(expectedRowsInNamespace)
		s.dbmock.ExpectRollback()

		actualSubscription, err := s.repository.Create(expectedSubscription, namespaceMock, providerMock, receiverMock)
		s.EqualError(err, "r.addReceiversConfiguration: failed to get receivers: random error")
		s.Nil(actualSubscription)
	})

	s.Run("should return error if receiver id not found", func() {
		providerMock := &mocks.ProviderService{}
		namespaceMock := &mocks.NamespaceService{}
		receiverMock := &mocks.ReceiverService{}
		dummySlackReceivers := &domain.Receiver{Id: 10, Type: "slack", Configurations: randomSlackReceiverConfig}

		providerMock.On("GetProvider", uint64(1)).Return(dummyProvider, nil).Once()
		namespaceMock.On("GetNamespace", uint64(1)).Return(dummyNamespace, nil).Once()
		receiverMock.On("ListReceivers").Return([]*domain.Receiver{dummySlackReceivers}, nil).Once()

		s.dbmock.ExpectBegin()
		s.dbmock.ExpectQuery(insertQuery).WithArgs(expectedSubscription.NamespaceId, expectedSubscription.Urn,
			expectedSubscription.Receiver, expectedSubscription.Match, expectedSubscription.CreatedAt,
			expectedSubscription.UpdatedAt, expectedSubscription.Id).WillReturnRows(sqlmock.NewRows(nil))

		expectedRows := sqlmock.
			NewRows([]string{"urn", "namespace_id", "receiver", "match", "created_at", "updated_at", "id"}).
			AddRow(expectedSubscription.Urn, expectedSubscription.NamespaceId,
				json.RawMessage(`[{"id":1 ,"configuration": {"channel_name": "test"}}, {"id":2 ,"configuration": {}}]`),
				json.RawMessage(`{"foo": "bar"}`), expectedSubscription.CreatedAt, expectedSubscription.UpdatedAt, expectedSubscription.Id)

		expectedRowsInNamespace := sqlmock.
			NewRows([]string{"urn", "namespace_id", "receiver", "match", "created_at", "updated_at", "id"}).
			AddRow(expectedSubscription.Urn, expectedSubscription.NamespaceId,
				json.RawMessage(`[{"id":1 ,"configuration": {"channel_name": "test"}}]`),
				json.RawMessage(`{"foo": "bar"}`), expectedSubscription.CreatedAt, expectedSubscription.UpdatedAt, expectedSubscription.Id).
			AddRow("bar", expectedSubscription.NamespaceId,
				json.RawMessage(`[{"id":2 ,"configuration": {"service_key": "test"}}]`),
				json.RawMessage(`{"foo": "bar"}`), expectedSubscription.CreatedAt, expectedSubscription.UpdatedAt, 2).
			AddRow("baz", expectedSubscription.NamespaceId,
				json.RawMessage(`[{"id":3 ,"configuration": {"url": "http://localhost:3000"}}]`),
				json.RawMessage(`{"foo": "bar"}`), expectedSubscription.CreatedAt, expectedSubscription.UpdatedAt, 3)

		s.dbmock.ExpectQuery(fetchLastInsertedQuery).WillReturnRows(expectedRows)
		s.dbmock.ExpectQuery(fetchSubscriptionsWithinNamespaceQuery).WillReturnRows(expectedRowsInNamespace)
		s.dbmock.ExpectRollback()

		actualSubscription, err := s.repository.Create(expectedSubscription, namespaceMock, providerMock, receiverMock)
		s.EqualError(err, "r.addReceiversConfiguration: receiver id 1 does not exist")
		s.Nil(actualSubscription)
	})

	s.Run("should return error if slack channel name not specified in subscription configs", func() {
		providerMock := &mocks.ProviderService{}
		namespaceMock := &mocks.NamespaceService{}
		receiverMock := &mocks.ReceiverService{}
		delete(inputRandomConfig, "channel_name")
		receiver1 = ReceiverMetadata{Id: 1, Configuration: inputRandomConfig}
		expectedSubscription = &Subscription{
			Id:          1,
			NamespaceId: 1,
			Urn:         "foo",
			Match:       match,
			Receiver:    []ReceiverMetadata{receiver1},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		dummySlackReceivers := &domain.Receiver{Id: 1, Type: "slack", Configurations: randomSlackReceiverConfig}
		dummyPagerdutyReceivers := &domain.Receiver{Id: 2, Type: "pagerduty", Configurations: randomPagerdutyReceiverConfig}
		dummyHTTPReceivers := &domain.Receiver{Id: 3, Type: "http", Configurations: randomHTTPReceiverConfig}
		providerMock.On("GetProvider", uint64(1)).Return(dummyProvider, nil).Once()
		namespaceMock.On("GetNamespace", uint64(1)).Return(dummyNamespace, nil).Once()
		receiverMock.On("ListReceivers").Return(
			[]*domain.Receiver{dummySlackReceivers, dummyPagerdutyReceivers, dummyHTTPReceivers}, nil).Once()

		s.dbmock.ExpectBegin()
		s.dbmock.ExpectQuery(insertQuery).WithArgs(expectedSubscription.NamespaceId, expectedSubscription.Urn,
			expectedSubscription.Receiver, expectedSubscription.Match, expectedSubscription.CreatedAt,
			expectedSubscription.UpdatedAt, expectedSubscription.Id).WillReturnRows(sqlmock.NewRows(nil))

		expectedRows := sqlmock.
			NewRows([]string{"urn", "namespace_id", "receiver", "match", "created_at", "updated_at", "id"}).
			AddRow(expectedSubscription.Urn, expectedSubscription.NamespaceId,
				json.RawMessage(`[{"id":1 ,"configuration": {"channel_name": "test"}}, {"id":2 ,"configuration": {}}]`),
				json.RawMessage(`{"foo": "bar"}`), expectedSubscription.CreatedAt, expectedSubscription.UpdatedAt, expectedSubscription.Id)

		expectedRowsInNamespace := sqlmock.
			NewRows([]string{"urn", "namespace_id", "receiver", "match", "created_at", "updated_at", "id"}).
			AddRow(expectedSubscription.Urn, expectedSubscription.NamespaceId,
				json.RawMessage(`[{"id":1 ,"configuration": {}}]`),
				json.RawMessage(`{"foo": "bar"}`), expectedSubscription.CreatedAt, expectedSubscription.UpdatedAt, expectedSubscription.Id).
			AddRow("bar", expectedSubscription.NamespaceId,
				json.RawMessage(`[{"id":2 ,"configuration": {"service_key": "test"}}]`),
				json.RawMessage(`{"foo": "bar"}`), expectedSubscription.CreatedAt, expectedSubscription.UpdatedAt, 2).
			AddRow("baz", expectedSubscription.NamespaceId,
				json.RawMessage(`[{"id":3 ,"configuration": {"url": "http://localhost:3000"}}]`),
				json.RawMessage(`{"foo": "bar"}`), expectedSubscription.CreatedAt, expectedSubscription.UpdatedAt, 3)

		s.dbmock.ExpectQuery(fetchLastInsertedQuery).WillReturnRows(expectedRows)
		s.dbmock.ExpectQuery(fetchSubscriptionsWithinNamespaceQuery).WillReturnRows(expectedRowsInNamespace)
		s.dbmock.ExpectRollback()

		actualSubscription, err := s.repository.Create(expectedSubscription, namespaceMock, providerMock, receiverMock)
		s.EqualError(err, "r.addReceiversConfiguration: configuration.channel_name missing from receiver with id 1")
		s.Nil(actualSubscription)
		inputRandomConfig["channel_name"] = "test"
	})

	s.Run("should return error for unsupported receiver type", func() {
		providerMock := &mocks.ProviderService{}
		namespaceMock := &mocks.NamespaceService{}
		receiverMock := &mocks.ReceiverService{}

		unsupportedReceiver := &domain.Receiver{Id: 1, Type: "email", Configurations: make(map[string]interface{})}
		providerMock.On("GetProvider", uint64(1)).Return(dummyProvider, nil).Once()
		namespaceMock.On("GetNamespace", uint64(1)).Return(dummyNamespace, nil).Once()
		receiverMock.On("ListReceivers").Return([]*domain.Receiver{unsupportedReceiver}, nil).Once()

		s.dbmock.ExpectBegin()
		s.dbmock.ExpectQuery(insertQuery).WithArgs(expectedSubscription.NamespaceId, expectedSubscription.Urn,
			expectedSubscription.Receiver, expectedSubscription.Match, expectedSubscription.CreatedAt,
			expectedSubscription.UpdatedAt, expectedSubscription.Id).WillReturnRows(sqlmock.NewRows(nil))

		expectedRows := sqlmock.
			NewRows([]string{"urn", "namespace_id", "receiver", "match", "created_at", "updated_at", "id"}).
			AddRow(expectedSubscription.Urn, expectedSubscription.NamespaceId,
				json.RawMessage(`[{"id":1 ,"configuration": {"channel_name": "test"}}, {"id":2 ,"configuration": {}}]`),
				json.RawMessage(`{"foo": "bar"}`), expectedSubscription.CreatedAt, expectedSubscription.UpdatedAt, expectedSubscription.Id)

		expectedRowsInNamespace := sqlmock.
			NewRows([]string{"urn", "namespace_id", "receiver", "match", "created_at", "updated_at", "id"}).
			AddRow(expectedSubscription.Urn, expectedSubscription.NamespaceId,
				json.RawMessage(`[{"id":1 ,"configuration": {"channel_name": "test"}}]`),
				json.RawMessage(`{"foo": "bar"}`), expectedSubscription.CreatedAt, expectedSubscription.UpdatedAt, expectedSubscription.Id).
			AddRow("bar", expectedSubscription.NamespaceId,
				json.RawMessage(`[{"id":2 ,"configuration": {"service_key": "test"}}]`),
				json.RawMessage(`{"foo": "bar"}`), expectedSubscription.CreatedAt, expectedSubscription.UpdatedAt, 2).
			AddRow("baz", expectedSubscription.NamespaceId,
				json.RawMessage(`[{"id":3 ,"configuration": {"url": "http://localhost:3000"}}]`),
				json.RawMessage(`{"foo": "bar"}`), expectedSubscription.CreatedAt, expectedSubscription.UpdatedAt, 3)

		s.dbmock.ExpectQuery(fetchLastInsertedQuery).WillReturnRows(expectedRows)
		s.dbmock.ExpectQuery(fetchSubscriptionsWithinNamespaceQuery).WillReturnRows(expectedRowsInNamespace)
		s.dbmock.ExpectRollback()

		actualSubscription, err := s.repository.Create(expectedSubscription, namespaceMock, providerMock, receiverMock)
		s.EqualError(err, "r.addReceiversConfiguration: subscriptions for receiver type email not supported via Siren inside Cortex")
		s.Nil(actualSubscription)
	})

	s.Run("should return error in alertmanager client initialization", func() {
		providerMock := &mocks.ProviderService{}
		namespaceMock := &mocks.NamespaceService{}
		receiverMock := &mocks.ReceiverService{}
		amClientMock := &ClientMock{}
		dummySlackReceivers := &domain.Receiver{Id: 1, Type: "slack", Configurations: randomSlackReceiverConfig}
		dummyPagerdutyReceivers := &domain.Receiver{Id: 2, Type: "pagerduty", Configurations: randomPagerdutyReceiverConfig}
		dummyHTTPReceivers := &domain.Receiver{Id: 3, Type: "http", Configurations: randomHTTPReceiverConfig}
		providerMock.On("GetProvider", uint64(1)).Return(dummyProvider, nil).Once()
		namespaceMock.On("GetNamespace", uint64(1)).Return(dummyNamespace, nil).Once()
		receiverMock.On("ListReceivers").Return(
			[]*domain.Receiver{dummySlackReceivers, dummyPagerdutyReceivers, dummyHTTPReceivers}, nil).Once()
		amClientMock.On("SyncConfig", mock.AnythingOfType("alertmanager.AMConfig"), "dummy").Run(func(args mock.Arguments) {
			rarg := args.Get(0)
			s.Require().IsType(alertmanager.AMConfig{}, rarg)
			r := rarg.(alertmanager.AMConfig)
			s.Equal(3, len(r.Receivers))
			s.Equal("foo_receiverId_1_idx_0", r.Receivers[0].Receiver)
			s.Equal("bar_receiverId_2_idx_0", r.Receivers[1].Receiver)
			s.Equal("baz_receiverId_3_idx_0", r.Receivers[2].Receiver)
		}).Return(nil).Once()

		oldAMClientCreator := alertmanagerClientCreator
		defer func() { alertmanagerClientCreator = oldAMClientCreator }()
		alertmanagerClientCreator = func(c domain.CortexConfig) (alertmanager.Client, error) {
			return nil, errors.New("random error")
		}

		s.dbmock.ExpectBegin()
		s.dbmock.ExpectQuery(insertQuery).WithArgs(expectedSubscription.NamespaceId, expectedSubscription.Urn,
			expectedSubscription.Receiver, expectedSubscription.Match, expectedSubscription.CreatedAt,
			expectedSubscription.UpdatedAt, expectedSubscription.Id).WillReturnRows(sqlmock.NewRows(nil))

		expectedRows := sqlmock.
			NewRows([]string{"urn", "namespace_id", "receiver", "match", "created_at", "updated_at", "id"}).
			AddRow(expectedSubscription.Urn, expectedSubscription.NamespaceId,
				json.RawMessage(`[{"id":1 ,"configuration": {"channel_name": "test"}}, {"id":2 ,"configuration": {}}]`),
				json.RawMessage(`{"foo": "bar"}`), expectedSubscription.CreatedAt, expectedSubscription.UpdatedAt, expectedSubscription.Id)

		expectedRowsInNamespace := sqlmock.
			NewRows([]string{"urn", "namespace_id", "receiver", "match", "created_at", "updated_at", "id"}).
			AddRow(expectedSubscription.Urn, expectedSubscription.NamespaceId,
				json.RawMessage(`[{"id":1 ,"configuration": {"channel_name": "test"}}]`),
				json.RawMessage(`{"foo": "bar"}`), expectedSubscription.CreatedAt, expectedSubscription.UpdatedAt, expectedSubscription.Id).
			AddRow("bar", expectedSubscription.NamespaceId,
				json.RawMessage(`[{"id":2 ,"configuration": {"service_key": "test"}}]`),
				json.RawMessage(`{"foo": "bar"}`), expectedSubscription.CreatedAt, expectedSubscription.UpdatedAt, 2).
			AddRow("baz", expectedSubscription.NamespaceId,
				json.RawMessage(`[{"id":3 ,"configuration": {"url": "http://localhost:3000"}}]`),
				json.RawMessage(`{"foo": "bar"}`), expectedSubscription.CreatedAt, expectedSubscription.UpdatedAt, 3)

		s.dbmock.ExpectQuery(fetchLastInsertedQuery).WillReturnRows(expectedRows)
		s.dbmock.ExpectQuery(fetchSubscriptionsWithinNamespaceQuery).WillReturnRows(expectedRowsInNamespace)
		s.dbmock.ExpectRollback()

		actualSubscription, err := s.repository.Create(expectedSubscription, namespaceMock, providerMock, receiverMock)
		s.EqualError(err, "alertmanagerClientCreator: : random error")
		s.Nil(actualSubscription)
	})

	s.Run("should return error syncing config with alertmanager", func() {
		providerMock := &mocks.ProviderService{}
		namespaceMock := &mocks.NamespaceService{}
		receiverMock := &mocks.ReceiverService{}
		amClientMock := &ClientMock{}

		dummySlackReceivers := &domain.Receiver{Id: 1, Type: "slack", Configurations: randomSlackReceiverConfig}
		dummyPagerdutyReceivers := &domain.Receiver{Id: 2, Type: "pagerduty", Configurations: randomPagerdutyReceiverConfig}
		dummyHTTPReceivers := &domain.Receiver{Id: 3, Type: "http", Configurations: randomHTTPReceiverConfig}

		providerMock.On("GetProvider", uint64(1)).Return(dummyProvider, nil).Once()
		namespaceMock.On("GetNamespace", uint64(1)).Return(dummyNamespace, nil).Once()
		receiverMock.On("ListReceivers").Return(
			[]*domain.Receiver{dummySlackReceivers, dummyPagerdutyReceivers, dummyHTTPReceivers}, nil).Once()
		amClientMock.On("SyncConfig", mock.AnythingOfType("alertmanager.AMConfig"), "dummy").
			Return(errors.New("random error")).Once()

		oldAMClientCreator := alertmanagerClientCreator
		defer func() { alertmanagerClientCreator = oldAMClientCreator }()
		alertmanagerClientCreator = func(c domain.CortexConfig) (alertmanager.Client, error) {
			return amClientMock, nil
		}

		s.dbmock.ExpectBegin()
		s.dbmock.ExpectQuery(insertQuery).WithArgs(expectedSubscription.NamespaceId, expectedSubscription.Urn,
			expectedSubscription.Receiver, expectedSubscription.Match, expectedSubscription.CreatedAt,
			expectedSubscription.UpdatedAt, expectedSubscription.Id).WillReturnRows(sqlmock.NewRows(nil))

		expectedRows := sqlmock.
			NewRows([]string{"urn", "namespace_id", "receiver", "match", "created_at", "updated_at", "id"}).
			AddRow(expectedSubscription.Urn, expectedSubscription.NamespaceId,
				json.RawMessage(`[{"id":1 ,"configuration": {"channel_name": "test"}}, {"id":2 ,"configuration": {}}]`),
				json.RawMessage(`{"foo": "bar"}`), expectedSubscription.CreatedAt, expectedSubscription.UpdatedAt, expectedSubscription.Id)

		expectedRowsInNamespace := sqlmock.
			NewRows([]string{"urn", "namespace_id", "receiver", "match", "created_at", "updated_at", "id"}).
			AddRow(expectedSubscription.Urn, expectedSubscription.NamespaceId,
				json.RawMessage(`[{"id":1 ,"configuration": {"channel_name": "test"}}]`),
				json.RawMessage(`{"foo": "bar"}`), expectedSubscription.CreatedAt, expectedSubscription.UpdatedAt, expectedSubscription.Id).
			AddRow("bar", expectedSubscription.NamespaceId,
				json.RawMessage(`[{"id":2 ,"configuration": {"service_key": "test"}}]`),
				json.RawMessage(`{"foo": "bar"}`), expectedSubscription.CreatedAt, expectedSubscription.UpdatedAt, 2).
			AddRow("baz", expectedSubscription.NamespaceId,
				json.RawMessage(`[{"id":3 ,"configuration": {"url": "http://localhost:3000"}}]`),
				json.RawMessage(`{"foo": "bar"}`), expectedSubscription.CreatedAt, expectedSubscription.UpdatedAt, 3)

		s.dbmock.ExpectQuery(fetchLastInsertedQuery).WillReturnRows(expectedRows)
		s.dbmock.ExpectQuery(fetchSubscriptionsWithinNamespaceQuery).WillReturnRows(expectedRowsInNamespace)
		s.dbmock.ExpectRollback()

		actualSubscription, err := s.repository.Create(expectedSubscription, namespaceMock, providerMock, receiverMock)
		s.EqualError(err, "r.amClient.SyncConfig: random error")
		s.Nil(actualSubscription)
	})

	s.Run("should return error for unsupported providers", func() {
		providerMock := &mocks.ProviderService{}
		namespaceMock := &mocks.NamespaceService{}
		receiverMock := &mocks.ReceiverService{}
		dummySlackReceivers := &domain.Receiver{Id: 1, Type: "slack", Configurations: randomSlackReceiverConfig}
		dummyPagerdutyReceivers := &domain.Receiver{Id: 2, Type: "pagerduty", Configurations: randomPagerdutyReceiverConfig}
		dummyHTTPReceivers := &domain.Receiver{Id: 3, Type: "http", Configurations: randomHTTPReceiverConfig}
		providerMock.On("GetProvider", uint64(1)).Return(&domain.Provider{Id: 1, Urn: "test", Type: "prometheus"}, nil).Once()
		namespaceMock.On("GetNamespace", uint64(1)).Return(dummyNamespace, nil).Once()
		receiverMock.On("ListReceivers").Return(
			[]*domain.Receiver{dummySlackReceivers, dummyPagerdutyReceivers, dummyHTTPReceivers}, nil).Once()

		s.dbmock.ExpectBegin()
		s.dbmock.ExpectQuery(insertQuery).WithArgs(expectedSubscription.NamespaceId, expectedSubscription.Urn,
			expectedSubscription.Receiver, expectedSubscription.Match, expectedSubscription.CreatedAt,
			expectedSubscription.UpdatedAt, expectedSubscription.Id).WillReturnRows(sqlmock.NewRows(nil))

		expectedRows := sqlmock.
			NewRows([]string{"urn", "namespace_id", "receiver", "match", "created_at", "updated_at", "id"}).
			AddRow(expectedSubscription.Urn, expectedSubscription.NamespaceId,
				json.RawMessage(`[{"id":1 ,"configuration": {"channel_name": "test"}}, {"id":2 ,"configuration": {}}]`),
				json.RawMessage(`{"foo": "bar"}`), expectedSubscription.CreatedAt, expectedSubscription.UpdatedAt, expectedSubscription.Id)

		expectedRowsInNamespace := sqlmock.
			NewRows([]string{"urn", "namespace_id", "receiver", "match", "created_at", "updated_at", "id"}).
			AddRow(expectedSubscription.Urn, expectedSubscription.NamespaceId,
				json.RawMessage(`[{"id":1 ,"configuration": {"channel_name": "test"}}]`),
				json.RawMessage(`{"foo": "bar"}`), expectedSubscription.CreatedAt, expectedSubscription.UpdatedAt, expectedSubscription.Id).
			AddRow("bar", expectedSubscription.NamespaceId,
				json.RawMessage(`[{"id":2 ,"configuration": {"service_key": "test"}}]`),
				json.RawMessage(`{"foo": "bar"}`), expectedSubscription.CreatedAt, expectedSubscription.UpdatedAt, 2).
			AddRow("baz", expectedSubscription.NamespaceId,
				json.RawMessage(`[{"id":3 ,"configuration": {"url": "http://localhost:3000"}}]`),
				json.RawMessage(`{"foo": "bar"}`), expectedSubscription.CreatedAt, expectedSubscription.UpdatedAt, 3)

		s.dbmock.ExpectQuery(fetchLastInsertedQuery).WillReturnRows(expectedRows)
		s.dbmock.ExpectQuery(fetchSubscriptionsWithinNamespaceQuery).WillReturnRows(expectedRowsInNamespace)
		s.dbmock.ExpectRollback()

		actualSubscription, err := s.repository.Create(expectedSubscription, namespaceMock, providerMock, receiverMock)
		s.EqualError(err, "subscriptions for provider type 'prometheus' not supported")
		s.Nil(actualSubscription)
	})
}

func (s *RepositoryTestSuite) TestGet() {
	expectedSubscription := &Subscription{
		Id:          1,
		NamespaceId: 1,
		Urn:         "foo",
		Match:       make(StringStringMap),
		Receiver:    []ReceiverMetadata{{Id: 1, Configuration: make(StringStringMap)}},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	s.Run("should get subscription by id", func() {
		selectQuery := regexp.QuoteMeta(`SELECT * FROM "subscriptions" WHERE id = 1`)
		expectedRows := sqlmock.
			NewRows([]string{"urn", "namespace_id", "receiver", "match", "created_at", "updated_at", "id"}).
			AddRow(expectedSubscription.Urn, expectedSubscription.NamespaceId,
				json.RawMessage(`[{"id":1 ,"configuration": {}}]`), json.RawMessage(`{}`),
				expectedSubscription.CreatedAt, expectedSubscription.UpdatedAt, expectedSubscription.Id)
		s.dbmock.ExpectQuery(selectQuery).WillReturnRows(expectedRows)
		actualSubscription, err := s.repository.Get(1)
		s.Equal(expectedSubscription.Id, actualSubscription.Id)
		s.Equal(expectedSubscription.NamespaceId, actualSubscription.NamespaceId)
		s.Equal(expectedSubscription.Urn, actualSubscription.Urn)
		s.Equal(expectedSubscription.Match, actualSubscription.Match)
		s.Nil(err)
	})

	s.Run("should return error in finding subscription", func() {
		selectQuery := regexp.QuoteMeta(`SELECT * FROM "subscriptions" WHERE id = 1`)
		s.dbmock.ExpectQuery(selectQuery).WillReturnError(errors.New("random error"))
		actualSubscription, err := s.repository.Get(1)
		s.EqualError(err, "random error")
		s.Nil(actualSubscription)
	})

	s.Run("should return nil if subscription not found", func() {
		selectQuery := regexp.QuoteMeta(`SELECT * FROM "subscriptions" WHERE id = 1`)
		expectedRows := sqlmock.
			NewRows([]string{"urn", "namespace_id", "receiver", "match", "created_at", "updated_at", "id"})
		s.dbmock.ExpectQuery(selectQuery).WillReturnRows(expectedRows)
		actualSubscription, err := s.repository.Get(1)
		s.Nil(actualSubscription)
		s.Nil(err)
	})

}

func (s *RepositoryTestSuite) TestList() {
	expectedSubscription := &Subscription{
		Id:          1,
		NamespaceId: 1,
		Urn:         "foo",
		Match:       make(StringStringMap),
		Receiver:    []ReceiverMetadata{{Id: 1, Configuration: make(StringStringMap)}},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	s.Run("should get all subscriptions", func() {
		selectQuery := regexp.QuoteMeta(`select * from subscriptions`)
		expectedRows := sqlmock.
			NewRows([]string{"urn", "namespace_id", "receiver", "match", "created_at", "updated_at", "id"}).
			AddRow(expectedSubscription.Urn, expectedSubscription.NamespaceId,
				json.RawMessage(`[{"id":1 ,"configuration": {}}]`), json.RawMessage(`{}`),
				expectedSubscription.CreatedAt, expectedSubscription.UpdatedAt, expectedSubscription.Id)
		s.dbmock.ExpectQuery(selectQuery).WillReturnRows(expectedRows)
		actualSubscriptions, err := s.repository.List()

		s.Equal(1, len(actualSubscriptions))
		s.Equal(expectedSubscription.Id, actualSubscriptions[0].Id)
		s.Equal(expectedSubscription.NamespaceId, actualSubscriptions[0].NamespaceId)
		s.Equal(expectedSubscription.Urn, actualSubscriptions[0].Urn)
		s.Equal(expectedSubscription.Match, actualSubscriptions[0].Match)
		s.Nil(err)
	})

	s.Run("should return error in fetching subscriptions", func() {
		selectQuery := regexp.QuoteMeta(`select * from subscriptions`)
		s.dbmock.ExpectQuery(selectQuery).WillReturnError(errors.New("random error"))
		actualSubscriptions, err := s.repository.List()

		s.EqualError(err, "random error")
		s.Nil(actualSubscriptions)
	})
}

func (s *RepositoryTestSuite) TestUpdate() {
	match := make(StringStringMap)
	inputRandomConfig := make(StringStringMap)
	randomSlackReceiverConfig := make(map[string]interface{})
	randomPagerdutyReceiverConfig := make(map[string]interface{})
	randomHTTPReceiverConfig := make(map[string]interface{})
	match["foo"] = "bar"
	inputRandomConfig["channel_name"] = "test"
	randomSlackReceiverConfig["token"] = "xoxb"
	randomPagerdutyReceiverConfig["service_key"] = "abcd"
	randomHTTPReceiverConfig["url"] = "http://localhost:3000"
	receiver := ReceiverMetadata{Id: 1, Configuration: inputRandomConfig}
	subscription := &Subscription{
		Id:          1,
		NamespaceId: 1,
		Urn:         "foo",
		Match:       match,
		Receiver:    []ReceiverMetadata{receiver},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	inputRandomConfig["channel_name"] = "updated_channel"

	input := &Subscription{
		Id:          1,
		NamespaceId: 1,
		Urn:         "foo",
		Match:       match,
		Receiver:    []ReceiverMetadata{receiver},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	dummyNamespace := &domain.Namespace{Id: 1, Provider: 1, Urn: "dummy"}
	dummyProvider := &domain.Provider{Id: 1, Urn: "test", Type: "cortex", Host: "http://localhost:8080"}

	firstSelectQuery := regexp.QuoteMeta(`SELECT * FROM "subscriptions" WHERE id = 1`)
	secondSelectQuery := regexp.QuoteMeta(`SELECT * FROM "subscriptions" WHERE id = 1`)
	updateQuery := regexp.QuoteMeta(`UPDATE "subscriptions" SET "namespace_id"=$1,"urn"=$2,"receiver"=$3,"match"=$4,"created_at"=$5,"updated_at"=$6 WHERE id = $7 AND "id" = $8`)
	fetchSubscriptionsWithinNamespaceQuery := regexp.QuoteMeta(`SELECT * FROM "subscriptions" WHERE namespace_id = 1`)

	s.Run("should update a subscription", func() {
		providerMock := &mocks.ProviderService{}
		namespaceMock := &mocks.NamespaceService{}
		receiverMock := &mocks.ReceiverService{}
		amClientMock := &ClientMock{}
		dummySlackReceivers := &domain.Receiver{Id: 1, Type: "slack", Configurations: randomSlackReceiverConfig}
		dummyPagerdutyReceivers := &domain.Receiver{Id: 2, Type: "pagerduty", Configurations: randomPagerdutyReceiverConfig}
		dummyHTTPReceivers := &domain.Receiver{Id: 3, Type: "http", Configurations: randomHTTPReceiverConfig}
		providerMock.On("GetProvider", uint64(1)).Return(dummyProvider, nil).Once()
		namespaceMock.On("GetNamespace", uint64(1)).Return(dummyNamespace, nil).Once()
		receiverMock.On("ListReceivers").Return(
			[]*domain.Receiver{dummySlackReceivers, dummyPagerdutyReceivers, dummyHTTPReceivers}, nil).Once()
		amClientMock.On("SyncConfig", mock.AnythingOfType("alertmanager.AMConfig"), "dummy").Run(func(args mock.Arguments) {
			rarg := args.Get(0)
			s.Require().IsType(alertmanager.AMConfig{}, rarg)
			r := rarg.(alertmanager.AMConfig)
			s.Equal(3, len(r.Receivers))
			s.Equal("foo_receiverId_1_idx_0", r.Receivers[0].Receiver)
			s.Equal("updated_channel", r.Receivers[0].Configuration["channel_name"])
			s.Equal("bar_receiverId_2_idx_0", r.Receivers[1].Receiver)
			s.Equal("baz_receiverId_3_idx_0", r.Receivers[2].Receiver)
		}).Return(nil).Once()

		oldAMClientCreator := alertmanagerClientCreator
		defer func() { alertmanagerClientCreator = oldAMClientCreator }()
		alertmanagerClientCreator = func(c domain.CortexConfig) (alertmanager.Client, error) {
			return amClientMock, nil
		}

		s.dbmock.ExpectBegin()
		expectedRowsBeforeUpdate := sqlmock.
			NewRows([]string{"urn", "namespace_id", "receiver", "match", "created_at", "updated_at", "id"}).
			AddRow(subscription.Urn, subscription.NamespaceId,
				json.RawMessage(`[{"id":1 ,"configuration": {"channel_name":"test"}}]`), json.RawMessage(`{"foo":"bar"}`),
				subscription.CreatedAt, subscription.UpdatedAt, subscription.Id)
		s.dbmock.ExpectQuery(firstSelectQuery).WillReturnRows(expectedRowsBeforeUpdate)
		s.dbmock.ExpectExec(updateQuery).WithArgs(subscription.NamespaceId, subscription.Urn,
			subscription.Receiver, subscription.Match, AnyTime{}, AnyTime{}, subscription.Id, subscription.Id).
			WillReturnResult(sqlmock.NewResult(1, 1))

		expectedRowsAfterUpdate := sqlmock.
			NewRows([]string{"urn", "namespace_id", "receiver", "match", "created_at", "updated_at", "id"}).
			AddRow(subscription.Urn, subscription.NamespaceId,
				json.RawMessage(`[{"id":1 ,"configuration": {"channel_name":"updated_channel"}}]`), json.RawMessage(`{"foo":"bar"}`),
				subscription.CreatedAt, subscription.UpdatedAt, subscription.Id)

		s.dbmock.ExpectQuery(secondSelectQuery).WillReturnRows(expectedRowsAfterUpdate)

		expectedRowsInNamespace := sqlmock.
			NewRows([]string{"urn", "namespace_id", "receiver", "match", "created_at", "updated_at", "id"}).
			AddRow(subscription.Urn, subscription.NamespaceId,
				json.RawMessage(`[{"id":1 ,"configuration": {"channel_name": "updated_channel"}}]`),
				json.RawMessage(`{"foo": "bar"}`), subscription.CreatedAt, subscription.UpdatedAt, subscription.Id).
			AddRow("bar", subscription.NamespaceId,
				json.RawMessage(`[{"id":2 ,"configuration": {"service_key": "test"}}]`),
				json.RawMessage(`{"foo": "bar"}`), subscription.CreatedAt, subscription.UpdatedAt, 2).
			AddRow("baz", subscription.NamespaceId,
				json.RawMessage(`[{"id":3 ,"configuration": {"url": "http://localhost:3000"}}]`),
				json.RawMessage(`{"foo": "bar"}`), subscription.CreatedAt, subscription.UpdatedAt, 3)

		s.dbmock.ExpectQuery(fetchSubscriptionsWithinNamespaceQuery).WillReturnRows(expectedRowsInNamespace)
		s.dbmock.ExpectCommit()

		actualSubscription, err := s.repository.Update(input, namespaceMock, providerMock, receiverMock)
		s.Equal(subscription, actualSubscription)
		s.Nil(err)
	})

	s.Run("should return error in fetching subscription before update", func() {
		providerMock := &mocks.ProviderService{}
		namespaceMock := &mocks.NamespaceService{}
		receiverMock := &mocks.ReceiverService{}
		s.dbmock.ExpectBegin()
		s.dbmock.ExpectQuery(firstSelectQuery).WillReturnError(errors.New("random error"))
		s.dbmock.ExpectRollback()

		actualSubscription, err := s.repository.Update(input, namespaceMock, providerMock, receiverMock)
		s.EqualError(err, "random error")
		s.Nil(actualSubscription)
	})

	s.Run("should return error if subscription does not exist", func() {
		providerMock := &mocks.ProviderService{}
		namespaceMock := &mocks.NamespaceService{}
		receiverMock := &mocks.ReceiverService{}
		s.dbmock.ExpectBegin()
		s.dbmock.ExpectQuery(firstSelectQuery).WillReturnRows(sqlmock.
			NewRows([]string{"urn", "namespace_id", "receiver", "match", "created_at", "updated_at", "id"}))
		s.dbmock.ExpectRollback()

		actualSubscription, err := s.repository.Update(input, namespaceMock, providerMock, receiverMock)
		s.EqualError(err, "subscription doesn't exist")
		s.Nil(actualSubscription)
	})

	s.Run("should return error if updating subscription fails", func() {
		providerMock := &mocks.ProviderService{}
		namespaceMock := &mocks.NamespaceService{}
		receiverMock := &mocks.ReceiverService{}
		s.dbmock.ExpectBegin()
		expectedRowsBeforeUpdate := sqlmock.
			NewRows([]string{"urn", "namespace_id", "receiver", "match", "created_at", "updated_at", "id"}).
			AddRow(subscription.Urn, subscription.NamespaceId,
				json.RawMessage(`[{"id":1 ,"configuration": {"channel_name":"test"}}]`), json.RawMessage(`{"foo":"bar"}`),
				subscription.CreatedAt, subscription.UpdatedAt, subscription.Id)
		s.dbmock.ExpectQuery(firstSelectQuery).WillReturnRows(expectedRowsBeforeUpdate)
		s.dbmock.ExpectExec(updateQuery).WithArgs(subscription.NamespaceId, subscription.Urn,
			subscription.Receiver, subscription.Match, AnyTime{}, AnyTime{}, subscription.Id, subscription.Id).
			WillReturnError(errors.New("random error"))
		s.dbmock.ExpectRollback()

		actualSubscription, err := s.repository.Update(input, namespaceMock, providerMock, receiverMock)
		s.EqualError(err, "random error")
		s.Nil(actualSubscription)
	})

	s.Run("should return error in fetching subscription after update", func() {
		providerMock := &mocks.ProviderService{}
		namespaceMock := &mocks.NamespaceService{}
		receiverMock := &mocks.ReceiverService{}

		s.dbmock.ExpectBegin()
		expectedRowsBeforeUpdate := sqlmock.
			NewRows([]string{"urn", "namespace_id", "receiver", "match", "created_at", "updated_at", "id"}).
			AddRow(subscription.Urn, subscription.NamespaceId,
				json.RawMessage(`[{"id":1 ,"configuration": {"channel_name":"test"}}]`), json.RawMessage(`{"foo":"bar"}`),
				subscription.CreatedAt, subscription.UpdatedAt, subscription.Id)
		s.dbmock.ExpectQuery(firstSelectQuery).WillReturnRows(expectedRowsBeforeUpdate)
		s.dbmock.ExpectExec(updateQuery).WithArgs(subscription.NamespaceId, subscription.Urn,
			subscription.Receiver, subscription.Match, AnyTime{}, AnyTime{}, subscription.Id, subscription.Id).
			WillReturnResult(sqlmock.NewResult(1, 1))

		s.dbmock.ExpectQuery(secondSelectQuery).WillReturnError(errors.New("random error"))
		s.dbmock.ExpectRollback()

		actualSubscription, err := s.repository.Update(input, namespaceMock, providerMock, receiverMock)
		s.EqualError(err, "random error")
		s.Nil(actualSubscription)
	})

	s.Run("should return error in syncing alertmanager config", func() {
		providerMock := &mocks.ProviderService{}
		namespaceMock := &mocks.NamespaceService{}
		receiverMock := &mocks.ReceiverService{}
		amClientMock := &ClientMock{}
		dummySlackReceivers := &domain.Receiver{Id: 1, Type: "slack", Configurations: randomSlackReceiverConfig}
		dummyPagerdutyReceivers := &domain.Receiver{Id: 2, Type: "pagerduty", Configurations: randomPagerdutyReceiverConfig}
		dummyHTTPReceivers := &domain.Receiver{Id: 3, Type: "http", Configurations: randomHTTPReceiverConfig}
		providerMock.On("GetProvider", uint64(1)).Return(dummyProvider, nil).Once()
		namespaceMock.On("GetNamespace", uint64(1)).Return(dummyNamespace, nil).Once()
		receiverMock.On("ListReceivers").Return(
			[]*domain.Receiver{dummySlackReceivers, dummyPagerdutyReceivers, dummyHTTPReceivers}, nil).Once()
		amClientMock.On("SyncConfig", mock.AnythingOfType("alertmanager.AMConfig"), "dummy").
			Return(errors.New("random error")).Once()

		oldAMClientCreator := alertmanagerClientCreator
		defer func() { alertmanagerClientCreator = oldAMClientCreator }()
		alertmanagerClientCreator = func(c domain.CortexConfig) (alertmanager.Client, error) {
			return amClientMock, nil
		}

		s.dbmock.ExpectBegin()
		expectedRowsBeforeUpdate := sqlmock.
			NewRows([]string{"urn", "namespace_id", "receiver", "match", "created_at", "updated_at", "id"}).
			AddRow(subscription.Urn, subscription.NamespaceId,
				json.RawMessage(`[{"id":1 ,"configuration": {"channel_name":"test"}}]`), json.RawMessage(`{"foo":"bar"}`),
				subscription.CreatedAt, subscription.UpdatedAt, subscription.Id)
		s.dbmock.ExpectQuery(firstSelectQuery).WillReturnRows(expectedRowsBeforeUpdate)
		s.dbmock.ExpectExec(updateQuery).WithArgs(subscription.NamespaceId, subscription.Urn,
			subscription.Receiver, subscription.Match, AnyTime{}, AnyTime{}, subscription.Id, subscription.Id).
			WillReturnResult(sqlmock.NewResult(1, 1))

		expectedRowsAfterUpdate := sqlmock.
			NewRows([]string{"urn", "namespace_id", "receiver", "match", "created_at", "updated_at", "id"}).
			AddRow(subscription.Urn, subscription.NamespaceId,
				json.RawMessage(`[{"id":1 ,"configuration": {"channel_name":"updated_channel"}}]`), json.RawMessage(`{"foo":"bar"}`),
				subscription.CreatedAt, subscription.UpdatedAt, subscription.Id)

		s.dbmock.ExpectQuery(secondSelectQuery).WillReturnRows(expectedRowsAfterUpdate)

		expectedRowsInNamespace := sqlmock.
			NewRows([]string{"urn", "namespace_id", "receiver", "match", "created_at", "updated_at", "id"}).
			AddRow(subscription.Urn, subscription.NamespaceId,
				json.RawMessage(`[{"id":1 ,"configuration": {"channel_name": "updated_channel"}}]`),
				json.RawMessage(`{"foo": "bar"}`), subscription.CreatedAt, subscription.UpdatedAt, subscription.Id).
			AddRow("bar", subscription.NamespaceId,
				json.RawMessage(`[{"id":2 ,"configuration": {"service_key": "test"}}]`),
				json.RawMessage(`{"foo": "bar"}`), subscription.CreatedAt, subscription.UpdatedAt, 2).
			AddRow("baz", subscription.NamespaceId,
				json.RawMessage(`[{"id":3 ,"configuration": {"url": "http://localhost:3000"}}]`),
				json.RawMessage(`{"foo": "bar"}`), subscription.CreatedAt, subscription.UpdatedAt, 3)

		s.dbmock.ExpectQuery(fetchSubscriptionsWithinNamespaceQuery).WillReturnRows(expectedRowsInNamespace)
		s.dbmock.ExpectCommit()

		actualSubscription, err := s.repository.Update(input, namespaceMock, providerMock, receiverMock)
		s.EqualError(err, "r.amClient.SyncConfig: random error")
		s.Nil(actualSubscription)
	})
}

func (s *RepositoryTestSuite) TestDelete() {
	match := make(StringStringMap)
	inputRandomConfig := make(StringStringMap)
	randomSlackReceiverConfig := make(map[string]interface{})
	randomPagerdutyReceiverConfig := make(map[string]interface{})
	randomHTTPReceiverConfig := make(map[string]interface{})
	match["foo"] = "bar"
	inputRandomConfig["channel_name"] = "test"
	randomSlackReceiverConfig["token"] = "xoxb"
	randomPagerdutyReceiverConfig["service_key"] = "abcd"
	randomHTTPReceiverConfig["url"] = "http://localhost:3000"
	receiver := ReceiverMetadata{Id: 1, Configuration: inputRandomConfig}
	subscription := &Subscription{
		Id:          1,
		NamespaceId: 1,
		Urn:         "foo",
		Match:       match,
		Receiver:    []ReceiverMetadata{receiver},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	inputRandomConfig["channel_name"] = "updated_channel"

	dummyNamespace := &domain.Namespace{Id: 1, Provider: 1, Urn: "dummy"}
	dummyProvider := &domain.Provider{Id: 1, Urn: "test", Type: "cortex", Host: "http://localhost:8080"}

	selectQuery := regexp.QuoteMeta(`SELECT * FROM "subscriptions" WHERE id = 1`)
	deleteQuery := regexp.QuoteMeta(`DELETE FROM "subscriptions" WHERE "subscriptions"."id" = $1`)
	fetchSubscriptionsWithinNamespaceQuery := regexp.QuoteMeta(`SELECT * FROM "subscriptions" WHERE namespace_id = 1`)

	s.Run("should delete a subscription", func() {
		providerMock := &mocks.ProviderService{}
		namespaceMock := &mocks.NamespaceService{}
		receiverMock := &mocks.ReceiverService{}
		amClientMock := &ClientMock{}
		dummySlackReceivers := &domain.Receiver{Id: 1, Type: "slack", Configurations: randomSlackReceiverConfig}
		dummyPagerdutyReceivers := &domain.Receiver{Id: 2, Type: "pagerduty", Configurations: randomPagerdutyReceiverConfig}
		dummyHTTPReceivers := &domain.Receiver{Id: 3, Type: "http", Configurations: randomHTTPReceiverConfig}
		providerMock.On("GetProvider", uint64(1)).Return(dummyProvider, nil).Once()
		namespaceMock.On("GetNamespace", uint64(1)).Return(dummyNamespace, nil).Once()
		receiverMock.On("ListReceivers").Return(
			[]*domain.Receiver{dummySlackReceivers, dummyPagerdutyReceivers, dummyHTTPReceivers}, nil).Once()
		amClientMock.On("SyncConfig", mock.AnythingOfType("alertmanager.AMConfig"), "dummy").Run(func(args mock.Arguments) {
			rarg := args.Get(0)
			s.Require().IsType(alertmanager.AMConfig{}, rarg)
			r := rarg.(alertmanager.AMConfig)
			s.Equal(2, len(r.Receivers))
			s.Equal("bar_receiverId_2_idx_0", r.Receivers[0].Receiver)
			s.Equal("baz_receiverId_3_idx_0", r.Receivers[1].Receiver)
		}).Return(nil).Once()

		oldAMClientCreator := alertmanagerClientCreator
		defer func() { alertmanagerClientCreator = oldAMClientCreator }()
		alertmanagerClientCreator = func(c domain.CortexConfig) (alertmanager.Client, error) {
			return amClientMock, nil
		}

		s.dbmock.ExpectBegin()
		expectedRows := sqlmock.
			NewRows([]string{"urn", "namespace_id", "receiver", "match", "created_at", "updated_at", "id"}).
			AddRow(subscription.Urn, subscription.NamespaceId,
				json.RawMessage(`[{"id":1 ,"configuration": {"channel_name":"test"}}]`), json.RawMessage(`{"foo":"bar"}`),
				subscription.CreatedAt, subscription.UpdatedAt, subscription.Id)
		s.dbmock.ExpectQuery(selectQuery).WillReturnRows(expectedRows)
		s.dbmock.ExpectExec(deleteQuery).WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))

		expectedRowsInNamespace := sqlmock.
			NewRows([]string{"urn", "namespace_id", "receiver", "match", "created_at", "updated_at", "id"}).
			AddRow("bar", subscription.NamespaceId,
				json.RawMessage(`[{"id":2 ,"configuration": {"service_key": "test"}}]`),
				json.RawMessage(`{"foo": "bar"}`), subscription.CreatedAt, subscription.UpdatedAt, 2).
			AddRow("baz", subscription.NamespaceId,
				json.RawMessage(`[{"id":3 ,"configuration": {"url": "http://localhost:3000"}}]`),
				json.RawMessage(`{"foo": "bar"}`), subscription.CreatedAt, subscription.UpdatedAt, 3)

		s.dbmock.ExpectQuery(fetchSubscriptionsWithinNamespaceQuery).WillReturnRows(expectedRowsInNamespace)
		s.dbmock.ExpectCommit()

		err := s.repository.Delete(1, namespaceMock, providerMock, receiverMock)
		s.Nil(err)
	})

	s.Run("should return error in fetching subscription", func() {
		providerMock := &mocks.ProviderService{}
		namespaceMock := &mocks.NamespaceService{}
		receiverMock := &mocks.ReceiverService{}

		s.dbmock.ExpectBegin()
		s.dbmock.ExpectQuery(selectQuery).WillReturnError(errors.New("random error"))
		s.dbmock.ExpectRollback()

		err := s.repository.Delete(1, namespaceMock, providerMock, receiverMock)
		s.EqualError(err, "random error")
	})

	s.Run("should return no error if subscription does not exist", func() {
		providerMock := &mocks.ProviderService{}
		namespaceMock := &mocks.NamespaceService{}
		receiverMock := &mocks.ReceiverService{}

		s.dbmock.ExpectBegin()
		expectedRows := sqlmock.
			NewRows([]string{"urn", "namespace_id", "receiver", "match", "created_at", "updated_at", "id"})
		s.dbmock.ExpectQuery(selectQuery).WillReturnRows(expectedRows)
		s.dbmock.ExpectCommit()

		err := s.repository.Delete(1, namespaceMock, providerMock, receiverMock)
		s.Nil(err)
	})

	s.Run("should return error in deleting subscription", func() {
		providerMock := &mocks.ProviderService{}
		namespaceMock := &mocks.NamespaceService{}
		receiverMock := &mocks.ReceiverService{}
		s.dbmock.ExpectBegin()
		expectedRows := sqlmock.
			NewRows([]string{"urn", "namespace_id", "receiver", "match", "created_at", "updated_at", "id"}).
			AddRow(subscription.Urn, subscription.NamespaceId,
				json.RawMessage(`[{"id":1 ,"configuration": {"channel_name":"test"}}]`), json.RawMessage(`{"foo":"bar"}`),
				subscription.CreatedAt, subscription.UpdatedAt, subscription.Id)
		s.dbmock.ExpectQuery(selectQuery).WillReturnRows(expectedRows)
		s.dbmock.ExpectExec(deleteQuery).WithArgs(1).WillReturnError(errors.New("random error"))
		s.dbmock.ExpectRollback()

		err := s.repository.Delete(1, namespaceMock, providerMock, receiverMock)
		s.EqualError(err, "r.deleteSubscriptionFromDB: random error")
	})

	s.Run("should return error in syncing configuration with alertmanager", func() {
		providerMock := &mocks.ProviderService{}
		namespaceMock := &mocks.NamespaceService{}
		receiverMock := &mocks.ReceiverService{}
		amClientMock := &ClientMock{}
		dummySlackReceivers := &domain.Receiver{Id: 1, Type: "slack", Configurations: randomSlackReceiverConfig}
		dummyPagerdutyReceivers := &domain.Receiver{Id: 2, Type: "pagerduty", Configurations: randomPagerdutyReceiverConfig}
		dummyHTTPReceivers := &domain.Receiver{Id: 3, Type: "http", Configurations: randomHTTPReceiverConfig}
		providerMock.On("GetProvider", uint64(1)).Return(dummyProvider, nil).Once()
		namespaceMock.On("GetNamespace", uint64(1)).Return(dummyNamespace, nil).Once()
		receiverMock.On("ListReceivers").Return(
			[]*domain.Receiver{dummySlackReceivers, dummyPagerdutyReceivers, dummyHTTPReceivers}, nil).Once()
		amClientMock.On("SyncConfig", mock.AnythingOfType("alertmanager.AMConfig"), "dummy").
			Return(errors.New("random error")).Once()

		oldAMClientCreator := alertmanagerClientCreator
		defer func() { alertmanagerClientCreator = oldAMClientCreator }()
		alertmanagerClientCreator = func(c domain.CortexConfig) (alertmanager.Client, error) {
			return amClientMock, nil
		}

		s.dbmock.ExpectBegin()
		expectedRows := sqlmock.
			NewRows([]string{"urn", "namespace_id", "receiver", "match", "created_at", "updated_at", "id"}).
			AddRow(subscription.Urn, subscription.NamespaceId,
				json.RawMessage(`[{"id":1 ,"configuration": {"channel_name":"test"}}]`), json.RawMessage(`{"foo":"bar"}`),
				subscription.CreatedAt, subscription.UpdatedAt, subscription.Id)
		s.dbmock.ExpectQuery(selectQuery).WillReturnRows(expectedRows)
		s.dbmock.ExpectExec(deleteQuery).WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))

		expectedRowsInNamespace := sqlmock.
			NewRows([]string{"urn", "namespace_id", "receiver", "match", "created_at", "updated_at", "id"}).
			AddRow("bar", subscription.NamespaceId,
				json.RawMessage(`[{"id":2 ,"configuration": {"service_key": "test"}}]`),
				json.RawMessage(`{"foo": "bar"}`), subscription.CreatedAt, subscription.UpdatedAt, 2).
			AddRow("baz", subscription.NamespaceId,
				json.RawMessage(`[{"id":3 ,"configuration": {"url": "http://localhost:3000"}}]`),
				json.RawMessage(`{"foo": "bar"}`), subscription.CreatedAt, subscription.UpdatedAt, 3)

		s.dbmock.ExpectQuery(fetchSubscriptionsWithinNamespaceQuery).WillReturnRows(expectedRowsInNamespace)
		s.dbmock.ExpectRollback()

		err := s.repository.Delete(1, namespaceMock, providerMock, receiverMock)
		s.EqualError(err, "r.amClient.SyncConfig: random error")
	})
}

func TestRepository(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}
