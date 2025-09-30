package repository

import (
	"jobfair-company-service/internal/models"

	"gorm.io/gorm"
)

type ApplicationRepository struct {
	db *gorm.DB
}

func NewApplicationRepository(db *gorm.DB) *ApplicationRepository {
	return &ApplicationRepository{db: db}
}

func (r *ApplicationRepository) Create(application *models.JobApplication) (*models.JobApplication, error) {
	if err := r.db.Create(application).Error; err != nil {
		return nil, err
	}
	return application, nil
}

func (r *ApplicationRepository) GetByID(id uint) (*models.JobApplication, error) {
	var application models.JobApplication
	if err := r.db.First(&application, id).Error; err != nil {
		return nil, err
	}
	return &application, nil
}

func (r *ApplicationRepository) Update(application *models.JobApplication) error {
	return r.db.Save(application).Error
}

func (r *ApplicationRepository) List(companyID uint, limit, offset int, filters map[string]interface{}) ([]*models.JobApplication, int64, error) {
	var applications []*models.JobApplication
	var total int64

	query := r.db.Model(&models.JobApplication{}).Where("company_id = ?", companyID)

	for key, value := range filters {
		if value != "" && value != nil {
			query = query.Where(key+" = ?", value)
		}
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(limit).Offset(offset).Order("applied_at DESC").Find(&applications).Error; err != nil {
		return nil, 0, err
	}

	return applications, total, nil
}

func (r *ApplicationRepository) GetByJobID(jobID uint) ([]*models.JobApplication, error) {
	var applications []*models.JobApplication
	if err := r.db.Where("job_id = ?", jobID).Order("applied_at DESC").Find(&applications).Error; err != nil {
		return nil, err
	}
	return applications, nil
}

func (r *ApplicationRepository) GetByCompanyID(companyID uint) ([]*models.JobApplication, error) {
	var applications []*models.JobApplication
	if err := r.db.Where("company_id = ?", companyID).Order("applied_at DESC").Find(&applications).Error; err != nil {
		return nil, err
	}
	return applications, nil
}

func (r *ApplicationRepository) CountByCompanyID(companyID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&models.JobApplication{}).Where("company_id = ?", companyID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *ApplicationRepository) CountByJobID(jobID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&models.JobApplication{}).Where("job_id = ?", jobID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *ApplicationRepository) CountByStatus(companyID uint, status models.ApplicationStatus) (int64, error) {
	var count int64
	if err := r.db.Model(&models.JobApplication{}).Where("company_id = ? AND status = ?", companyID, status).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *ApplicationRepository) UpdateStatus(id uint, status models.ApplicationStatus) error {
	return r.db.Model(&models.JobApplication{}).Where("id = ?", id).Update("status", status).Error
}