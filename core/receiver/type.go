package receiver

const (
	TypeSlack        string = "slack"
	TypeSlackChannel string = "slack_channel"
	TypeHTTP         string = "http"
	TypePagerDuty    string = "pagerduty"
	TypeFile         string = "file"
)

var SupportedTypes = []string{
	TypeSlack,
	TypeSlackChannel,
	TypeHTTP,
	TypePagerDuty,
	TypeFile,
}

func IsTypeSupported(receiverType string) bool {
	for _, st := range SupportedTypes {
		if st == receiverType {
			return true
		}
	}
	return false
}
