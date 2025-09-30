package repository

import (
	"jobfair-auth-service/internal/models"

	"gorm.io/gorm"
)

type JobSeekerProfileRepository struct {
	db *gorm.DB
}

func NewJobSeekerProfileRepository(db *gorm.DB) *JobSeekerProfileRepository {
	return &JobSeekerProfileRepository{db: db}
}

func (r *JobSeekerProfileRepository) Create(profile *models.JobSeekerProfile) error {
	return r.db.Create(profile).Error
}

func (r *JobSeekerProfileRepository) GetByUserID(userID uint) (*models.JobSeekerProfile, error) {
	var profile models.JobSeekerProfile
	if err := r.db.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *JobSeekerProfileRepository) Update(profile *models.JobSeekerProfile) error {
	return r.db.Save(profile).Error
}

func (r *JobSeekerProfileRepository) Delete(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.JobSeekerProfile{}).Error
}
