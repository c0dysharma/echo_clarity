package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name                 string
	Email                string `gorm:"uniqueIndex"`
	RefreshToken         string
	AccessToken          string
	AccessTokenExpiresAt time.Time
}
