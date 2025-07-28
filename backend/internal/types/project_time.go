package types

import (
	"math"
	"time"
)

type ProjectTime struct {
	Source       string     `json:"source"`
	ProjectID    string     `json:"project_id"`
	ProjectTitle string     `json:"project_title"`
	Seconds      int        `json:"seconds"`
	Datetime     *time.Time `json:"datetime,omitempty"`
}

func (pt *ProjectTime) GetHours() float64 {
	hours := float64(pt.Seconds) / 3600
	return math.Round(hours*100) / 100
}

func (pt *ProjectTime) ToMap() map[string]interface{} {
	result := map[string]interface{}{
		"source":        pt.Source,
		"project_id":    pt.ProjectID,
		"project_title": pt.ProjectTitle,
		"seconds":       pt.Seconds,
		"hours":         pt.GetHours(),
	}
	
	if pt.Datetime != nil {
		result["datetime"] = pt.Datetime
	}
	
	return result
}