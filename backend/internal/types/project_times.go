package types

import (
	"math"
	"time"
)

// ProjectTimeList is a slice of ProjectTime, using idiomatic Go patterns
type ProjectTimeList []ProjectTime

// GetHours returns the total hours for all project times
func (ptl ProjectTimeList) GetHours() float64 {
	total := 0.0
	for _, item := range ptl {
		total += item.GetHours()
	}
	return math.Round(total*100) / 100
}

// GetDailyHours returns a map of daily hours for the given month
func (ptl ProjectTimeList) GetDailyHours(dayOfMonth time.Time) map[string]*float64 {
	som := time.Date(dayOfMonth.Year(), dayOfMonth.Month(), 1, 0, 0, 0, 0, dayOfMonth.Location())
	eom := som.AddDate(0, 1, -1)

	days := make(map[string]*float64)

	// Initialize all days of the month with nil
	for d := som; !d.After(eom); d = d.AddDate(0, 0, 1) {
		days[d.Format("2006-01-02")] = nil
	}

	// Sum up hours for each day
	for _, item := range ptl {
		if item.Datetime != nil {
			day := item.Datetime.Format("2006-01-02")
			if days[day] == nil {
				hours := 0.0
				days[day] = &hours
			}
			*days[day] += item.GetHours()
		}
	}

	return days
}

// ToArray converts the list to an array of maps (for JSON serialization)
func (ptl ProjectTimeList) ToArray() []map[string]interface{} {
	result := make([]map[string]interface{}, len(ptl))
	for i, item := range ptl {
		result[i] = item.ToMap()
	}
	return result
}

// Merge appends all items from another ProjectTimeList
func (ptl *ProjectTimeList) Merge(other ProjectTimeList) {
	*ptl = append(*ptl, other...)
}

// Add appends a ProjectTime to the list
func (ptl *ProjectTimeList) Add(projectTime ProjectTime) {
	*ptl = append(*ptl, projectTime)
}
