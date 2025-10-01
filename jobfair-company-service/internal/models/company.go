package models

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type CompanySize string
type SubscriptionTier string

const (
	CompanySize1to10     CompanySize = "1-10"
	CompanySize11to50    CompanySize = "11-50"
	CompanySize51to200   CompanySize = "51-200"
	CompanySize201to500  CompanySize = "201-500"
	CompanySize501to1000 CompanySize = "501-1000"
	CompanySize1000Plus  CompanySize = "1000+"
)

const (
	SubscriptionFree    SubscriptionTier = "free"
	SubscriptionBasic   SubscriptionTier = "basic"
	SubscriptionPremium SubscriptionTier = "premium"
	SubscriptionPro     SubscriptionTier = "pro"
)

type Company struct {
	ID                uint             `json:"id" gorm:"primaryKey"`
	UserID            uint             `json:"user_id" gorm:"uniqueIndex;not null"`
	Name              string           `json:"name" gorm:"not null"`
	Description       string           `json:"description" gorm:"type:text"`
	Industry          string           `json:"industry"`
	CompanySize       CompanySize      `json:"company_size"`
	FoundedYear       int              `json:"founded_year"`
	Email             string           `json:"email"`
	Phone             string           `json:"phone"`
	Website           string           `json:"website"`
	Address           string           `json:"address"`
	City              string           `json:"city"`
	State             string           `json:"state"`
	Country           string           `json:"country"`
	PostalCode        string           `json:"postal_code"`
	Latitude          float64          `json:"latitude,omitempty"`
	Longitude         float64          `json:"longitude,omitempty"`
	LogoURL           string           `json:"logo_url"`
	BannerURL         string           `json:"banner_url"`
	VideoURLs         pq.StringArray   `json:"video_urls" gorm:"type:text[]"`
	GalleryURLs       pq.StringArray   `json:"gallery_urls" gorm:"type:text[]"`
	LinkedinURL       string           `json:"linkedin_url"`
	FacebookURL       string           `json:"facebook_url"`
	TwitterURL        string           `json:"twitter_url"`
	InstagramURL      string           `json:"instagram_url"`
	IsVerified        bool             `json:"is_verified" gorm:"default:false"`
	VerifiedAt        *time.Time       `json:"verified_at"`
	VerificationBadge string           `json:"verification_badge"`
	IsFeatured        bool             `json:"is_featured" gorm:"default:false"`
	IsPremium         bool             `json:"is_premium" gorm:"default:false"`
	SubscriptionTier  SubscriptionTier `json:"subscription_tier" gorm:"default:'free'"`
	Slug              string           `json:"slug" gorm:"uniqueIndex"`
	MetaTitle         string           `json:"meta_title"`
	MetaDescription   string           `json:"meta_description"`
	Tags              pq.StringArray   `json:"tags" gorm:"type:text[]"`
	CreatedAt         time.Time        `json:"created_at"`
	UpdatedAt         time.Time        `json:"updated_at"`
	DeletedAt         gorm.DeletedAt   `json:"-" gorm:"index"`
}

type CreateCompanyRequest struct {
	Name        string      `json:"name" binding:"required"`
	Description string      `json:"description"`
	Industry    string      `json:"industry"`
	CompanySize CompanySize `json:"company_size"`
	FoundedYear int         `json:"founded_year"`
	Email       string      `json:"email" binding:"required,email"`
	Phone       string      `json:"phone"`
	Website     string      `json:"website"`
	Address     string      `json:"address"`
	City        string      `json:"city"`
	Country     string      `json:"country"`
}

type UpdateCompanyRequest struct {
	Name         *string      `json:"name"`
	Description  *string      `json:"description"`
	Industry     *string      `json:"industry"`
	CompanySize  *CompanySize `json:"company_size"`
	FoundedYear  *int         `json:"founded_year"`
	Email        *string      `json:"email"`
	Phone        *string      `json:"phone"`
	Website      *string      `json:"website"`
	Address      *string      `json:"address"`
	City         *string      `json:"city"`
	State        *string      `json:"state"`
	Country      *string      `json:"country"`
	PostalCode   *string      `json:"postal_code"`
	LinkedinURL  *string      `json:"linkedin_url"`
	FacebookURL  *string      `json:"facebook_url"`
	TwitterURL   *string      `json:"twitter_url"`
	InstagramURL *string      `json:"instagram_url"`
}

type CompanyAnalytics struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	CompanyID       uint      `json:"company_id" gorm:"uniqueIndex;not null"`
	BoothVisits     int       `json:"booth_visits" gorm:"default:0"`
	ProfileViews    int       `json:"profile_views" gorm:"default:0"`
	JobViews        int       `json:"job_views" gorm:"default:0"`
	Applications    int       `json:"applications" gorm:"default:0"`
	TotalJobsPosted int       `json:"total_jobs_posted" gorm:"default:0"`
	ActiveJobs      int       `json:"active_jobs" gorm:"default:0"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type CompanyMedia struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	CompanyID    uint           `json:"company_id" gorm:"not null;index"`
	MediaType    string         `json:"media_type"`
	MediaURL     string         `json:"media_url" gorm:"not null"`
	Title        string         `json:"title"`
	Description  string         `json:"description"`
	DisplayOrder int            `json:"display_order" gorm:"default:0"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

type DashboardStats struct {
	TotalJobsPosted    int            `json:"total_jobs_posted"`
	TotalApplicants    int            `json:"total_applicants"`
	JobViews           int            `json:"job_views"`
	TotalPositions     int            `json:"total_positions"`
	JobsGrowth         float64        `json:"jobs_growth"`
	ApplicantsGrowth   float64        `json:"applicants_growth"`
	ViewsGrowth        float64        `json:"views_growth"`
	PositionsGrowth    float64        `json:"positions_growth"`
	RecentJobs         []JobStatus    `json:"recent_jobs"`
	ApplicantTrend     []TrendData    `json:"applicant_trend"`
	ApplicantsByStatus map[string]int `json:"applicants_by_status"`
}

type JobProgress struct {
	Title          string `json:"title"`
	Status         string `json:"status"`
	StepsTotal     int    `json:"steps_total"`
	StepsCompleted int    `json:"steps_completed"`
}

type TrendData struct {
	Month string `json:"month"`
	Value int    `json:"value"`
}
