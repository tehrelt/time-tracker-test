package filters

import "time"

type Activity struct {
	UserId    string
	StartTime *time.Time
	EndTime   *time.Time
}
