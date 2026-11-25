package utils

import "testing"

func TestResolveSort(t *testing.T) {
	allowed := map[string]string{
		"createdAt": "users.created_at",
		"username":  "users.username",
	}

	testCases := []struct {
		name          string
		key           string
		dir           string
		wantColumn    string
		wantDirection string
		wantErr       bool
	}{
		{
			name:          "valid key and lowercase direction",
			key:           "username",
			dir:           "asc",
			wantColumn:    "users.username",
			wantDirection: "ASC",
			wantErr:       false,
		},
		{
			name:          "invalid key and lowercase direction",
			key:           "tokenVersion",
			dir:           "asc",
			wantDirection: "ASC",
			wantErr:       true,
		},
		{
			name:          "valid key and uppercase direction",
			key:           "createdAt",
			dir:           "DESC",
			wantColumn:    "users.created_at",
			wantDirection: "DESC",
			wantErr:       false,
		},
		{
			name:    "valid key and invalid direction",
			key:     "createdAt",
			dir:     "INVALID",
			wantErr: true,
		},
		{
			name:    "invalid key and invalid direction",
			key:     "invalidKey",
			dir:     "INVALID",
			wantErr: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			column, direction, err := ResolveSort(allowed, testCase.key, testCase.dir)

			if (err != nil) != testCase.wantErr {
				t.Errorf("ResolveSort() error = %v, wantErr %v", err, testCase.wantErr)
			}

			if testCase.wantErr {
				return
			}

			if column != testCase.wantColumn || direction != testCase.wantDirection {
				t.Errorf("got (%q %q), want (%q %q)", column, direction, testCase.wantColumn, testCase.wantDirection)
			}

		})

	}
}
