package database

import (
	"time"
)

type User struct {
	ID                uint      `json:"id" gorm:"primaryKey"`
	Name              string    `json:"name"`
	Email             string    `json:"email" gorm:"uniqueIndex"`
	EmailVerifiedAt   *time.Time `json:"email_verified_at"`
	Password          string    `json:"-"`
	RememberToken     string    `json:"-"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type Tracker struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Config    string    `json:"config" gorm:"type:json"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Projects  []Project `json:"projects" gorm:"foreignKey:TrackerID"`
	Tracks    []Track   `json:"tracks" gorm:"foreignKey:TrackerID"`
}

type Project struct {
	ID        uint    `json:"id" gorm:"primaryKey"`
	TrackerID uint    `json:"tracker_id"`
	Name      string  `json:"name"`
	Token     string  `json:"token" gorm:"type:text"`
	Tracker   Tracker `json:"tracker" gorm:"foreignKey:TrackerID"`
	Tracks    []Track `json:"tracks" gorm:"foreignKey:ProjectID"`
}

type Track struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	TrackerID uint      `json:"tracker_id"`
	ProjectID uint      `json:"project_id"`
	Seconds   int       `json:"seconds"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Tracker   Tracker   `json:"tracker" gorm:"foreignKey:TrackerID"`
	Project   Project   `json:"project" gorm:"foreignKey:ProjectID"`
}

type Setting struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Key       string    `json:"key" gorm:"uniqueIndex"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}