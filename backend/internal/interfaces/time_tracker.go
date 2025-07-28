package interfaces

import (
	"myspace/backend/internal/types"
	"time"
)

type TimeTracker interface {
	GetUserID() string
	GetSeconds(from, to time.Time) (int, error)
	GetRunningSeconds() (int, error)
	GetMonthlyTimeByProject(dayOfMonth time.Time) (*types.ProjectTimes, error)
	GetMonthIntervals(dayOfMonth time.Time) (*types.ProjectTimes, error)
}