package models

import (
	"time"
)

type Patient struct {
	ID           uint      `gorm:"primaryKey"`
	FirstName    string    `gorm:"type:varchar(100);not null"`
	LastName     string    `gorm:"type:varchar(100)"`
	DOB          time.Time `gorm:"not null"`
	Gender       string    `gorm:"type:varchar(10);check:gender IN ('male','female','other')"`
	Phone        string    `gorm:"type:varchar(20)"`
	Email        string    `gorm:"type:varchar(100)"`
	Address      string
	RegisteredBy *uint // foreign key (nullable in case receptionist is deleted)
	Receptionist *User `gorm:"foreignKey:RegisteredBy;constraint:OnDelete:SET NULL"` // association
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Notes        []PatientNote `gorm:"foreignKey:PatientID"` // has-many
}
