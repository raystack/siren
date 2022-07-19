package namespace

import "testing"

func TestError(t *testing.T) {
	t.Run("should return error with id if id is not empty", func(t *testing.T) {
		expectedErrString := "namespace with id 1 not found"
		err := NotFoundError{ID: 1}
		if err.Error() != expectedErrString {
			t.Fatalf("got error %v, expected was %v", err.Error(), expectedErrString)
		}
	})

	t.Run("should return error with no id if id is empty", func(t *testing.T) {
		expectedErrString := "namespace not found"
		err := NotFoundError{}
		if err.Error() != expectedErrString {
			t.Fatalf("got error %v, expected was %v", err.Error(), expectedErrString)
		}
	})
}
