package jwt

import "time"

// TokenExpiry helper to convert duration to time.Duration
func TokenExpiry(value int64, unit string) time.Duration {
	switch unit {
	case "minutes":
		return time.Duration(value) * time.Minute
	case "hours":
		return time.Duration(value) * time.Hour
	case "days":
		return time.Duration(value) * 24 * time.Hour
	default:
		return time.Duration(value) * time.Minute
	}
}
