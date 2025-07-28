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
	
	hours, err := h.trackersRepo.Hours(date, tomorrow)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get hours"})
		return
	}
	
	runningHours, err := h.trackersRepo.RunningHours()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get running hours"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"date":          date.Format("2006-01-02"),
		"hours":         hours,
		"running_hours": runningHours,
	})
}