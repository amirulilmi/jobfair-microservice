package models

import (
	"time"

	"gorm.io/gorm"
)

// Company represents a company profile
type Company struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UserID      uint           `json:"user_id" gorm:"not null"` // Foreign key to auth service
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description" gorm:"type:text"`
	Industry    string         `json:"industry"`
	Location    string         `json:"location"`
	Website     string         `json:"website"`
	Email       string         `json:"email"`
	Phone       string         `json:"phone"`
	LogoURL     string         `json:"logo_url"`
	BannerURL   string         `json:"banner_url"`
	VideoURLs   []string       `json:"video_urls" gorm:"type:text[]"`
	IsVerified  bool           `json:"is_verified" gorm:"default:false"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// CreateCompanyRequest used for creating a new company
type CreateCompanyRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Industry    string `json:"industry"`
	Location    string `json:"location"`
	Website     string `json:"website"`
	Email       string `json:"email" binding:"required,email"`
	Phone       string `json:"phone"`
}

// UpdateCompanyRequest used for updating existing company
type UpdateCompanyRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Industry    string `json:"industry"`
	Location    string `json:"location"`
	Website     string `json:"website"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
}

// CompanyAnalytics stores analytics data for a company
type CompanyAnalytics struct {
	CompanyID    uint      `json:"company_id"`
	BoothVisits  int       `json:"booth_visits"`
	ProfileViews int       `json:"profile_views"`
	JobViews     int       `json:"job_views"`
	Applications int       `json:"applications"`
	LastUpdated  time.Time `json:"last_updated"`
}
