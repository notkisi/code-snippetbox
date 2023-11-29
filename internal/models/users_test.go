package models

import (
	"testing"

	"github.com/notkisi/snippetbox/internal/assert"
)

func TestUserModelExists(t *testing.T) {
	tests := []struct {
		name     string
		userID   int
		expected bool
	}{
		{
			name:     "Valid ID",
			userID:   1,
			expected: true,
		},
		{
			name:     "Zero ID",
			userID:   0,
			expected: false,
		},
		{
			name:     "Non existent ID",
			userID:   123,
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db := newTestDB(t)

			m := UserModel{db}

			exists, err := m.Exists(tc.userID)
			assert.Equal(t, exists, tc.expected)
			// verify that there was no errors
			assert.NilError(t, err)
		})
	}
}
