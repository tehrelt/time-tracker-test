package dto

import "time"

type SaveActivity struct {
	UserId    string
	StartTime time.Time
}

type StopActivityDto struct {
	UserId  string
	EndTime time.Time
}
