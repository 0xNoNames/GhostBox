package utils

import (
	"fmt"
	"time"
)

const (
	_  = iota             // ignore first value by assigning to blank identifier
	KB = 1 << (10 * iota) // 1024
	MB
	GB
	TB
)

// byteSuffixes returns the appropriate suffix and unit based on the provided number of bytes.
func byteSuffixes(i int64) (suffix string, unit float64) {
	switch {
	case i < KB:
		suffix = "B"
		unit = 1
	case i < MB:
		suffix = "KB"
		unit = KB
	case i < GB:
		suffix = "MB"
		unit = MB
	case i < TB:
		suffix = "GB"
		unit = GB
	default:
		suffix = "TB"
		unit = TB
	}
	return
}

// FormatBytesProgress formats the bytes completed and total length into a string representation with the appropriate suffix.
func FormatBytesProgress(bytesCompleted, totalLength int64) string {
	suffix, unit := byteSuffixes(totalLength)
	return fmt.Sprintf("%.1f/%.1f%s",
		float64(bytesCompleted)/unit,
		float64(totalLength)/unit,
		suffix)
}

// FormatDuration formats a duration in a human-readable format.
func FormatDuration(duration time.Duration) string {
	duration *= time.Second
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60

	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}
