package subscription

import (
	"context"
	"fmt"
	"sort"

	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/pkg/subscription/alertmanager"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var (
	transactionContextKey = struct{}{}
)

// Repository talks to the store to read or insert data
type Repository struct {
	db       *gorm.DB
	amClient alertmanager.Client
}

// NewRepository returns repository struct
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db, nil}
}

func (r *Repository) WithTransaction(ctx context.Context) context.Context {
	tx := r.db.Begin()
	return context.WithValue(ctx, transactionContextKey, tx)
}

func getTransaction(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(transactionContextKey).(*gorm.DB); !ok {
		return nil
	} else {
		return tx
	}
}

func (r *Repository) Rollback(ctx context.Context) error {
	if tx := getTransaction(ctx); tx != nil {
		tx = tx.Rollback()
		if tx.Error != nil {
			return r.db.Error
		}
		return nil
	}
	return errors.New("no transaction")
}

func (r *Repository) Commit(ctx context.Context) error {
	if tx := getTransaction(ctx); tx != nil {
		tx = tx.Commit()
		if tx.Error != nil {
			return r.db.Error
		}
		return nil
	}
	return errors.New("no transaction")
}

func (r *Repository) List(ctx context.Context) ([]*domain.Subscription, error) {
	var subscriptionModels []*Subscription
	selectQuery := "select * from subscriptions"
	result := r.db.Raw(selectQuery).Find(&subscriptionModels)
	if result.Error != nil {
		return nil, result.Error
	}

	var subscriptions []*domain.Subscription
	for _, s := range subscriptionModels {
		subscriptions = append(subscriptions, s.toDomain())
	}

	return subscriptions, nil
}

func (r *Repository) Create(ctx context.Context, sub *domain.Subscription) (*domain.Subscription, error) {
	var model Subscription
	model.fromDomain(sub)
	var newSubscription *Subscription
	var err error

	db := r.db
	if tx := getTransaction(ctx); tx != nil {
		db = tx
	}

	newSubscription, err = r.insertSubscriptionIntoDB(db, &model)
	if err != nil {
		return nil, errors.Wrap(err, "r.insertSubscriptionIntoDB")
	}

	return newSubscription.toDomain(), nil
}

func (r *Repository) Get(ctx context.Context, id uint64) (*domain.Subscription, error) {
	var subscription Subscription
	result := r.db.Where(fmt.Sprintf("id = %d", id)).Find(&subscription)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}

	return subscription.toDomain(), nil
}

func (r *Repository) Update(ctx context.Context, sub *domain.Subscription) (*domain.Subscription, error) {
	db := r.db
	if tx := getTransaction(ctx); tx != nil {
		db = tx
	}

	model := new(Subscription)
	model.fromDomain(sub)
	result := db.Where("id = ?", model.Id).Updates(model)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("subscription doesn't exist")
	}
	return model.toDomain(), nil
}

func (r *Repository) Delete(ctx context.Context, id uint64, namespaceService domain.NamespaceService, providerService domain.ProviderService,
	receiverService domain.ReceiverService) error {
	deleteError := r.db.Transaction(func(tx *gorm.DB) error {
		var subscription Subscription
		result := r.db.Where(fmt.Sprintf("id = %d", id)).Find(&subscription)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return nil
		}
		err := r.deleteSubscriptionFromDB(tx, id)
		if err != nil {
			return errors.Wrap(err, "r.deleteSubscriptionFromDB")
		}
		return r.syncInUpstreamCurrentSubscriptionsOfNamespace(tx, subscription.NamespaceId,
			namespaceService, providerService, receiverService)
	})
	return deleteError
}

func (r *Repository) Migrate() error {
	err := r.db.AutoMigrate(&Subscription{})
	if err != nil {
		return err
	}
	return nil
}

var alertmanagerClientCreator = alertmanager.NewClient

func (r *Repository) syncInUpstreamCurrentSubscriptionsOfNamespace(tx *gorm.DB, namespaceId uint64, namespaceService domain.NamespaceService,
	providerService domain.ProviderService, receiverService domain.ReceiverService) error {
	// fetch all subscriptions in this namespace.
	subscriptionsInNamespace, err := r.getAllSubscriptionsWithinNamespace(tx, namespaceId)
	if err != nil {
		return errors.Wrap(err, "r.getAllSubscriptionsWithinNamespace")
	}
	// check provider type of the namespace
	providerInfo, namespaceInfo, err := r.getProviderAndNamespaceInfoFromNamespaceId(namespaceId, namespaceService, providerService)
	if err != nil {
		return errors.Wrap(err, "r.getProviderAndNamespaceInfoFromNamespaceId")
	}
	subscriptionsInNamespaceEnrichedWithReceivers, err := r.addReceiversConfiguration(subscriptionsInNamespace, receiverService)
	if err != nil {
		return errors.Wrap(err, "r.addReceiversConfiguration")
	}
	// do upstream call to create subscriptions as per provider type
	switch providerInfo.Type {
	case "cortex":
		amConfig := getAmConfigFromSubscriptions(subscriptionsInNamespaceEnrichedWithReceivers)
		newAMClient, err := alertmanagerClientCreator(domain.CortexConfig{Address: providerInfo.Host})
		if err != nil {
			return errors.Wrap(err, "alertmanagerClientCreator: ")
		}
		r.amClient = newAMClient
		err = r.amClient.SyncConfig(amConfig, namespaceInfo.Urn)
		if err != nil {
			return errors.Wrap(err, "r.amClient.SyncConfig")
		}
	default:
		return errors.New(fmt.Sprintf("subscriptions for provider type '%s' not supported", providerInfo.Type))
	}
	return nil
}

func (r *Repository) getProviderAndNamespaceInfoFromNamespaceId(id uint64, namespaceService domain.NamespaceService,
	providerService domain.ProviderService) (*domain.Provider, *domain.Namespace, error) {
	namespaceInfo, err := namespaceService.GetNamespace(id)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to get namespace details")
	}
	providerInfo, err := providerService.GetProvider(namespaceInfo.Provider)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to get provider details")
	}
	return providerInfo, namespaceInfo, nil
}

func (r *Repository) insertSubscriptionIntoDB(tx *gorm.DB, sub *Subscription) (*Subscription, error) {
	var newSubscription Subscription
	sortReceivers(sub)
	result := tx.Create(sub)
	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "failed to insert subscription")
	}

	result = tx.Where(fmt.Sprintf("id = %d", sub.Id)).Find(&newSubscription)
	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "failed to get newly inserted subscription")
	}
	return &newSubscription, nil
}

func sortReceivers(sub *Subscription) {
	sort.Slice(sub.Receiver, func(i, j int) bool {
		return sub.Receiver[i].Id < sub.Receiver[j].Id
	})
}

func (r Repository) getAllSubscriptionsWithinNamespace(tx *gorm.DB, id uint64) ([]Subscription, error) {
	var subscriptionsInNamespace []Subscription
	result := tx.Where(fmt.Sprintf("namespace_id = %d", id)).Find(&subscriptionsInNamespace)
	if result.Error != nil {
		return nil, result.Error
	}
	return subscriptionsInNamespace, nil
}

func (r *Repository) deleteSubscriptionFromDB(tx *gorm.DB, id uint64) error {
	result := tx.Delete(Subscription{}, id)
	return result.Error
}
