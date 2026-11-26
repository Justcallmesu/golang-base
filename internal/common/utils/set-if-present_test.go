package utils

import "testing"

func TestSetIfPresent(t *testing.T) {
	testCases := []struct {
		name          string
		value         any
		valueExpected bool
	}{
		{
			name:          "source nil, target should remain unchanged",
			value:         nil,
			valueExpected: false,
		},
		{
			name:          "source not nil in function format, target should be updated",
			value:         func() *string { s := "new value"; return &s },
			valueExpected: true,
		},
		{
			name:          "source not nil in bool format, target should be updated",
			value:         true,
			valueExpected: true,
		},
		{
			name:          "source not nil in string format, target should be updated",
			value:         "new value",
			valueExpected: true,
		},
	}

	for _, testCase := range testCases {
		var target any

		SetIfPresent(&target, &testCase.value)

		if (target != nil) != testCase.valueExpected {
			t.Fatalf("Test case '%s' failed: expected target to be %v, but got %v", testCase.name, testCase.valueExpected, target != nil)
		}
	}
}
