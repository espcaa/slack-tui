package utils

import (
	"strconv"
	"time"
)

func TimestampToString(ts string) string {
	// Convert timestamp to a human-readable format
	if len(ts) < 10 {
		return ts // Not a valid timestamp
	}
	// Convert the first 10 characters to an integer
	unixTime, err := strconv.ParseInt(ts[:10], 10, 64)
	if err != nil {
		return ts // Return original if conversion fails
	}
	return time.Unix(unixTime, 0).Format("01/02/26 15:04")
}
