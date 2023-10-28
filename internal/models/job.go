package models

import "gorm.io/gorm"

type NewCompany struct {
	CompanyName string `json:"company_name" validate:"required"`
	Location    string `json:"location" validate:"required"`
}

type Company struct {
	gorm.Model
	CompanyName string `json:"company_name"  gorm:"unique;not null"`
	Location    string `json:"location"`
	UserID      uint   `json:"user_id"`
	User        User   `json:"-" gorm:"foreignkey:UserID"`
}

type NewJob struct {
	Title       string  `json:"title" validate:"required"`
	Salary      float64 `json:"salary"`
	Description string  `json:"description"`
}

type Job struct {
	gorm.Model
	Title       string  `json:"title" gorm:"not null"`
	Salary      float64 `json:"salary"`
	Description string  `json:"description"`
	CompanyID   uint    `json:"company_id" gorm:"not null"`
	Company     Company `json:"-" gorm:"foreignkey:CompanyID"`
}
