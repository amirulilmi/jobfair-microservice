// File: internal/repository/company_basic_profile_repository.go
package repository

import (
	"jobfair-auth-service/internal/models"

	"gorm.io/gorm"
)

type CompanyBasicProfileRepository struct {
	db *gorm.DB
}

func NewCompanyBasicProfileRepository(db *gorm.DB) *CompanyBasicProfileRepository {
	return &CompanyBasicProfileRepository{db: db}
}

func (r *CompanyBasicProfileRepository) Create(profile *models.CompanyBasicProfile) error {
	return r.db.Create(profile).Error
}

func (r *CompanyBasicProfileRepository) GetByUserID(userID uint) (*models.CompanyBasicProfile, error) {
	var profile models.CompanyBasicProfile
	if err := r.db.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *CompanyBasicProfileRepository) Update(profile *models.CompanyBasicProfile) error {
	return r.db.Save(profile).Error
}

func (r *CompanyBasicProfileRepository) Delete(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.CompanyBasicProfile{}).Error
}