package receiver

const (
	TypeSlack     string = "slack"
	TypeHTTP      string = "http"
	TypePagerDuty string = "pagerduty"
	TypeFile      string = "file"
)

var SupportedTypes = []string{
	TypeSlack,
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
