package services

import (
	"context"
	"job-portal-api/internal/models"
)

func (s *Conn) CreateCompany(ctx context.Context, nc models.NewCompany, uid uint) (models.Company, error) {
	comp := models.Company{
		CompanyName: nc.CompanyName,
		Location:    nc.Location,
		UserID:      uid,
	}

	tx := s.db.WithContext(ctx).Create(&comp)
	if tx.Error != nil {
		return models.Company{}, tx.Error
	}

	return comp, nil
}

func (s *Conn) FetchCompanies(ctx context.Context) ([]models.Company, error) {
	var comps []models.Company
	tx := s.db.WithContext(ctx).Find(&comps)
	if tx.Error != nil {
		return []models.Company{}, tx.Error
	}
	return comps, nil
}

func (s *Conn) FetchCompany(ctx context.Context, cID uint) (models.Company, error) {
	var comp models.Company
	tx := s.db.WithContext(ctx).Where("id = ?", cID)
	err := tx.First(&comp).Error
	if err != nil {
		return models.Company{}, err
	}
	return comp, nil
}

func (s *Conn) AddJobToCompany(ctx context.Context, nj models.NewJob, uID uint, cID uint) (models.Job, error) {
	var comp models.Company

	tx := s.db.WithContext(ctx).Where("id = ? and user_id = ?", cID, uID)
	err := tx.First(&comp).Error
	if err != nil {
		return models.Job{}, err
	}

	job := models.Job{
		Title:       nj.Title,
		Salary:      nj.Salary,
		Description: nj.Description,
		CompanyID:   cID,
	}

	tx = s.db.WithContext(ctx).Create(&job)
	if tx.Error != nil {
		return models.Job{}, tx.Error
	}

	return job, nil
}

func (s *Conn) ListCompanyJobs(ctx context.Context, cID uint) ([]models.Job, error) {
	var jobs []models.Job

	tx := s.db.WithContext(ctx).
		Where("company_id = ?", cID).Find(&jobs)

	if tx.Error != nil {
		return []models.Job{}, tx.Error
	}
	return jobs, nil
}

func (s *Conn) ListJobs(ctx context.Context) ([]models.Job, error) {
	var jobs []models.Job

	tx := s.db.WithContext(ctx).Find(&jobs)
	if tx.Error != nil {
		return []models.Job{}, tx.Error
	}
	return jobs, nil
}

func (s *Conn) GetJobByID(ctx context.Context, jID uint) (models.Job, error) {
	var job models.Job

	tx := s.db.WithContext(ctx).
		Where("id = ?", jID).First(&job)

	if tx.Error != nil {
		return models.Job{}, tx.Error
	}
	return job, nil
}
