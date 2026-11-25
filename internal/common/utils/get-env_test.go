package utils

import (
	"testing"
)

func TestGetEnv(t *testing.T) {
	t.Setenv("APP_ENV", "development") // Set an environment variable for testing

	testCases := []struct {
		name         string
		key          string
		defaultValue string
		expected     string
	}{
		{
			name:         "existing environment variable",
			key:          "APP_ENV",
			defaultValue: "default_value",
			expected:     "development",
		},
		{
			name:         "non-existing environment variable",
			key:          "NON_EXISTING_ENV",
			defaultValue: "default_value",
			expected:     "default_value",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			value := GetEnv(testCase.key, testCase.defaultValue)

			if value != testCase.expected {
				t.Errorf("GetEnv(%q, %q); want %q", testCase.key, testCase.defaultValue, testCase.expected)
			}
		})

	}

}
