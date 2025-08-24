package repositories

import (
	"fmt"
	"myspace/backend/internal/config"
	"myspace/backend/internal/interfaces"
	"myspace/backend/internal/trackers"
	"myspace/backend/internal/types"
	"time"
)

type TrackersRepository struct {
	trackers []interfaces.TimeTracker
	config   *config.Config
}

func NewTrackersRepository(cfg *config.Config) *TrackersRepository {
	repo := &TrackersRepository{
		trackers: make([]interfaces.TimeTracker, 0),
		config:   cfg,
	}
	repo.hydrate()
	return repo
}

func (tr *TrackersRepository) hydrate() {
	if tr.config.Mayven.Auth != "" {
		fmt.Printf("Adding Mayven tracker with auth: %s...\n", tr.config.Mayven.Auth[:50])
		tr.addTracker(trackers.NewMayven(tr.config))
	} else {
		fmt.Println("No Mayven auth found, skipping Mayven tracker")
	}
	
	if tr.config.Everhour.Token != "" {
		tr.addTracker(trackers.NewEverhour(tr.config))
	}
	
	if tr.config.Clockify.Token != "" {
		tr.addTracker(trackers.NewClockify(tr.config))
	}
	
	fmt.Printf("Total trackers initialized: %d\n", len(tr.trackers))
}

func (tr *TrackersRepository) addTracker(tracker interfaces.TimeTracker) {
	tr.trackers = append(tr.trackers, tracker)
}

func (tr *TrackersRepository) Hours(from, to time.Time) (float64, error) {
	totalSeconds := 0
	
	for _, tracker := range tr.trackers {
		seconds, err := tracker.GetSeconds(from, to)
		if err != nil {
			continue
		}
		totalSeconds += seconds
	}
	
	return float64(totalSeconds) / 3600, nil
}

func (tr *TrackersRepository) RunningHours() (float64, error) {
	totalSeconds := 0
	
	for _, tracker := range tr.trackers {
		seconds, err := tracker.GetRunningSeconds()
		if err != nil {
			continue
		}
		totalSeconds += seconds
	}
	
	return float64(totalSeconds) / 3600, nil
}

func (tr *TrackersRepository) GetMonthlyTimeByProject(dayOfMonth time.Time) (types.ProjectTimeList, error) {
	var projectTimes types.ProjectTimeList
	
	for _, tracker := range tr.trackers {
		times, err := tracker.GetMonthlyTimeByProject(dayOfMonth)
		if err != nil {
			continue
		}
		projectTimes.Merge(times)
	}
	
	return projectTimes, nil
}

func (tr *TrackersRepository) GetDailyHours(dayOfMonth time.Time) (map[string]*float64, error) {
	var projectTimes types.ProjectTimeList
	
	for _, tracker := range tr.trackers {
		times, err := tracker.GetMonthIntervals(dayOfMonth)
		if err != nil {
			continue
		}
		projectTimes.Merge(times)
	}
	
	return projectTimes.GetDailyHours(dayOfMonth), nil
}