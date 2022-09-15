package subscription

type ServiceOption func(*Service)

func RegisterProviderPlugin(typeName string, service ProviderPlugin) ServiceOption {
	return func(s *Service) {
		s.subscriptionProviderRegistry[typeName] = service
	}
}
