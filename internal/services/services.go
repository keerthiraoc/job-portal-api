package services

import (
	"errors"

	"gorm.io/gorm"
)

type Conn struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) (*Conn, error) {
	if db == nil {
		return nil, errors.New("please provide a valid connection")
	}
	return &Conn{db: db}, nil
}
