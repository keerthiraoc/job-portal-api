package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name         string `json:"username" gorm:"not null"`
	Email        string `json:"email" gorm:"not null;unique"`
	PasswordHash string `json:"-"`
}

type NewUser struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email"  validate:"required,email"`
	Password string `json:"password"  validate:"required"`
}

type UserLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type NewToken struct {
	Token string `json:"token"`
}
