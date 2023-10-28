package services

import (
	"job-portal-api/internal/models"
)

func (s *Conn) AutoMigrate() error {
	err := s.db.Migrator().DropTable(&models.Company{}, &models.User{}, &models.Job{})
	if err != nil {
		return err
	}

	err = s.db.Migrator().AutoMigrate(&models.Company{}, &models.User{}, &models.Job{})
	if err != nil {
		return err
	}
	return nil
}
