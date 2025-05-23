package entities

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email            string `gorm:"unique;not null" json:"email"`
	Password         string `gorm:"not null" json:"password"`
	Role             string `gorm:"type:varchar(20);not null" json:"role"`
	IsVerified       bool   `gorm:"default:false" json:"is_verified"`
	VerificationCode string `gorm:"size:100" json:"-"`
}
