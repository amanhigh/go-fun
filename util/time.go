package util

import "time"

/* Layouts */
const PRINT_LAYOUT = "Jan 2, 2006 at 3:04pm (MST)"
const SLASH_MILLISECOND_LAYOUT = "[02/Jan/2006:15:04:05.000]"

func FormatTime(time time.Time, layout string) string {
	return time.Format(layout)
}

func TimeAgo(duration time.Duration) time.Time {
	return time.Now().Add(-duration)
}

func TimeAfter(duration time.Duration) time.Time {
	return time.Now().Add(duration)
}

/**
Using Hour String Compute Todays Data
"4:05AM"
 */
func TodaysHour(string string) time.Time {
	parsedHourMin, _ := time.Parse(time.Kitchen, string)
	year, month, day := time.Now().Date()
	return time.Date(year, month, day, parsedHourMin.Hour(), parsedHourMin.Minute(), 0, 0, time.UTC)
}
