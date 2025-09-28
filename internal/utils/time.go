package utils

import "time"

func MinutesToNanoseconds(minutes uint) time.Duration {
	return time.Duration(minutes) * time.Minute
}
