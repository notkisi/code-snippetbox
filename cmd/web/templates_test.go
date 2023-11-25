package main

import (
	"testing"
	"time"

	"github.com/notkisi/snippetbox/internal/assert"
)

func TestHumanDate(t *testing.T) {
	tests := []struct {
		name     string
		param    time.Time
		expected string
	}{
		{
			name:     "UTC",
			param:    time.Date(2022, 3, 17, 10, 15, 0, 0, time.UTC),
			expected: "17 Mar 2022 at 10:15",
		},
		{
			name:     "Empty",
			param:    time.Time{},
			expected: "",
		},
		{
			name:     "CET",
			param:    time.Date(2022, 3, 17, 10, 15, 0, 0, time.FixedZone("CET", 1*60*60)),
			expected: "17 Mar 2022 at 09:15",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			actual := humanDate(testCase.param)
			assert.Equal(t, actual, testCase.expected)
		})
	}

}
