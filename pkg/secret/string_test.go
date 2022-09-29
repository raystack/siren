package secret

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	testString := "testtest"

	t.Run("should print masked string by default", func(t *testing.T) {
		maskedString := MaskableString(testString)
		result := fmt.Sprintf("%v", maskedString)
		assert.Equal(t, result, strings.Repeat("*", len(testString)))
	})

	t.Run("should print unmasked string with unmasked function", func(t *testing.T) {
		maskedString := MaskableString(testString)
		result := fmt.Sprintf("%v", maskedString.UnmaskedString())
		assert.Equal(t, result, testString)
	})
}
