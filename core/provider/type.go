package provider

const (
	TypeCortex string = "cortex"
)

var SupportedTypes = []string{
	TypeCortex,
}

func IsTypeSupported(providerType string) bool {
	for _, st := range SupportedTypes {
		if st == providerType {
			return true
		}
	}
	return false
}
