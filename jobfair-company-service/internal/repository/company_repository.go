package repository

import (
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

func (r *CompanyRepository) List(limit, offset int) ([]*models.Company, error) {
    var companies []*models.Company
    if err := r.db.Limit(limit).Offset(offset).Find(&companies).Error; err != nil {
        return nil, err
    }
    return companies, nil
}

func (r *CompanyRepository) UpdateAnalytics(companyID uint, analytics *models.CompanyAnalytics) error {
    return r.db.Model(&models.CompanyAnalytics{}).Where("company_id = ?", companyID).Updates(analytics).Error
}

func (r *CompanyRepository) GetAnalytics(companyID uint) (*models.CompanyAnalytics, error) {
    var analytics models.CompanyAnalytics
    if err := r.db.Where("company_id = ?", companyID).First(&analytics).Error; err != nil {
        return nil, err
    }
    return &analytics, nil
}