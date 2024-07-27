package pkg

import "time"

func SecondsToTime(seconds int64) time.Time {
	duration := time.Duration(seconds) * time.Second
	return time.Unix(0, 0).Add(duration)
}
