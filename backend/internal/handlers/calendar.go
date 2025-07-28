package handlers

import (
	"myspace/backend/internal/repositories"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type CalendarHandler struct {
	trackersRepo *repositories.TrackersRepository
}

func NewCalendarHandler(trackersRepo *repositories.TrackersRepository) *CalendarHandler {
	return &CalendarHandler{
		trackersRepo: trackersRepo,
	}
}

func (h *CalendarHandler) Index(c *gin.Context) {
	yearStr := c.Param("year")
	monthStr := c.Param("month")
	
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
	
	date := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	
	dailyHours, err := h.trackersRepo.GetDailyHours(date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get daily hours"})
		return
	}
	
	days := h.getDays(date, dailyHours)
	
	c.JSON(http.StatusOK, gin.H{
		"year":  year,
		"month": month,
		"days":  days,
	})
}

func (h *CalendarHandler) getDays(date time.Time, dailyHours map[string]*float64) []map[string]interface{} {
	som := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
	eom := som.AddDate(0, 1, -1)
	
	var days []map[string]interface{}
	
	firstWeekday := int(som.Weekday())
	if firstWeekday == 0 {
		firstWeekday = 7
	}
	
	for i := 1; i < firstWeekday; i++ {
		days = append(days, map[string]interface{}{
			"day":   nil,
			"hours": nil,
		})
	}
	
	for d := som; !d.After(eom); d = d.AddDate(0, 0, 1) {
		dayStr := d.Format("2006-01-02")
		hours := dailyHours[dayStr]
		
		days = append(days, map[string]interface{}{
			"day":   d.Day(),
			"hours": hours,
		})
	}
	
	return days
}