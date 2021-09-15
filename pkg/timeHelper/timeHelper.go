package timeHelper

import "time"

// TimeIsDayBeforeTime will return true if the time provided was of a different day prior to today.
// Example:
//		- TimeIsDayBeforeTime(2021-09-11T09:10, 2021-09-12T09:15) returns true
//		- TimeIsDayBeforeTime(2021-09-09T09:10, 2021-09-12T09:15) returns true
//		- TimeIsDayBeforeTime(2021-09-12T09:10, 2021-09-12T09:15) returns false
//		- TimeIsDayBeforeTime(2021-09-12T09:30, 2021-09-12T09:15) returns false
func TimeIsDayBeforeTime(time1, time2 time.Time) bool {
	time2Midnight := time.Date(time2.Year(), time2.Month(), time2.Day(), 0, 0, 0, 0, time2.Location())
	return time1.Before(time2Midnight)
}
