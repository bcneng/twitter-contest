package contest

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestDo(t *testing.T) {
	tests := []struct {
		name           string
		candidates     []string
		pick           int
		expectedPicked int
		expectedList   []string
	}{
		{
			name:           "Happy path",
			candidates:     []string{"@sdecandelario", "@Linkita", "@gonzaloserrano", "@koesystems", "@smoyac"},
			pick:           2,
			expectedPicked: 2,
		},
		{
			name:           "Candidates list smaller than winners to pick",
			candidates:     []string{"@sdecandelario", "@Linkita", "@gonzaloserrano", "@koesystems", "@smoyac"},
			pick:           10,
			expectedPicked: 5,
			expectedList:   []string{"@sdecandelario", "@Linkita", "@gonzaloserrano", "@koesystems", "@smoyac"},
		},
		{
			name:           "0 to pick means 0 picked",
			candidates:     []string{"@sdecandelario", "@Linkita", "@gonzaloserrano", "@koesystems", "@smoyac"},
			pick:           0,
			expectedPicked: 0,
		},
		{
			name:           "contest on an empty candidates list",
			candidates:     []string{},
			pick:           5,
			expectedPicked: 0,
		},
	}
	Version = "test"
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := Do(test.candidates, test.pick)
			if test.expectedPicked == 0 {
				require.Nil(t, result)
				return
			} else {
				require.NotNil(t, result)
			}
			assert.Len(t, result.Winners, test.expectedPicked)

			if len(test.expectedList) > 0 {
				assert.ElementsMatch(t, result.Winners, test.expectedList)
			}

			assert.Equal(t, "test", result.Version)
			assert.NotEmpty(t, result.Time)
		})
	}
}
