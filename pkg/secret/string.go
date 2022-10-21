package secret

import "strings"

type MaskableString string

func (m MaskableString) UnmaskedString() string {
	return string(m)
}

func (m MaskableString) String() string {
	return strings.Repeat("*", len(m))
}
