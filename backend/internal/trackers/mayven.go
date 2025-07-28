package trackers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
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
	Seconds int    `json:"-"`
	Date    string `json:"_date"`
	SecondsStr string `json:"seconds"`
}

func (m *MayvenChartData) UnmarshalJSON(data []byte) error {
	type Alias MayvenChartData
	aux := &struct {
		SecondsStr string `json:"seconds"`
		*Alias
	}{
		Alias: (*Alias)(m),
	}
	
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	
	seconds, err := strconv.Atoi(aux.SecondsStr)
	if err != nil {
		return err
	}
	
	m.Seconds = seconds
	return nil
}

type MayvenAggregatedData struct {
	ItemID  string `json:"-"`
	Title   string `json:"title"`
	Seconds int    `json:"-"`
}

func (m *MayvenAggregatedData) UnmarshalJSON(data []byte) error {
	aux := struct {
		ItemIDInt  int    `json:"item_id"`
		ItemIDStr  string `json:"item_id"`
		Title      string `json:"title"`
		SecondsStr string `json:"seconds"`
	}{}
	
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	
	// Handle item_id as either string or int
	if aux.ItemIDStr != "" {
		m.ItemID = aux.ItemIDStr
	} else {
		m.ItemID = strconv.Itoa(aux.ItemIDInt)
	}
	
	m.Title = aux.Title
	
	seconds, err := strconv.Atoi(aux.SecondsStr)
	if err != nil {
		return err
	}
	
	m.Seconds = seconds
	return nil
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
	return m.config.Mayven.ApiURL
}

func (m *Mayven) headers() map[string]string {
	return map[string]string{
		"Authorization": m.config.Mayven.Auth,
		"Accept":        "application/json",
	}
}

func (m *Mayven) GetUserID() string {
	if m.userID != nil {
		log.Printf("[MAYVEN DEBUG] Using cached user ID: %d", *m.userID)
		return strconv.Itoa(*m.userID)
	}
	
	log.Printf("[MAYVEN DEBUG] Fetching user ID from /api/hydrate")
	
	resp, err := m.client.Get(m.baseURI(), "/api/hydrate", m.headers(), nil)
	if err != nil {
		log.Printf("[MAYVEN DEBUG] Hydrate API error: %v", err)
		return ""
	}
	defer resp.Body.Close()
	
	log.Printf("[MAYVEN DEBUG] Hydrate response status: %d", resp.StatusCode)
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[MAYVEN DEBUG] Error reading hydrate response: %v", err)
		return ""
	}
	
	log.Printf("[MAYVEN DEBUG] Hydrate response body: %s", string(body))
	
	var hydrate MayvenHydrate
	if err := json.Unmarshal(body, &hydrate); err != nil {
		log.Printf("[MAYVEN DEBUG] Error unmarshaling hydrate response: %v", err)
		return ""
	}
	
	m.userID = &hydrate.Data.Me.Data.ID
	log.Printf("[MAYVEN DEBUG] Got user ID: %d", *m.userID)
	return strconv.Itoa(*m.userID)
}

func (m *Mayven) GetSeconds(from, to time.Time) (int, error) {
	userID := m.GetUserID()
	log.Printf("[MAYVEN DEBUG] GetUserID returned: '%s'", userID)
	
	if userID == "" {
		log.Printf("[MAYVEN DEBUG] No user ID, returning 0 seconds")
		return 0, nil
	}
	
	params := map[string]string{
		"dateStart": from.Format("2006-01-02") + " 00:00:00",
		"dateEnd":   to.Format("2006-01-02") + " 23:59:59",
		"users[]":   userID,
	}
	
	log.Printf("[MAYVEN DEBUG] API call params: %+v", params)
	
	resp, err := m.client.Get(m.baseURI(), "/api/time-statistics", m.headers(), params)
	if err != nil {
		log.Printf("[MAYVEN DEBUG] API request failed: %v", err)
		return 0, fmt.Errorf("failed to get time statistics: %w", err)
	}
	defer resp.Body.Close()
	
	log.Printf("[MAYVEN DEBUG] API response status: %d", resp.StatusCode)
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read response body: %w", err)
	}
	
	log.Printf("[MAYVEN DEBUG] API response body: %s", string(body))
	
	var stats MayvenTimeStats
	if err := json.Unmarshal(body, &stats); err != nil {
		log.Printf("[MAYVEN DEBUG] JSON unmarshal error: %v", err)
		return 0, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	
	totalSeconds := 0
	for _, item := range stats.Data.ChartData {
		totalSeconds += item.Seconds
	}
	
	log.Printf("[MAYVEN DEBUG] Chart data items: %d, total seconds: %d", len(stats.Data.ChartData), totalSeconds)
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