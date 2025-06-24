package models

import "time"

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"type:varchar(100);not null"`
	Email     string `gorm:"type:varchar(100);uniqueIndex;not null"`
	Password  string `gorm:"not null"`
	Role      string `gorm:"type:varchar(20);not null;check:role IN ('doctor','receptionist')"`
	CreatedAt time.Time
}
