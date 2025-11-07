package mysql

import "time"

func timeNowUTC() time.Time {
	return time.Now().UTC()
}
