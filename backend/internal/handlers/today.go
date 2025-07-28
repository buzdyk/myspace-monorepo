package handlers

import (
	"myspace/backend/internal/repositories"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type TodayHandler struct {
	trackersRepo *repositories.TrackersRepository
}

func NewTodayHandler(trackersRepo *repositories.TrackersRepository) *TodayHandler {
	return &TodayHandler{
		trackersRepo: trackersRepo,
	}
}

func (h *TodayHandler) Redirect(c *gin.Context) {
	now := time.Now()
	url := "/" + strconv.Itoa(now.Year()) + "/" + strconv.Itoa(int(now.Month())) + "/" + strconv.Itoa(now.Day())
	c.Redirect(http.StatusFound, url)
}

func (h *TodayHandler) Index(c *gin.Context) {
	yearStr := c.Param("year")
	monthStr := c.Param("month")
	dayStr := c.Param("day")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid year"})
		return
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid month"})
		return
	}

	day, err := strconv.Atoi(dayStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid day"})
		return
	}

	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	tomorrow := date.AddDate(0, 0, 1)
	now := time.Now()

	// Get today's hours
	todayHours, err := h.trackersRepo.Hours(date, tomorrow)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get hours"})
		return
	}

	// Get running hours
	runningHours, err := h.trackersRepo.RunningHours()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get running hours"})
		return
	}

	// Get month hours
	monthStart := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	monthEnd := date.AddDate(0, 0, 1) // End of the current day
	monthHours, err := h.trackersRepo.Hours(monthStart, monthEnd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get month hours"})
		return
	}

	// Add running hours if it's today
	isToday := date.Year() == now.Year() && date.Month() == now.Month() && date.Day() == now.Day()
	if isToday {
		monthHours += runningHours
	}

	// Settings (hardcoded for now - should come from database)
	dailyGoal := 8.0
	monthlyGoal := 160.0

	// Calculate percentages
	todayPercent := int((todayHours / dailyGoal) * 100)
	monthPercent := (monthHours / monthlyGoal) * 100

	// Calculate pace (simplified)
	daysInMonth := float64(time.Date(year, time.Month(month+1), 0, 0, 0, 0, 0, time.UTC).Day())
	dayOfMonth := float64(day)
	remainingDays := daysInMonth - dayOfMonth
	expectedHours := remainingDays * dailyGoal
	remainingHours := monthlyGoal - monthHours
	pace := expectedHours - remainingHours

	// Navigation links
	prevDay := date.AddDate(0, 0, -1)
	nextDay := date.AddDate(0, 0, 1)

	nav := gin.H{
		"month": date.Format("January"),
		"day":   date.Format("2nd"),
		"year": func() string {
			if date.Year() == now.Year() {
				return ""
			} else {
				return date.Format("2006")
			}
		}(),
		"month_link": "/" + strconv.Itoa(year) + "/" + strconv.Itoa(month) + "/calendar",
		"prev_link":  "/" + strconv.Itoa(prevDay.Year()) + "/" + strconv.Itoa(int(prevDay.Month())) + "/" + strconv.Itoa(prevDay.Day()),
		"next_link":  "/" + strconv.Itoa(nextDay.Year()) + "/" + strconv.Itoa(int(nextDay.Month())) + "/" + strconv.Itoa(nextDay.Day()),
	}

	c.JSON(http.StatusOK, gin.H{
		"date":          date.Format("2006-01-02"),
		"hours":         todayHours,
		"running_hours": runningHours,
		"today_percent": todayPercent,
		"month_percent": monthPercent,
		"month_hours":   monthHours,
		"pace":          pace,
		"daily_goal":    dailyGoal,
		"is_today":      isToday,
		"nav":           nav,
	})
}
