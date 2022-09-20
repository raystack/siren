package subscription

type ServiceOption func(*Service)

func RegisterProviderPlugin(typeName string, service SubscriptionSyncer) ServiceOption {
	return func(s *Service) {
		s.subscriptionProviderRegistry[typeName] = service
	}
}
