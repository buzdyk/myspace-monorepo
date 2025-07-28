package config

import (
	"os"
)

type Config struct {
	Port string
	
	Database struct {
		Path string
	}
	
	Clockify struct {
		Token       string
		WorkspaceID string
		UserID      string
	}
	
	Everhour struct {
		Token string
	}
	
	Mayven struct {
		Auth string
	}
}

func Load() *Config {
	cfg := &Config{
		Port: getEnv("PORT", "8080"),
	}
	
	cfg.Database.Path = getEnv("DB_PATH", "./database.sqlite")
	
	cfg.Clockify.Token = getEnv("CLOCKIFY_TOKEN", "")
	cfg.Clockify.WorkspaceID = getEnv("CLOCKIFY_WORKSPACE_ID", "")
	cfg.Clockify.UserID = getEnv("CLOCKIFY_USER_ID", "")
	
	cfg.Everhour.Token = getEnv("EVERHOUR_TOKEN", "")
	cfg.Mayven.Auth = getEnv("MAYVEN_AUTH", "")
	
	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}