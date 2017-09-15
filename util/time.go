package util

import "time"

/* Date/Time */
const PRINT_LAYOUT = "Jan 2, 2006 at 3:04pm (MST)"

func FormatTime(time time.Time, layout string) string {
	return time.Format(layout)
}
