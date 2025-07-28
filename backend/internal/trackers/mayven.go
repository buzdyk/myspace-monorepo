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

type Mayven struct {
	client *RestClient
	config *config.Config
	userID *int
}

type MayvenTimeStats struct {
	Data struct {
		ChartData           []MayvenChartData       `json:"chartData"`
		AggregatedIntervals []MayvenAggregatedData  `json:"aggregatedIntervals"`
	} `json:"data"`
}

type MayvenChartData struct {
	Seconds int    `json:"seconds"`
	Date    string `json:"_date"`
}

type MayvenAggregatedData struct {
	ItemID  string `json:"item_id"`
	Title   string `json:"title"`
	Seconds int    `json:"seconds"`
}

type MayvenHydrate struct {
	Data struct {
		Me struct {
			Data struct {
				ID int `json:"id"`
			} `json:"data"`
		} `json:"me"`
	} `json:"data"`
}

type MayvenTimer struct {
	Data struct {
		StartedAt string `json:"started_at"`
	} `json:"data"`
}

func NewMayven(cfg *config.Config) *Mayven {
	return &Mayven{
		client: NewRestClient(),
		config: cfg,
	}
}

func (m *Mayven) baseURI() string {
	return "https://api.mayven.io"
}

func (m *Mayven) headers() map[string]string {
	return map[string]string{
		"Authorization": m.config.Mayven.Auth,
		"Accept":        "application/json",
	}
}

func (m *Mayven) GetUserID() string {
	if m.userID != nil {
		return strconv.Itoa(*m.userID)
	}
	
	resp, err := m.client.Get(m.baseURI(), "/api/hydrate", m.headers(), nil)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	
	var hydrate MayvenHydrate
	if err := json.Unmarshal(body, &hydrate); err != nil {
		return ""
	}
	
	m.userID = &hydrate.Data.Me.Data.ID
	return strconv.Itoa(*m.userID)
}

func (m *Mayven) GetSeconds(from, to time.Time) (int, error) {
	params := map[string]string{
		"dateStart": from.Format("2006-01-02") + " 00:00:00",
		"dateEnd":   to.Format("2006-01-02") + " 23:59:59",
		"users[]":   m.GetUserID(),
	}
	
	resp, err := m.client.Get(m.baseURI(), "/api/time-statistics", m.headers(), params)
	if err != nil {
		return 0, fmt.Errorf("failed to get time statistics: %w", err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read response body: %w", err)
	}
	
	var stats MayvenTimeStats
	if err := json.Unmarshal(body, &stats); err != nil {
		return 0, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	
	totalSeconds := 0
	for _, item := range stats.Data.ChartData {
		totalSeconds += item.Seconds
	}
	
	return totalSeconds, nil
}

func (m *Mayven) GetRunningSeconds() (int, error) {
	resp, err := m.client.Get(m.baseURI(), "/api/timer", m.headers(), nil)
	if err != nil {
		return 0, fmt.Errorf("failed to get timer: %w", err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read response body: %w", err)
	}
	
	var timer MayvenTimer
	if err := json.Unmarshal(body, &timer); err != nil {
		return 0, nil
	}
	
	if timer.Data.StartedAt == "" {
		return 0, nil
	}
	
	startedAt, err := time.Parse(time.RFC3339, timer.Data.StartedAt)
	if err != nil {
		return 0, nil
	}
	
	return int(time.Since(startedAt).Seconds()), nil
}

func (m *Mayven) GetMonthIntervals(dayOfMonth time.Time) (*types.ProjectTimes, error) {
	som := time.Date(dayOfMonth.Year(), dayOfMonth.Month(), 1, 0, 0, 0, 0, dayOfMonth.Location())
	eom := som.AddDate(0, 1, -1)
	
	params := map[string]string{
		"dateStart": som.Format("2006-01-02") + " 00:00:00",
		"dateEnd":   eom.Format("2006-01-02") + " 23:59:59",
		"users[]":   m.GetUserID(),
	}
	
	resp, err := m.client.Get(m.baseURI(), "/api/time-statistics", m.headers(), params)
	if err != nil {
		return nil, fmt.Errorf("failed to get time statistics: %w", err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	
	var stats MayvenTimeStats
	if err := json.Unmarshal(body, &stats); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	
	projectTimes := types.NewProjectTimes()
	for _, item := range stats.Data.ChartData {
		date, err := time.Parse("2006-01-02", item.Date)
		if err != nil {
			continue
		}
		
		projectTimes.Add(types.ProjectTime{
			Source:       "x",
			ProjectID:    "x",
			ProjectTitle: "x",
			Seconds:      item.Seconds,
			Datetime:     &date,
		})
	}
	
	return projectTimes, nil
}

func (m *Mayven) GetMonthlyTimeByProject(dayOfMonth time.Time) (*types.ProjectTimes, error) {
	som := time.Date(dayOfMonth.Year(), dayOfMonth.Month(), 1, 0, 0, 0, 0, dayOfMonth.Location())
	eom := som.AddDate(0, 1, -1)
	
	params := map[string]string{
		"dateStart":            som.Format("2006-01-02") + " 00:00:00",
		"dateEnd":              eom.Format("2006-01-02") + " 23:59:59",
		"users[]":              m.GetUserID(),
		"groupByPrimaryValue":  "project_id",
		"groupBySecondValue":   "todo_id",
		"groupBy":              "project_id",
		"orderBy":              "seconds:desc",
	}
	
	resp, err := m.client.Get(m.baseURI(), "/api/time-statistics", m.headers(), params)
	if err != nil {
		return nil, fmt.Errorf("failed to get time statistics: %w", err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	
	var stats MayvenTimeStats
	if err := json.Unmarshal(body, &stats); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	
	projectTimes := types.NewProjectTimes()
	for _, item := range stats.Data.AggregatedIntervals {
		projectTimes.Add(types.ProjectTime{
			Source:       "mayven",
			ProjectID:    item.ItemID,
			ProjectTitle: item.Title,
			Seconds:      item.Seconds,
		})
	}
	
	return projectTimes, nil
}