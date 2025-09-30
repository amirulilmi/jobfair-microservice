package models

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type UserType string
type EmploymentStatus string
type JobSearchStatus string

const (
	UserTypeJobSeeker UserType = "job_seeker"
	UserTypeCompany   UserType = "company"
	UserTypeAdmin     UserType = "admin"
)

const (
	EmploymentStatusStudent      EmploymentStatus = "student"
	EmploymentStatusEmployed     EmploymentStatus = "employed"
	EmploymentStatusSelfEmployed EmploymentStatus = "self_employed"
	EmploymentStatusUnemployed   EmploymentStatus = "unemployed"
)

const (
	JobSearchStatusActivelyLooking JobSearchStatus = "actively_looking"
	JobSearchStatusOpen            JobSearchStatus = "open_to_opportunities"
	JobSearchStatusNotLooking      JobSearchStatus = "not_looking"
)

// Main User Model
type User struct {
	ID       uint     `json:"id" gorm:"primaryKey"`
	Email    string   `json:"email" gorm:"uniqueIndex;not null"`
	Password string   `json:"-" gorm:"not null"`
	UserType UserType `json:"user_type" gorm:"type:varchar(20);not null"`

	// Basic Profile
	FirstName    string  `json:"first_name"`
	LastName     string  `json:"last_name"`
	PhoneNumber  *string `json:"phone_number" gorm:"index:idx_phone_unique,unique,where:phone_number IS NOT NULL"`
	CountryCode  string  `json:"country_code" gorm:"default:'+62'"`
	Country      string  `json:"country"`
	ProfilePhoto string  `json:"profile_photo"`

	// Verification Status
	IsEmailVerified bool       `json:"is_email_verified" gorm:"default:false"`
	IsPhoneVerified bool       `json:"is_phone_verified" gorm:"default:false"`
	EmailVerifiedAt *time.Time `json:"email_verified_at"`
	PhoneVerifiedAt *time.Time `json:"phone_verified_at"`

	// Account Status
	IsActive          bool `json:"is_active" gorm:"default:true"`
	IsProfileComplete bool `json:"is_profile_complete" gorm:"default:false"`

	// Timestamps
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// Job Seeker Profile Extension
type JobSeekerProfile struct {
	ID     uint `json:"id" gorm:"primaryKey"`
	UserID uint `json:"user_id" gorm:"uniqueIndex;not null"`

	// Employment Information
	CurrentJobTitle  string           `json:"current_job_title"`
	CurrentCompany   string           `json:"current_company"`
	EmploymentStatus EmploymentStatus `json:"employment_status"`
	JobSearchStatus  JobSearchStatus  `json:"job_search_status"`

	// Job Preferences
	DesiredPositions   pq.StringArray `json:"desired_positions" gorm:"type:text[]"`
	PreferredLocations pq.StringArray `json:"preferred_locations" gorm:"type:text[]"`
	JobTypes           pq.StringArray `json:"job_types" gorm:"type:text[]"`

	// Preferences & Settings
	NotificationsEnabled bool `json:"notifications_enabled" gorm:"default:true"`
	LocationEnabled      bool `json:"location_enabled" gorm:"default:false"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Company Basic Profile (stored in auth service)
type CompanyBasicProfile struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UserID      uint      `json:"user_id" gorm:"uniqueIndex;not null"`
	CompanyName string    `json:"company_name" gorm:"not null"`
	Industry    string    `json:"industry"`
	PhoneNumber string    `json:"phone_number"`
	Address     string    `json:"address"`
	Website     string    `json:"website"`
	LogoURL     string    `json:"logo_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// OTP Verification
type OTPVerification struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UserID      uint      `json:"user_id" gorm:"not null"`
	PhoneNumber string    `json:"phone_number" gorm:"not null"`
	OTPCode     string    `json:"-" gorm:"not null"`
	Purpose     string    `json:"purpose"`
	ExpiresAt   time.Time `json:"expires_at"`
	IsUsed      bool      `json:"is_used" gorm:"default:false"`
	CreatedAt   time.Time `json:"created_at"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type RegisterRequest struct {
	Email    string   `json:"email" binding:"required,email"`
	Password string   `json:"password" binding:"required,min=6"`
	UserType UserType `json:"user_type" binding:"required,oneof=job_seeker company"`
}

type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	User         User   `json:"user"`
}

// ==========================================
// DTOs for API Requests/Responses
// ==========================================

// Step 1: Initial Registration
type RegisterStep1Request struct {
	Email    string   `json:"email" binding:"required,email"`
	Password string   `json:"password" binding:"required,min=6"`
	UserType UserType `json:"user_type" binding:"required,oneof=job_seeker company"`
}

// Step 2: Basic Profile Setup - JOB SEEKER
type RegisterStep2JobSeekerRequest struct {
	FirstName   string `json:"first_name""`
	LastName    string `json:"last_name""`
	PhoneNumber string `json:"phone_number""`
	CountryCode string `json:"country_code""`
	Country     string `json:"country""`
}

// Step 2: Basic Profile Setup - COMPANY
type RegisterStep2CompanyRequest struct {
	CompanyName string `json:"company_name""`
	Industry    string `json:"industry""`
	PhoneNumber string `json:"phone_number""`
	CountryCode string `json:"country_code""`
	Address     string `json:"address""`
	Website     string `json:"website"`
}

// Step 3: Phone Verification
type PhoneVerificationRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
}

type VerifyOTPRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	OTPCode     string `json:"otp_code" binding:"required,len=6"`
}

// Step 4: Job Seeker - Employment Status
type JobSeekerStep1Request struct {
	EmploymentStatus EmploymentStatus `json:"employment_status""`
	CurrentJobTitle  string           `json:"current_job_title"`
	CurrentCompany   string           `json:"current_company"`
}

// Step 5: Job Seeker - Job Preferences
type JobSeekerStep2Request struct {
	JobSearchStatus    JobSearchStatus `json:"job_search_status" "`
	DesiredPositions   []string        `json:"desired_positions" "`
	PreferredLocations []string        `json:"preferred_locations" "`
	JobTypes           []string        `json:"job_types" "`
}

// Step 6: Notification & Location Permissions
type PermissionsRequest struct {
	NotificationsEnabled bool `json:"notifications_enabled"`
	LocationEnabled      bool `json:"location_enabled"`
}

// Step 7: Profile Photo/Logo Upload
type ProfilePhotoUploadResponse struct {
	PhotoURL string `json:"photo_url"`
}

type BasicProfileData struct {
	FirstName   string  `json:"first_name,omitempty"`
	LastName    string  `json:"last_name,omitempty"`
	CompanyName string  `json:"company_name,omitempty"`
	Industry    string  `json:"industry,omitempty"`
	PhoneNumber *string `json:"phone_number,omitempty"`
	CountryCode string  `json:"country_code,omitempty"`
	Country     string  `json:"country,omitempty"`
	Address     string  `json:"address,omitempty"`
	Website     string  `json:"website,omitempty"`
}

type EmploymentData struct {
    EmploymentStatus string `json:"employment_status"`
    CurrentJobTitle  string `json:"current_job_title,omitempty"`
    CurrentCompany   string `json:"current_company,omitempty"`
}

type JobPreferencesData struct {
	JobSearchStatus    string   `json:"job_search_status"`
	DesiredPositions   []string `json:"desired_positions"`
	PreferredLocations []string `json:"preferred_locations"`
	JobTypes           []string `json:"job_types"`
}

type OTPSentData struct {
	PhoneNumber string `json:"phone_number"`
	OTPCode     string `json:"otp_code,omitempty"`
	ExpiresAt   int64  `json:"expires_at"`
}

type ProfilePhotoData struct {
	PhotoURL string `json:"photo_url"`
}

// Responses
type RegisterStep1Response struct {
	UserID       uint     `json:"user_id"`
	Email        string   `json:"email"`
	UserType     UserType `json:"user_type"`
	NextStep     string   `json:"next_step"`
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
}

type OTPSentResponse struct {
	ExpiresIn int `json:"expires_in"`
}

type ProfileSetupResponse struct {
	IsProfileComplete bool        `json:"is_profile_complete"`
	NextStep          string      `json:"next_step"`
	PhotoURL          string      `json:"photo_url,omitempty"`
	Data              interface{} `json:"data,omitempty"`
}
