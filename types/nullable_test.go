package types

import (
	"database/sql/driver"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Nullable_Scan(t *testing.T) {
	num := 800
	tests := []struct {
		input    driver.Value
		name     string
		expected *Nullable[int]
	}{
		{
			name:  "Valid Int",
			input: 800,
			expected: &Nullable[int]{
				Set:   true,
				Valid: true,
				Val:   &num,
			},
		},
		{
			name:  "Null",
			input: nil,
			expected: &Nullable[int]{
				Set:   true,
				Valid: false,
				Val:   nil,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			json := &Nullable[int]{}
			json.Scan(tc.input)
			assert.Equal(t, tc.expected, json)
		})
	}
}
