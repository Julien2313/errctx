package errctx

import (
	"errors"
	"fmt"
	"testing"
)

func TestIs(t *testing.T) {
	baseErr := errors.New("base error")
	targetErr := errors.New("target error")

	tests := []struct {
		name     string
		err      error
		target   error
		expected bool
	}{
		{
			name:     "simple error match",
			err:      targetErr,
			target:   targetErr,
			expected: true,
		},
		{
			name:     "simple error no match",
			err:      baseErr,
			target:   targetErr,
			expected: false,
		},
		{
			name:     "ErrCtx with matching inner error",
			err:      Errorf("wrapped: %w", targetErr),
			target:   targetErr,
			expected: true,
		},
		{
			name:     "ErrCtx with non-matching inner error",
			err:      Errorf("wrapped: %w", baseErr),
			target:   targetErr,
			expected: false,
		},
		{
			name:     "nested ErrCtx with matching error",
			err:      Errorf("outer: %w", Errorf("inner: %w", targetErr)),
			target:   targetErr,
			expected: true,
		},
		{
			name:     "nested ErrCtx with non-matching error",
			err:      Errorf("outer: %w", Errorf("inner: %w", baseErr)),
			target:   targetErr,
			expected: false,
		},
		{
			name:     "ErrCtx matching itself",
			err:      New("custom error"),
			target:   New("custom error"),
			expected: false,
		},
		{
			name:     "nil error with nil target",
			err:      nil,
			target:   nil,
			expected: errors.Is(nil, nil),
		},
		{
			name:     "nil error with non-nil target",
			err:      nil,
			target:   targetErr,
			expected: false,
		},
		{
			name:     "non-nil error with nil target",
			err:      baseErr,
			target:   nil,
			expected: false,
		},
		{
			name:     "ErrCtx with joined errors containing target",
			err:      New("main").Join(targetErr),
			target:   targetErr,
			expected: true,
		},
		{
			name:     "ErrCtx with joined errors not containing target",
			err:      New("main").Join(baseErr),
			target:   targetErr,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println(tt.name)
			result := Is(tt.err, tt.target)
			if result != tt.expected {
				t.Errorf("Is(%v, %v) = %v, want %v", tt.err, tt.target, result, tt.expected)
			}
		})
	}
}

func TestAs(t *testing.T) {
	customErr := CustomError{Msg: "custom error"}
	anotherCustomErr := AnotherCustomError{Code: 404}

	tests := []struct {
		name     string
		err      error
		target   any
		expected bool
	}{
		{
			name:     "simple error type match",
			err:      customErr,
			target:   &CustomError{},
			expected: true,
		},
		{
			name:     "simple error type no match",
			err:      customErr,
			target:   &AnotherCustomError{},
			expected: false,
		},
		{
			name:     "ErrCtx with matching inner error type",
			err:      Errorf("wrapped: %w", customErr),
			target:   &CustomError{},
			expected: true,
		},
		{
			name:     "ErrCtx with non-matching inner error type",
			err:      Errorf("wrapped: %w", customErr),
			target:   &AnotherCustomError{},
			expected: false,
		},
		{
			name:     "nested ErrCtx with matching error type",
			err:      Errorf("outer: %w", Errorf("inner: %w", customErr)),
			target:   &CustomError{},
			expected: true,
		},
		{
			name:     "nested ErrCtx with non-matching error type",
			err:      Errorf("outer: %w", Errorf("inner: %w", customErr)),
			target:   &AnotherCustomError{},
			expected: false,
		},
		{
			name:     "ErrCtx matching itself",
			err:      New("custom error"),
			target:   &ErrCtx{},
			expected: true,
		},
		{
			name:     "nil error with nil target",
			err:      nil,
			target:   nil,
			expected: errors.As(nil, func() any { return nil }()),
		},
		{
			name:     "nil error with non-nil target",
			err:      nil,
			target:   &CustomError{},
			expected: false,
		},
		{
			name:     "ErrCtx with joined errors containing target type",
			err:      New("main").Join(customErr),
			target:   &CustomError{},
			expected: true,
		},
		{
			name:     "ErrCtx with joined errors not containing target type",
			err:      New("main").Join(customErr),
			target:   &AnotherCustomError{},
			expected: false,
		},
		{
			name:     "mixed error types in joined errors",
			err:      New("main").Join(customErr).Join(anotherCustomErr),
			target:   &CustomError{},
			expected: true,
		},
		{
			name:     "mixed error types in joined errors - second type",
			err:      New("main").Join(customErr).Join(anotherCustomErr),
			target:   &AnotherCustomError{},
			expected: true,
		},
		{
			name:     "standard error with custom error target",
			err:      errors.New("standard error"),
			target:   &CustomError{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := As(tt.err, tt.target)
			if result != tt.expected {
				t.Errorf("As(%v, %T) = %v, want %v", tt.err, tt.target, result, tt.expected)
			}
		})
	}
}

type CustomError struct {
	Msg string
}

func (e CustomError) Error() string {
	return e.Msg
}

type AnotherCustomError struct {
	Code int
}

func (e AnotherCustomError) Error() string {
	return fmt.Sprintf("error with code %d", e.Code)
}

func TestWith(t *testing.T) {
	errCtx := Errorf("custom error: %w", errors.New("base error")).With("field", "value")
	if errCtx.Value("field") != "value" {
		t.Errorf("Mismatch for key 'field' expected '%#v' got '%#v'", "value", errCtx.Value("field"))
	}

	if ErrToCtx(error(errCtx)).Value("field") != "value" {
		t.Errorf("Mismatch for key 'field' expected '%#v' got '%#v'", "value", ErrToCtx(error(errCtx)).Value("field"))
	}
}
