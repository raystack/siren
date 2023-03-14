package provider_test

import (
	"testing"

	"github.com/goto/siren/core/provider"
)

func TestError(t *testing.T) {
	t.Run("should return error with id if id is not empty", func(t *testing.T) {
		expectedErrString := "provider with id 1 not found"
		err := provider.NotFoundError{ID: 1}
		if err.Error() != expectedErrString {
			t.Fatalf("got error %v, expected was %v", err.Error(), expectedErrString)
		}
	})

	t.Run("should return error with no id if id is empty", func(t *testing.T) {
		expectedErrString := "provider not found"
		err := provider.NotFoundError{}
		if err.Error() != expectedErrString {
			t.Fatalf("got error %v, expected was %v", err.Error(), expectedErrString)
		}
	})
}
