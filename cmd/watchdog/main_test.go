package main

import (
	"testing"
	"time"
)

func TestGetCheckInterval(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		expected time.Duration
	}{
		{
			name:     "uses default when environment variable is empty",
			envValue: "",
			expected: 30 * time.Second,
		},
		{
			name:     "uses configured interval",
			envValue: "5",
			expected: 5 * time.Second,
		},
		{
			name:     "uses default for invalid string",
			envValue: "abc",
			expected: 30 * time.Second,
		},
		{
			name:     "uses default for zero",
			envValue: "0",
			expected: 30 * time.Second,
		},
		{
			name:     "uses default for negative value",
			envValue: "-10",
			expected: 30 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("CHECK_INTERVAL", tt.envValue)

			actual := getCheckInterval()

			if actual != tt.expected {
				t.Fatalf(
					"getCheckInterval() = %v; expected %v",
					actual,
					tt.expected,
				)
			}
		})
	}
}
