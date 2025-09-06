package num

import (
	"fmt"
	"testing"
	"time"
)

func TestParseFreq(t *testing.T) {
	testCases := []string{"2H", "30T", "1D", "5S", "10min", "15ms", "1000U", "2D"}

	for _, tc := range testCases {
		dur, err := ParseFreq(tc)
		if err != nil {
			fmt.Printf("Parse %s → ERROR: %v\n", tc, err)
		} else {
			fmt.Printf("Parse %s → %v\n", tc, dur)
		}
	}
}

func TestDateRange(t *testing.T) {
	times, _ := DateRange(time.Now(), 5, "2H")
	for _, t := range times {
		fmt.Println(t.Format("2006-01-02 15:04:05"))
	}
}
