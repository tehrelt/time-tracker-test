package domain

import "time"

type ActivityRecord struct {
	Id        int64      `json:"id" db:"id"`
	User      User       `json:"user" db:"user"`
	StartTime time.Time  `json:"startTime" db:"start_time"`
	EndTime   *time.Time `json:"endTime,omitempty" db:"end_time"`
}

type Session struct {
	StartTime time.Time  `json:"startTime" db:"start_time"`
	EndTime   *time.Time `json:"endTime,omitempty" db:"end_time"`
}

type ActivitySummary struct {
	UserId      string        `json:"userId" db:"user_id"`
	IsActiveNow bool          `json:"isActiveNow" db:"is_active_now"`
	Sessions    []*Session    `json:"sessions"`
	TotalTime   time.Duration `json:"totalTime" db:"total_time"`
	TotalCount  int           `json:"totalCount" db:"total_count"`
}
