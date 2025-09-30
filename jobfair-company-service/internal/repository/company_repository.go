package repository

import (
	"errors"
	"time"

	"jobfair-company-service/internal/models"
	"gorm.io/gorm"
)

type CompanyRepository struct {
	db *gorm.DB
}

func NewCompanyRepository(db *gorm.DB) *CompanyRepository {
	return &CompanyRepository{db: db}
}

func (r *CompanyRepository) Create(company *models.Company) (*models.Company, error) {
	if err := r.db.Create(company).Error; err != nil {
		return nil, err
	}
	return company, nil
}

func (r *CompanyRepository) GetByID(id uint) (*models.Company, error) {
	var company models.Company
	if err := r.db.First(&company, id).Error; err != nil {
		return nil, err
	}
	return &company, nil
}

func (r *CompanyRepository) GetByUserID(userID uint) (*models.Company, error) {
	var company models.Company
	if err := r.db.Where("user_id = ?", userID).First(&company).Error; err != nil {
		return nil, err
	}
	return &company, nil
}

func (r *CompanyRepository) Update(company *models.Company) error {
	return r.db.Save(company).Error
}

func (r *CompanyRepository) Delete(id uint) error {
	return r.db.Delete(&models.Company{}, id).Error
}

func (r *CompanyRepository) List(limit, offset int, filters map[string]interface{}) ([]*models.Company, int64, error) {
	var companies []*models.Company
	var total int64

	query := r.db.Model(&models.Company{})

	for key, value := range filters {
		if value != "" && value != nil {
			query = query.Where(key+" = ?", value)
		}
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(limit).Offset(offset).Order("created_at DESC").Find(&companies).Error; err != nil {
		return nil, 0, err
	}

	return companies, total, nil
}

func (r *CompanyRepository) GetAnalytics(companyID uint) (*models.CompanyAnalytics, error) {
	var analytics models.CompanyAnalytics
	if err := r.db.Where("company_id = ?", companyID).First(&analytics).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newAnalytics := &models.CompanyAnalytics{
				CompanyID: companyID,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			if err := r.db.Create(newAnalytics).Error; err != nil {
				return nil, err
			}
			return newAnalytics, nil
		}
		return nil, err
	}
	return &analytics, nil
}

func (r *CompanyRepository) UpdateAnalytics(analytics *models.CompanyAnalytics) error {
	return r.db.Save(analytics).Error
}

func (r *CompanyRepository) IncrementProfileViews(companyID uint) error {
	return r.db.Model(&models.CompanyAnalytics{}).Where("company_id = ?", companyID).UpdateColumn("profile_views", gorm.Expr("profile_views + ?", 1)).Error
}

func (r *CompanyRepository) IncrementBoothVisits(companyID uint) error {
	return r.db.Model(&models.CompanyAnalytics{}).Where("company_id = ?", companyID).UpdateColumn("booth_visits", gorm.Expr("booth_visits + ?", 1)).Error
}