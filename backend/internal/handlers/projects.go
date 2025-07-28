package handlers

import (
	"myspace/backend/internal/repositories"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ProjectsHandler struct {
	trackersRepo *repositories.TrackersRepository
}

func NewProjectsHandler(trackersRepo *repositories.TrackersRepository) *ProjectsHandler {
	return &ProjectsHandler{
		trackersRepo: trackersRepo,
	}
}

func (h *ProjectsHandler) Redirect(c *gin.Context) {
	now := time.Now()
	url := "/" + strconv.Itoa(now.Year()) + "/" + strconv.Itoa(int(now.Month())) + "/projects"
	c.Redirect(http.StatusFound, url)
}

func (h *ProjectsHandler) Index(c *gin.Context) {
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
	
	projectTimes, err := h.trackersRepo.GetMonthlyTimeByProject(date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get project times"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"year":     year,
		"month":    month,
		"projects": projectTimes.ToArray(),
		"total_hours": projectTimes.GetHours(),
	})
}