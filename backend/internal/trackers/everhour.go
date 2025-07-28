package trackers

import (
	"encoding/json"
	"fmt"
	"io"
	"myspace/backend/internal/config"
	"myspace/backend/internal/types"
	"strconv"
	"time"
)

type Everhour struct {
	client *RestClient
	config *config.Config
	userID *int
}

type EverhourTimeEntry struct {
	Time      int    `json:"time"`
	CreatedAt string `json:"createdAt"`
}

type EverhourTimer struct {
	Status   string `json:"status"`
	Duration int    `json:"duration"`
}

type EverhourUser struct {
	ID int `json:"id"`
}

type EverhourProject struct {
	Name string `json:"name"`
}

func NewEverhour(cfg *config.Config) *Everhour {
	return &Everhour{
		client: NewRestClient(),
		config: cfg,
	}
}

func (e *Everhour) baseURI() string {
	return "https://api.everhour.com"
}

func (e *Everhour) headers() map[string]string {
	return map[string]string{
		"X-Api-Key": e.config.Everhour.Token,
	}
}

func (e *Everhour) GetUserID() string {
	if e.userID != nil {
		return strconv.Itoa(*e.userID)
	}
	
	resp, err := e.client.Get(e.baseURI(), "/users/me", e.headers(), nil)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	
	var user EverhourUser
	if err := json.Unmarshal(body, &user); err != nil {
		return ""
	}
	
	e.userID = &user.ID
	return strconv.Itoa(user.ID)
}

func (e *Everhour) GetSeconds(from, to time.Time) (int, error) {
	params := map[string]string{
		"from": from.Format("2006-01-02"),
		"to":   to.Format("2006-01-02"),
	}
	
	path := fmt.Sprintf("/users/%s/time", e.GetUserID())
	resp, err := e.client.Get(e.baseURI(), path, e.headers(), params)
	if err != nil {
		return 0, fmt.Errorf("failed to get time entries: %w", err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read response body: %w", err)
	}
	
	var entries []EverhourTimeEntry
	if err := json.Unmarshal(body, &entries); err != nil {
		return 0, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	
	totalSeconds := 0
	for _, entry := range entries {
		totalSeconds += entry.Time
	}
	
	return totalSeconds, nil
}

func (e *Everhour) GetRunningSeconds() (int, error) {
	resp, err := e.client.Get(e.baseURI(), "/timers/current", e.headers(), nil)
	if err != nil {
		return 0, fmt.Errorf("failed to get current timer: %w", err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read response body: %w", err)
	}
	
	var timer EverhourTimer
	if err := json.Unmarshal(body, &timer); err != nil {
		return 0, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	
	if timer.Status != "active" {
		return 0, nil
	}
	
	return timer.Duration, nil
}

func (e *Everhour) GetMonthIntervals(dayOfMonth time.Time) (*types.ProjectTimes, error) {
	som := time.Date(dayOfMonth.Year(), dayOfMonth.Month(), 1, 0, 0, 0, 0, dayOfMonth.Location())
	eom := som.AddDate(0, 1, -1)
	
	params := map[string]string{
		"from": som.Format("2006-01-02"),
		"to":   eom.Format("2006-01-02"),
	}
	
	path := fmt.Sprintf("/users/%s/time", e.GetUserID())
	resp, err := e.client.Get(e.baseURI(), path, e.headers(), params)
	if err != nil {
		return nil, fmt.Errorf("failed to get time entries: %w", err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	
	var entries []EverhourTimeEntry
	if err := json.Unmarshal(body, &entries); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	
	projectTimes := types.NewProjectTimes()
	for _, entry := range entries {
		createdAt, err := time.Parse(time.RFC3339, entry.CreatedAt)
		if err != nil {
			continue
		}
		
		projectTimes.Add(types.ProjectTime{
			Source:       "x",
			ProjectID:    "x",
			ProjectTitle: "x",
			Seconds:      entry.Time,
			Datetime:     &createdAt,
		})
	}
	
	return projectTimes, nil
}

func (e *Everhour) GetMonthlyTimeByProject(dayOfMonth time.Time) (*types.ProjectTimes, error) {
	som := time.Date(dayOfMonth.Year(), dayOfMonth.Month(), 1, 0, 0, 0, 0, dayOfMonth.Location())
	eom := som.AddDate(0, 1, -1)
	
	seconds, err := e.GetSeconds(som, eom)
	if err != nil {
		return nil, err
	}
	
	projectTimes := types.NewProjectTimes()
	if seconds > 0 {
		projectName, _ := e.getProjectName("one-and-only")
		projectTimes.Add(types.ProjectTime{
			Source:       "everhour",
			ProjectID:    "one-and-only",
			ProjectTitle: projectName,
			Seconds:      seconds,
		})
	}
	
	return projectTimes, nil
}

func (e *Everhour) getProjectName(projectID string) (string, error) {
	resp, err := e.client.Get(e.baseURI(), "/projects", e.headers(), nil)
	if err != nil {
		return "", fmt.Errorf("failed to get projects: %w", err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}
	
	var projects []EverhourProject
	if err := json.Unmarshal(body, &projects); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}
	
	if len(projects) > 0 {
		return projects[0].Name, nil
	}
	
	return "Unknown Project", nil
}