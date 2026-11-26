package models

import (
	"errors"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestComparePassword(t *testing.T) {
	user := &User{}

	if err := user.SetPassword("correct-dummy-password"); err != nil {
		t.Fatalf("SetPassword() error = %v", err)
	}

	testCases := []struct {
		name          string
		candidate     string
		expectedError error
	}{
		{
			name:          "compare correct password",
			candidate:     "correct-dummy-password",
			expectedError: nil,
		},
		{
			name:          "compare incorrect password",
			candidate:     "incorrect-dummy-password",
			expectedError: bcrypt.ErrMismatchedHashAndPassword,
		},
		{
			name:          "compare empty password",
			candidate:     "",
			expectedError: bcrypt.ErrMismatchedHashAndPassword,
		},
		{
			name:          "compare correct case password",
			candidate:     "Correct-Dummy-Password",
			expectedError: bcrypt.ErrMismatchedHashAndPassword,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			err := user.ComparePassword(testCase.candidate)

			if !errors.Is(err, testCase.expectedError) {
				t.Errorf("ComparePassword() error = %v, wantErr %v", err, testCase.expectedError)
			}

		})
	}
}

func TestSetPassword(t *testing.T) {

	testCases := []struct {
		name        string
		password    string
		expectError bool
	}{
		{
			name:        "set valid password",
			password:    "valid-password",
			expectError: false,
		},
		{
			name:        "set empty password",
			password:    "",
			expectError: true,
		},
		{
			name:        "set long password",
			password:    "this-is-a-very-long-password-that-should-still-be-accepted-by-the-set-password-function",
			expectError: true,
		},
	}

	for _, testCase := range testCases {
		user := &User{}

		t.Run(testCase.name, func(t *testing.T) {
			err := user.SetPassword(testCase.password)

			if (err != nil) != testCase.expectError {
				t.Errorf("SetPassword() error = %v, expectError %v", err, testCase.expectError)
			}
		})
	}
}

func TestCompareDummyPassword(t *testing.T) {
	cost, err := bcrypt.Cost(dummyPasswordHash)

	if err != nil {
		t.Fatalf("Failed to get cost of dummy password hash: %v", err)
	}

	if cost != passwordHashCost {
		t.Errorf("dummy cost = %d, want %d (timing would differ from real users)", cost, passwordHashCost)
	}
}
