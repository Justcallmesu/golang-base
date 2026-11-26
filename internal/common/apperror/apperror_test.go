package apperror

import (
	"fmt"
	"testing"
)

func TestIsKind(t *testing.T) {
	testCases := []struct {
		name     string
		err      error
		kind     Kind
		expected bool
	}{
		{name: "isKind with matching kind", err: BadRequest("test"), kind: KindBadRequest, expected: true},
		{name: "isKind with non-matching kind", err: Forbidden("test"), kind: KindUnauthorized, expected: false},
		{name: "isKind with nil error", err: nil, kind: KindBadRequest, expected: false},
		{name: "wrapped with %w", err: fmt.Errorf("find user: %w", Unauthorized("test")), kind: KindUnauthorized, expected: true},
		{name: "wrapped with %w nested", err: fmt.Errorf("find user: %w", fmt.Errorf("find user: %w", Conflict("test"))), kind: KindConflict, expected: true},
		{name: "wrapped with %v loses error identity", err: fmt.Errorf("find user: %v", UnprocessableEntity("test")), kind: KindNotFound, expected: false},
		{name: "plain error", err: fmt.Errorf("plain error"), kind: KindNotFound, expected: false},
		{name: "nil error with kind", err: nil, kind: KindNotFound, expected: false},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if got := IsKind(testCase.err, testCase.kind); got != testCase.expected {
				t.Errorf("IsKind() = %v, want %v", got, testCase.expected)
			}
		})

	}
}
