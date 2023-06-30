package errors_test

import (
	goerrors "errors"
	"testing"

	"github.com/raystack/siren/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func Test_E(t *testing.T) {
	t.Parallel()

	t.Run("NotError", func(t *testing.T) {
		want := errors.Error{
			Code:    "internal_error",
			Cause:   "some native error",
			Message: "Some unexpected error occurred",
		}

		err := errors.New("some native error")
		got := errors.E(err)

		assert.Equal(t, want, got)
	})

	t.Run("ErrorValue", func(t *testing.T) {
		want := errors.Error{
			Code:    "not_found",
			Cause:   "foo not found",
			Message: "something was not found",
		}

		err := error(errors.ErrNotFound.
			WithMsgf("something was not found").
			WithCausef("foo not found"))

		got := errors.E(err)
		assert.Equal(t, want, got)
	})
}

func Test_Verbose(t *testing.T) {
	t.Parallel()

	t.Run("NonError", func(t *testing.T) {
		err := errors.New("some native error")

		got := errors.Verbose(err)
		assert.EqualError(t, got, "some native error")
	})

	t.Run("CustomError", func(t *testing.T) {
		err := errors.ErrInvalid.
			WithMsgf("request is not valid").
			WithCausef("invalid parameter value")

		got := errors.Verbose(err)
		assert.EqualError(t, got, "bad_request: request is not valid (cause: invalid parameter value)")
	})
}

func Test_Errorf(t *testing.T) {
	t.Parallel()
	e := errors.Errorf("failed: %d", 100)
	assert.Error(t, e)
	assert.EqualError(t, e, "failed: 100")
}

func Test_OneOf(t *testing.T) {
	t.Parallel()

	group := []error{
		errors.ErrNotFound,
		errors.ErrInvalid,
		errors.ErrUnsupported,
	}

	assert.False(t, errors.OneOf(errors.New("failed"), group...))
	assert.True(t, errors.OneOf(errors.ErrNotFound.WithMsgf("object not found"), group...))
}

func TestError_Error(t *testing.T) {
	t.Parallel()

	table := []struct {
		title string
		err   errors.Error
		want  string
	}{
		{
			title: "WithoutCause",
			err:   errors.ErrInvalid,
			want:  "request is not valid",
		},
		{
			title: "WithCause",
			err:   errors.ErrInvalid.WithMsgf("").WithCausef("foo"),
			want:  "bad_request: foo",
		},
	}

	for _, tt := range table {
		tt := tt
		t.Run(tt.title, func(t *testing.T) {
			t.Parallel()
			got := tt.err.Error()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestError_Is(t *testing.T) {
	t.Parallel()

	table := []struct {
		title string
		err   errors.Error
		other error
		want  bool
	}{
		{
			title: "WithDifferentCode",
			err:   errors.ErrInternal,
			other: errors.ErrInvalid,
			want:  false,
		},
		{
			title: "NonEntropyErr",
			err:   errors.ErrInternal,
			other: goerrors.New("foo"), // nolint
			want:  true,
		},
		{
			title: "WithSameCode",
			err:   errors.ErrInvalid.WithCausef("cause 1"),
			other: errors.ErrInvalid.WithCausef("cause 2"),
			want:  true,
		},
	}

	for _, tt := range table {
		tt := tt
		t.Run(tt.title, func(t *testing.T) {
			t.Parallel()
			got := goerrors.Is(tt.err, tt.other)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestError_WithCausef(t *testing.T) {
	t.Parallel()

	table := []struct {
		title string
		err   errors.Error
		want  errors.Error
	}{
		{
			title: "WithCauseString",
			err:   errors.ErrInvalid.WithCausef("foo"),
			want: errors.Error{
				Code:    "bad_request",
				Message: "Request is not valid",
				Cause:   "foo",
			},
		},
		{
			title: "WithCauseFormatted",
			err:   errors.ErrConflict.WithCausef("hello %s", "world"),
			want: errors.Error{
				Code:    "conflict",
				Message: "An entity with conflicting identifier exists",
				Cause:   "hello world",
			},
		},
	}

	for _, tt := range table {
		tt := tt
		t.Run(tt.title, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, tt.err)
		})
	}
}

func TestError_WithMsgf(t *testing.T) {
	t.Parallel()

	table := []struct {
		title string
		err   errors.Error
		want  errors.Error
	}{
		{
			title: "WithCauseString",
			err:   errors.ErrInvalid.WithMsgf("foo"),
			want: errors.Error{
				Code:    "bad_request",
				Message: "foo",
			},
		},
		{
			title: "WithCauseFormatted",
			err:   errors.ErrInvalid.WithMsgf("hello %s", "world"),
			want: errors.Error{
				Code:    "bad_request",
				Message: "hello world",
			},
		},
	}

	for _, tt := range table {
		tt := tt
		t.Run(tt.title, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, tt.err)
		})
	}
}
