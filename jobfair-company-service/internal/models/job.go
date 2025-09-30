package models

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type JobType string
type JobLevel string
type JobStatus string
type ApplicationStatus string

const (
	JobTypeFullTime   JobType = "full_time"
	JobTypePartTime   JobType = "part_time"
	JobTypeContract   JobType = "contract"
	JobTypeInternship JobType = "internship"
	JobTypeFreelance  JobType = "freelance"
)

const (
	JobLevelEntry      JobLevel = "entry"
	JobLevelJunior     JobLevel = "junior"
	JobLevelMid        JobLevel = "mid"
	JobLevelSenior     JobLevel = "senior"
	JobLevelLead       JobLevel = "lead"
	JobLevelManager    JobLevel = "manager"
	JobLevelDirector   JobLevel = "director"
	JobLevelExecutive  JobLevel = "executive"
)

const (
	JobStatusDraft     JobStatus = "draft"
	JobStatusActive    JobStatus = "active"
	JobStatusPaused    JobStatus = "paused"
	JobStatusClosed    JobStatus = "closed"
	JobStatusExpired   JobStatus = "expired"
)

const (
	ApplicationStatusShortlisted ApplicationStatus = "shortlisted"
	ApplicationStatusHired       ApplicationStatus = "hired"
	ApplicationStatusInterview   ApplicationStatus = "interview"
	ApplicationStatusRejected    ApplicationStatus = "rejected"
	ApplicationStatusApplied     ApplicationStatus = "applied"
)

type Job struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	CompanyID   uint           `json:"company_id" gorm:"not null;index"`
	Title       string         `json:"title" gorm:"not null"`
	Description string         `json:"description" gorm:"type:text"`
	Requirements string        `json:"requirements" gorm:"type:text"`
	Responsibilities string    `json:"responsibilities" gorm:"type:text"`
	JobType     JobType        `json:"job_type" gorm:"not null"`
	JobLevel    JobLevel       `json:"job_level"`
	Status      JobStatus      `json:"status" gorm:"default:'draft'"`
	Location    string         `json:"location"`
	City        string         `json:"city"`
	Country     string         `json:"country"`
	IsRemote    bool           `json:"is_remote" gorm:"default:false"`
	SalaryMin   int            `json:"salary_min"`
	SalaryMax   int            `json:"salary_max"`
	SalaryCurrency string      `json:"salary_currency" gorm:"default:'IDR'"`
	ShowSalary  bool           `json:"show_salary" gorm:"default:false"`
	Skills      pq.StringArray `json:"skills" gorm:"type:text[]"`
	Benefits    pq.StringArray `json:"benefits" gorm:"type:text[]"`
	Positions   int            `json:"positions" gorm:"default:1"`
	ViewCount   int            `json:"view_count" gorm:"default:0"`
	ApplicationCount int       `json:"application_count" gorm:"default:0"`
	ExpiresAt   *time.Time     `json:"expires_at"`
	PublishedAt *time.Time     `json:"published_at"`
	Slug        string         `json:"slug" gorm:"index"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type CreateJobRequest struct {
	Title       string         `json:"title" binding:"required"`
	Description string         `json:"description" binding:"required"`
	Requirements string        `json:"requirements"`
	Responsibilities string    `json:"responsibilities"`
	JobType     JobType        `json:"job_type" binding:"required"`
	JobLevel    JobLevel       `json:"job_level"`
	Location    string         `json:"location"`
	City        string         `json:"city"`
	Country     string         `json:"country"`
	IsRemote    bool           `json:"is_remote"`
	SalaryMin   int            `json:"salary_min"`
	SalaryMax   int            `json:"salary_max"`
	ShowSalary  bool           `json:"show_salary"`
	Skills      []string       `json:"skills"`
	Benefits    []string       `json:"benefits"`
	Positions   int            `json:"positions"`
	ExpiresAt   *time.Time     `json:"expires_at"`
}

type UpdateJobRequest struct {
	Title       *string        `json:"title"`
	Description *string        `json:"description"`
	Requirements *string       `json:"requirements"`
	Responsibilities *string   `json:"responsibilities"`
	JobType     *JobType       `json:"job_type"`
	JobLevel    *JobLevel      `json:"job_level"`
	Status      *JobStatus     `json:"status"`
	Location    *string        `json:"location"`
	City        *string        `json:"city"`
	Country     *string        `json:"country"`
	IsRemote    *bool          `json:"is_remote"`
	SalaryMin   *int           `json:"salary_min"`
	SalaryMax   *int           `json:"salary_max"`
	ShowSalary  *bool          `json:"show_salary"`
	Skills      []string       `json:"skills"`
	Benefits    []string       `json:"benefits"`
	Positions   *int           `json:"positions"`
	ExpiresAt   *time.Time     `json:"expires_at"`
}

type JobApplication struct {
	ID            uint              `json:"id" gorm:"primaryKey"`
	JobID         uint              `json:"job_id" gorm:"not null;index"`
	UserID        uint              `json:"user_id" gorm:"not null;index"`
	CompanyID     uint              `json:"company_id" gorm:"not null;index"`
	Status        ApplicationStatus `json:"status" gorm:"default:'applied'"`
	CoverLetter   string            `json:"cover_letter" gorm:"type:text"`
	ResumeURL     string            `json:"resume_url"`
	AppliedAt     time.Time         `json:"applied_at"`
	ViewedAt      *time.Time        `json:"viewed_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
	Notes         string            `json:"notes" gorm:"type:text"`
	CreatedAt     time.Time         `json:"created_at"`
}

type JobApplicationDetail struct {
	ID            uint              `json:"id"`
	JobID         uint              `json:"job_id"`
	JobTitle      string            `json:"job_title"`
	ApplicantName string            `json:"applicant_name"`
	ApplicantEmail string           `json:"applicant_email"`
	ApplicantPhone string           `json:"applicant_phone"`
	Position      string            `json:"position"`
	Location      string            `json:"location"`
	DateApplied   time.Time         `json:"date_applied"`
	Experience    string            `json:"experience"`
	Status        ApplicationStatus `json:"status"`
	ResumeURL     string            `json:"resume_url"`
	CoverLetter   string            `json:"cover_letter"`
}

type JobListResponse struct {
	ID              uint      `json:"id"`
	Title           string    `json:"title"`
	JobType         JobType   `json:"job_type"`
	Location        string    `json:"location"`
	Status          JobStatus `json:"status"`
	Positions       int       `json:"positions"`
	ApplicationCount int      `json:"application_count"`
	ViewCount       int       `json:"view_count"`
	PublishedAt     *time.Time `json:"published_at"`
	ExpiresAt       *time.Time `json:"expires_at"`
	CreatedAt       time.Time `json:"created_at"`
}