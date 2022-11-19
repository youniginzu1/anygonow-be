package lib

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetTimeRange(t *testing.T) {
	TEST_CASE := []map[string]interface{}{
		{
			"quantity": 1,
			"unit":     "WEEK",
			"expected": 1 * 7 * 24 * time.Hour,
		},
		{
			"quantity": 2,
			"unit":     "WEEK",
			"expected": 2 * 7 * 24 * time.Hour,
		},
		{
			"quantity": 1,
			"unit":     "DAY",
			"expected": 1 * 24 * time.Hour,
		},
		{
			"quantity": 3,
			"unit":     "HOUR",
			"expected": 3 * time.Hour,
		},
	}

	for _, c := range TEST_CASE {
		timeRange := GetTimeRange(c["quantity"].(int), c["unit"].(string))
		assert.Equal(t, c["expected"], timeRange,
			"actual %v, expected %v, case quantity %d unit %s", c["expected"], timeRange, c["quantity"].(int), c["unit"].(string))
	}
}
