package types

import (
	"math"
	"time"
)

type ProjectTimes struct {
	items        []ProjectTime
	currentIndex int
}

func NewProjectTimes() *ProjectTimes {
	return &ProjectTimes{
		items:        make([]ProjectTime, 0),
		currentIndex: 0,
	}
}

func (pts *ProjectTimes) Add(projectTime ProjectTime) *ProjectTimes {
	pts.items = append(pts.items, projectTime)
	return pts
}

func (pts *ProjectTimes) Merge(other *ProjectTimes) *ProjectTimes {
	other.Rewind()
	for {
		item := other.Next()
		if item == nil {
			break
		}
		pts.Add(*item)
	}
	return pts
}

func (pts *ProjectTimes) Rewind() {
	pts.currentIndex = 0
}

func (pts *ProjectTimes) Next() *ProjectTime {
	if pts.currentIndex >= len(pts.items) {
		return nil
	}
	item := &pts.items[pts.currentIndex]
	pts.currentIndex++
	return item
}

func (pts *ProjectTimes) Count() int {
	return len(pts.items)
}

func (pts *ProjectTimes) GetHours() float64 {
	total := 0.0
	for _, item := range pts.items {
		total += item.GetHours()
	}
	return math.Round(total*100) / 100
}

func (pts *ProjectTimes) GetDailyHours(dayOfMonth time.Time) map[string]*float64 {
	som := time.Date(dayOfMonth.Year(), dayOfMonth.Month(), 1, 0, 0, 0, 0, dayOfMonth.Location())
	eom := som.AddDate(0, 1, -1)
	
	days := make(map[string]*float64)
	
	for d := som; !d.After(eom); d = d.AddDate(0, 0, 1) {
		days[d.Format("2006-01-02")] = nil
	}
	
	for _, item := range pts.items {
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

func (pts *ProjectTimes) ToArray() []map[string]interface{} {
	result := make([]map[string]interface{}, len(pts.items))
	for i, item := range pts.items {
		result[i] = item.ToMap()
	}
	return result
}