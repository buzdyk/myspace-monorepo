package main

import (
	"log"
	"myspace/backend/internal/config"
	"myspace/backend/internal/database"
	"myspace/backend/internal/handlers"
	"myspace/backend/internal/repositories"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	
	db, err := database.Connect(cfg.Database.Path)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	
	if err := database.Migrate(db); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	
	trackersRepo := repositories.NewTrackersRepository(cfg)
	
	todayHandler := handlers.NewTodayHandler(trackersRepo)
	projectsHandler := handlers.NewProjectsHandler(trackersRepo)
	calendarHandler := handlers.NewCalendarHandler(trackersRepo)
	
	r := gin.Default()
	
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})
	
	r.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/today")
	})
	
	r.GET("/today", todayHandler.Redirect)
	r.GET("/:year/:month/:day", todayHandler.Index)
	
	r.GET("/month", projectsHandler.Redirect)
	r.GET("/:year/:month/projects", projectsHandler.Index)
	r.GET("/:year/:month/calendar", calendarHandler.Index)
	
	log.Printf("Server starting on port %s", cfg.Port)
	r.Run(":" + cfg.Port)
}