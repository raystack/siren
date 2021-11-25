package subscription

import (
	"fmt"
	"github.com/odpf/siren/domain"
	"github.com/odpf/siren/pkg/subscription/alertmanager"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"sort"
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

func (r Repository) List() ([]*Subscription, error) {
	var subscriptions []*Subscription
	selectQuery := fmt.Sprintf("select * from subscriptions")
	result := r.db.Raw(selectQuery).Find(&subscriptions)
	if result.Error != nil {
		return nil, result.Error
	}

	return subscriptions, nil
}

func (r Repository) Create(sub *Subscription, namespaceService domain.NamespaceService,
	providerService domain.ProviderService, receiverService domain.ReceiverService) (*Subscription, error) {
	var newSubscription *Subscription
	var err error
	createError := r.db.Transaction(func(tx *gorm.DB) error {
		newSubscription, err = r.insertSubscriptionIntoDB(tx, sub)
		if err != nil {
			return errors.Wrap(err, "r.insertSubscriptionIntoDB")
		}
		return r.syncInUpstreamCurrentSubscriptionsOfNamespace(tx, newSubscription.NamespaceId,
			namespaceService, providerService, receiverService)
	})
	return newSubscription, createError
}

func (r Repository) Get(id uint64) (*Subscription, error) {
	var subscription Subscription
	result := r.db.Where(fmt.Sprintf("id = %d", id)).Find(&subscription)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}

	return &subscription, nil
}

func (r Repository) Update(sub *Subscription, namespaceService domain.NamespaceService,
	providerService domain.ProviderService, receiverService domain.ReceiverService) (*Subscription, error) {
	var newSubscription, existingSubscription Subscription
	updateError := r.db.Transaction(func(tx *gorm.DB) error {
		result := r.db.Where(fmt.Sprintf("id = %d", sub.Id)).Find(&existingSubscription)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errors.New("subscription doesn't exist")
		} else {
			sortReceivers(sub)
			result = r.db.Where("id = ?", sub.Id).Updates(sub)
			if result.Error != nil {
				return result.Error
			}
		}
		result = r.db.Where(fmt.Sprintf("id = %d", sub.Id)).Find(&newSubscription)
		if result.Error != nil {
			return result.Error
		}
		return r.syncInUpstreamCurrentSubscriptionsOfNamespace(tx, sub.NamespaceId, namespaceService,
			providerService, receiverService)
	})
	return &newSubscription, updateError
}

func (r Repository) Delete(id uint64, namespaceService domain.NamespaceService, providerService domain.ProviderService,
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

func (r Repository) Migrate() error {
	err := r.db.AutoMigrate(&Subscription{})
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) syncInUpstreamCurrentSubscriptionsOfNamespace(tx *gorm.DB, namespaceId uint64, namespaceService domain.NamespaceService,
	providerService domain.ProviderService, receiverService domain.ReceiverService) error {
	// fetch all subscriptions in this namespace.
	subscriptionsInNamespace, err := r.getAllSubscriptionsWithinNamespace(tx, namespaceId)
	if err != nil {
		return errors.Wrap(err, "r.getAllSubscriptionsWithinNamespace")
	}
	// check provider type of the namespace
	providerInfo, namespaceInfo, err := r.getProviderAndNamespaceInfoFromNamespaceId(namespaceId, namespaceService, providerService)
	if err != nil {
		return err
	}
	subscriptionsInNamespaceEnrichedWithReceivers, err := r.addReceiversConfiguration(subscriptionsInNamespace, receiverService)
	if err != nil {
		return err
	}
	// do upstream call to create subscriptions as per provider type
	switch providerInfo.Type {
	case "cortex":
		amConfig := getAmConfigFromSubscriptions(subscriptionsInNamespaceEnrichedWithReceivers)
		newAMClient, err := alertmanager.NewClient(domain.CortexConfig{Address: providerInfo.Host})
		if err != nil {
			return errors.Wrap(err, "failed to initialize alertmanager client")
		}
		r.amClient = newAMClient
		err = r.amClient.SyncConfig(amConfig, namespaceInfo.Urn)
		if err != nil {
			return err
		}
	default:
		return errors.New(fmt.Sprintf("subscriptions for provider type '%s' not supported", providerInfo.Type))
	}
	return nil
}

func (r Repository) getProviderAndNamespaceInfoFromNamespaceId(id uint64, namespaceService domain.NamespaceService,
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

func (r Repository) insertSubscriptionIntoDB(tx *gorm.DB, sub *Subscription) (*Subscription, error) {
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

func (r Repository) deleteSubscriptionFromDB(tx *gorm.DB, id uint64) error {
	result := tx.Delete(Subscription{}, id)
	return result.Error
}
