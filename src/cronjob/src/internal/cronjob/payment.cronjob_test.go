package cronjob

import (
	"testing"
	"time"
)

var Layout string = "2006-01-02 15:04:05"

func mustParseTime(layout string, value string) time.Time {
	s, err := time.Parse(layout, value)
	if err != nil {
		panic(err)
	}
	return s
}

func TestTime(t *testing.T) {
	TEST_CASE := []map[string]interface{}{
		{
			"day":           "Sunday",
			"unit":          "DAY",
			"now":           mustParseTime(Layout, "2022-04-26 15:25:00"),
			"expected_time": mustParseTime(Layout, "2022-04-27 00:00:00"),
		},
		{
			"day":           "Sunday",
			"unit":          "WEEK",
			"now":           mustParseTime(Layout, "2022-04-26 15:25:00"),
			"expected_time": mustParseTime(Layout, "2022-05-01 00:00:00"),
		},
		{
			"day":           "Sunday",
			"unit":          "HOUR",
			"now":           mustParseTime(Layout, "2022-04-26 15:25:00"),
			"expected_time": mustParseTime(Layout, "2022-04-26 16:00:00"),
		},
		{
			"day":           "Saturday",
			"unit":          "WEEK",
			"now":           mustParseTime(Layout, "2022-04-26 15:25:00"),
			"expected_time": mustParseTime(Layout, "2022-04-30 00:00:00"),
		},
	}

	for _, c := range TEST_CASE {
		timee := externalGetNextPeriod(c["day"].(string), c["unit"].(string), c["now"].(time.Time))
		if !timee.Equal(c["expected_time"].(time.Time)) {
			t.Errorf("not equal, actual: %v ,expected: %v, day: %s, unit: %s", timee, c["expected_time"], c["day"], c["unit"])
		}
	}
}

func ThrottleConcurrency(limit int, arr []string, handle func(string)) <-chan struct{} {
	result := make(chan struct{}, len(arr))
	go func() {
		blocking := make(chan struct{}, limit)
		for _, v := range arr {
			blocking <- struct{}{}
			go func(s string) {
				handle(s)
				result <- <-blocking
			}(v)
		}
	}()
	return result
}

func TestThrottle(t *testing.T) {
	arr := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
	outChan := ThrottleConcurrency(3, arr, func(s string) {
		t.Log(s)
		time.Sleep(1 * time.Second)
	})
	for i := 0; i < len(arr); i++ {
		<-outChan
	}
}
