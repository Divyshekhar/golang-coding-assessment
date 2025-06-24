package models

import "time"

type PatientNote struct {
	ID        uint    `gorm:"primaryKey"`
	PatientID uint    `gorm:"not null"`
	Patient   Patient `gorm:"constraint:OnDelete:CASCADE"`
	DoctorID  *uint   // nullable if doctor is deleted
	Doctor    *User   `gorm:"foreignKey:DoctorID;constraint:OnDelete:SET NULL"`
	Note      string  `gorm:"type:text;not null"`
	CreatedAt time.Time
}
