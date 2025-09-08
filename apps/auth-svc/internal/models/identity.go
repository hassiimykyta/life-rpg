package models

import "time"

type Identity struct {
	UserId       string    `gorm:"primaryKey;size36"`
	Email        string    `gorm:"size:255;uniqueIndex;not null"`
	Username     string    `gorm:"size:64;uniqueIndex;not null"`
	PasswordHash string    `gorm:"not null"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}

func (Identity) TableName() string { return "identity" }
