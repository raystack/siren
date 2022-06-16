package receiver_test

import (
	"testing"

	"github.com/odpf/siren/core/receiver"
)

func TestError(t *testing.T) {
	t.Run("should return error with id if id is not empty", func(t *testing.T) {
		expectedErrString := "receiver with id 1 not found"
		err := receiver.NotFoundError{ID: 1}
		if err.Error() != expectedErrString {
			t.Fatalf("got error %v, expected was %v", err.Error(), expectedErrString)
		}
	})

	t.Run("should return error with no id if id is empty", func(t *testing.T) {
		expectedErrString := "receiver not found"
		err := receiver.NotFoundError{}
		if err.Error() != expectedErrString {
			t.Fatalf("got error %v, expected was %v", err.Error(), expectedErrString)
		}
	})
}
