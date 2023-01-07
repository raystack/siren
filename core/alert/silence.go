package alert

const (
	SilenceStatusTotal   = "total"
	SilenceStatusPartial = "partial"
)

func silenceStatus(hasSilenced, hasNonSilenced bool) string {
	if hasSilenced && !hasNonSilenced {
		return SilenceStatusTotal
	} else if hasSilenced && hasNonSilenced {
		return SilenceStatusPartial
	}
	return ""
}
