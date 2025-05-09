package security

import "time"

const nanosInSecond = 1e9

func timeToFloat64(t time.Time) float64 {
	return float64(t.UTC().UnixNano()) / nanosInSecond
}
