package services

import (
	"fmt"
	"job-portal-api/internal/models"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func (s *Conn) CreateUser(nu models.NewUser) (models.User, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(nu.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, fmt.Errorf("generating password hash: %w", err)
	}

	u := models.User{
		Name:         nu.Name,
		Email:        nu.Email,
		PasswordHash: string(hashedPass),
	}
	tx := s.db.Create(&u)
	if tx.Error != nil {
		return models.User{}, tx.Error
	}

	return u, nil
}

func (s *Conn) Authenticate(lu models.UserLogin) (jwt.RegisteredClaims, error) {
	var u models.User
	tx := s.db.Where("email = ?", lu.Email).First(&u)
	if tx.Error != nil {
		return jwt.RegisteredClaims{}, tx.Error
	}

	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(lu.Password))
	if err != nil {
		return jwt.RegisteredClaims{}, err
	}

	claims := jwt.RegisteredClaims{
		Issuer:    "job portal",
		Subject:   strconv.FormatUint(uint64(u.ID), 10),
		Audience:  jwt.ClaimStrings{"job seekers"},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	return claims, nil
}
