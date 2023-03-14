package template_test

import (
	"testing"

	"github.com/goto/siren/core/template"
)

func TestError(t *testing.T) {
	t.Run("should return error with name if name is not empty", func(t *testing.T) {
		expectedErrString := "template with name \"some-name\" not found"
		err := template.NotFoundError{Name: "some-name"}
		if err.Error() != expectedErrString {
			t.Fatalf("got error %v, expected was %v", err.Error(), expectedErrString)
		}
	})

	t.Run("should return error with no name if name is empty", func(t *testing.T) {
		expectedErrString := "template not found"
		err := template.NotFoundError{}
		if err.Error() != expectedErrString {
			t.Fatalf("got error %v, expected was %v", err.Error(), expectedErrString)
		}
	})
}
