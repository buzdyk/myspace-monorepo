package trackers

import (
	"encoding/json"
	"fmt"
	"io"
	"myspace/backend/internal/config"
	"myspace/backend/internal/types"
	"strings"
	"time"
)

type Clockify struct {
	client *RestClient
	config *config.Config
}

type ClockifyTimeEntry struct {
	TimeInterval struct {
		Start string `json:"start"`
		End   string `json:"end"`
	} `json:"timeInterval"`
}

func NewClockify(cfg *config.Config) *Clockify {
	return &Clockify{
		client: NewRestClient(),
		config: cfg,
	}
}

func (c *Clockify) baseURI() string {
	return "https://api.clockify.me"
}

func (c *Clockify) headers() map[string]string {
	return map[string]string{
		"x-api-key": c.config.Clockify.Token,
		"Accept":    "application/json",
	}
}

func (c *Clockify) getPathWithWorkspace(path string) string {
	if path == "" {
		path = "/"
	}
	return fmt.Sprintf("/api/v1/workspaces/%s%s", c.config.Clockify.WorkspaceID, path)
}

func (c *Clockify) getSecondsForTimeEntry(entry ClockifyTimeEntry) (int, error) {
	start, err := time.Parse(time.RFC3339, entry.TimeInterval.Start)
	if err != nil {
		return 0, fmt.Errorf("failed to parse start time: %w", err)
	}
	
	end, err := time.Parse(time.RFC3339, entry.TimeInterval.End)
	if err != nil {
		return 0, fmt.Errorf("failed to parse end time: %w", err)
	}
	
	return int(end.Sub(start).Seconds()), nil
}

func (c *Clockify) GetUserID() string {
	return c.config.Clockify.UserID
}

func (c *Clockify) GetSeconds(from, to time.Time) (int, error) {
	path := c.getPathWithWorkspace(fmt.Sprintf("/user/%s/time-entries", c.GetUserID()))
	params := map[string]string{
		"start": from.Format("2006-01-02") + "T00:00:00Z",
		"end":   to.Format("2006-01-02") + "T23:59:59Z",
	}
	
	resp, err := c.client.Get(c.baseURI(), path, c.headers(), params)
	if err != nil {
		return 0, fmt.Errorf("failed to get time entries: %w", err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read response body: %w", err)
	}
	
	var entries []ClockifyTimeEntry
	if err := json.Unmarshal(body, &entries); err != nil {
		return 0, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	
	totalSeconds := 0
	for _, entry := range entries {
		seconds, err := c.getSecondsForTimeEntry(entry)
		if err != nil {
			continue
		}
		totalSeconds += seconds
	}
	
	return totalSeconds, nil
}

func (c *Clockify) GetRunningSeconds() (int, error) {
	return 0, nil
}

func (c *Clockify) GetMonthlyTimeByProject(dayOfMonth time.Time) (types.ProjectTimeList, error) {
	som := time.Date(dayOfMonth.Year(), dayOfMonth.Month(), 1, 0, 0, 0, 0, dayOfMonth.Location())
	eom := som.AddDate(0, 1, -1)
	
	seconds, err := c.GetSeconds(som, eom)
	if err != nil {
		return nil, err
	}
	
	var projectTimes types.ProjectTimeList
	projectTimes.Add(types.ProjectTime{
		Source:       "clockify",
		ProjectID:    "x",
		ProjectTitle: c.generateRandomString(10),
		Seconds:      seconds,
	})
	
	return projectTimes, nil
}

func (c *Clockify) GetMonthIntervals(dayOfMonth time.Time) (types.ProjectTimeList, error) {
	som := time.Date(dayOfMonth.Year(), dayOfMonth.Month(), 1, 0, 0, 0, 0, dayOfMonth.Location())
	eom := som.AddDate(0, 1, -1)
	
	path := c.getPathWithWorkspace(fmt.Sprintf("/user/%s/time-entries", c.GetUserID()))
	params := map[string]string{
		"start": som.Format("2006-01-02") + "T00:00:00Z",
		"end":   eom.Format("2006-01-02") + "T23:59:59Z",
	}
	
	resp, err := c.client.Get(c.baseURI(), path, c.headers(), params)
	if err != nil {
		return nil, fmt.Errorf("failed to get time entries: %w", err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	
	var entries []ClockifyTimeEntry
	if err := json.Unmarshal(body, &entries); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	
	var projectTimes types.ProjectTimeList
	for _, entry := range entries {
		seconds, err := c.getSecondsForTimeEntry(entry)
		if err != nil {
			continue
		}
		
		start, err := time.Parse(time.RFC3339, entry.TimeInterval.Start)
		if err != nil {
			continue
		}
		
		projectTimes.Add(types.ProjectTime{
			Source:       "x",
			ProjectID:    "x", 
			ProjectTitle: "x",
			Seconds:      seconds,
			Datetime:     &start,
		})
	}
	
	return projectTimes, nil
}

func (c *Clockify) generateRandomString(length int) string {
	return strings.Repeat("x", length)
}