package num

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func ParseFreq(freq string) (time.Duration, error) {
	freq = strings.TrimSpace(freq)
	if len(freq) == 0 {
		return 0, errors.New("empty freq string")
	}

	// 解析数字部分
	i := 0
	for i < len(freq) && (freq[i] >= '0' && freq[i] <= '9') {
		i++
	}

	var n int
	var err error
	if i == 0 {
		n = 1 // 默认倍数为1
	} else {
		n, err = strconv.Atoi(freq[:i])
		if err != nil {
			return 0, fmt.Errorf("invalid number in freq: %w", err)
		}
	}

	unit := freq[i:]
	if unit == "" {
		return 0, errors.New("missing unit in freq")
	}

	// 映射 Pandas 单位到 time.Duration
	switch unit {
	case "N", "ns":
		return time.Duration(n), nil
	case "U", "us", "µs":
		return time.Microsecond * time.Duration(n), nil
	case "L", "ms":
		return time.Millisecond * time.Duration(n), nil
	case "S":
		return time.Second * time.Duration(n), nil
	case "T", "min":
		return time.Minute * time.Duration(n), nil
	case "H":
		return time.Hour * time.Duration(n), nil
	case "D":
		return 24 * time.Hour * time.Duration(n), nil
	default:
		return 0, fmt.Errorf("unsupported freq unit: %s", unit)
	}
}

// DateRange 生成时间序列
func DateRange(start time.Time, periods int, freqStr string) ([]time.Time, error) {
	dur, err := ParseFreq(freqStr)
	if err != nil {
		return nil, err
	}

	var result []time.Time
	t := start
	for i := 0; i < periods; i++ {
		result = append(result, t)
		t = t.Add(dur)
	}
	return result, nil
}
