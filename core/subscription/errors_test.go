package subscription_test

import (
	"testing"

	"github.com/odpf/siren/core/subscription"
)

func TestError(t *testing.T) {
	t.Run("should return error with id if id is not empty", func(t *testing.T) {
		expectedErrString := "subscription with id 1 not found"
		err := subscription.NotFoundError{ID: 1}
		if err.Error() != expectedErrString {
			t.Fatalf("got error %v, expected was %v", err.Error(), expectedErrString)
		}
	})

	t.Run("should return error with no id if id is empty", func(t *testing.T) {
		expectedErrString := "subscription not found"
		err := subscription.NotFoundError{}
		if err.Error() != expectedErrString {
			t.Fatalf("got error %v, expected was %v", err.Error(), expectedErrString)
		}
	})
}
