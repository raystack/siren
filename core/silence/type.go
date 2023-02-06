package silence

const (
	TypeMatchers     = "Matchers"
	TypeSubscription = "subscription"
)

func IsTypeValid(silenceTypeStr string) bool {
	if silenceTypeStr == TypeMatchers ||
		silenceTypeStr == TypeSubscription {
		return true
	}
	return false
}
