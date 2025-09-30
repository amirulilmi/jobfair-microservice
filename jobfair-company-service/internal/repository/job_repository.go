package repository

import (
	"jobfair-company-service/internal/models"
	"time"

	"gorm.io/gorm"
)

type JobRepository struct {
	db *gorm.DB
}

func NewJobRepository(db *gorm.DB) *JobRepository {
	return &JobRepository{db: db}
}

func (r *JobRepository) Create(job *models.Job) (*models.Job, error) {
	if err := r.db.Create(job).Error; err != nil {
		return nil, err
	}
	return job, nil
}

func (r *JobRepository) GetByID(id uint) (*models.Job, error) {
	var job models.Job
	if err := r.db.First(&job, id).Error; err != nil {
		return nil, err
	}
	return &job, nil
}

func (r *JobRepository) Update(job *models.Job) error {
	return r.db.Save(job).Error
}

func (r *JobRepository) Delete(id uint) error {
	return r.db.Delete(&models.Job{}, id).Error
}

func (r *JobRepository) List(companyID uint, limit, offset int, filters map[string]interface{}) ([]*models.Job, int64, error) {
	var jobs []*models.Job
	var total int64

	query := r.db.Model(&models.Job{}).Where("company_id = ?", companyID)

	for key, value := range filters {
		if value != "" && value != nil {
			query = query.Where(key+" = ?", value)
		}
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&jobs).Error; err != nil {
		return nil, 0, err
	}

	return jobs, total, nil
}

func (r *JobRepository) GetByCompanyID(companyID uint) ([]*models.Job, error) {
	var jobs []*models.Job
	if err := r.db.Where("company_id = ? AND deleted_at IS NULL", companyID).Order("created_at DESC").Find(&jobs).Error; err != nil {
		return nil, err
	}
	return jobs, nil
}

func (r *JobRepository) IncrementViewCount(jobID uint) error {
	return r.db.Model(&models.Job{}).Where("id = ?", jobID).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
}

func (r *JobRepository) UpdateStatus(jobID uint, status models.JobStatus) error {
	updates := map[string]interface{}{
		"status": status,
	}

	if status == models.JobStatusActive {
		updates["published_at"] = time.Now()
	}

	return r.db.Model(&models.Job{}).Where("id = ?", jobID).Updates(updates).Error
}

func (r *JobRepository) GetActiveJobsCount(companyID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&models.Job{}).Where("company_id = ? AND status = ?", companyID, models.JobStatusActive).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *JobRepository) GetTotalJobsCount(companyID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&models.Job{}).Where("company_id = ?", companyID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}