package services

import (
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"jobfair-company-service/internal/models"
	"jobfair-company-service/internal/repository"
	"github.com/gosimple/slug"
)

type CompanyService struct {
	companyRepo     *repository.CompanyRepository
	jobRepo         *repository.JobRepository
	applicationRepo *repository.ApplicationRepository
}

func NewCompanyService(companyRepo *repository.CompanyRepository, jobRepo *repository.JobRepository, applicationRepo *repository.ApplicationRepository) *CompanyService {
	return &CompanyService{
		companyRepo:     companyRepo,
		jobRepo:         jobRepo,
		applicationRepo: applicationRepo,
	}
}

func (s *CompanyService) CreateCompany(userID uint, req *models.CreateCompanyRequest) (*models.Company, error) {
	if existing, _ := s.companyRepo.GetByUserID(userID); existing != nil {
		return nil, errors.New("company already exists for this user")
	}

	company := &models.Company{
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
		Industry:    req.Industry,
		CompanySize: req.CompanySize,
		FoundedYear: req.FoundedYear,
		Email:       req.Email,
		Phone:       req.Phone,
		Website:     req.Website,
		Address:     req.Address,
		City:        req.City,
		Country:     req.Country,
		Slug:        slug.Make(req.Name),
		IsVerified:  false,
	}

	createdCompany, err := s.companyRepo.Create(company)
	if err != nil {
		return nil, err
	}

	analytics := &models.CompanyAnalytics{
		CompanyID: createdCompany.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	s.companyRepo.UpdateAnalytics(analytics)

	return createdCompany, nil
}

func (s *CompanyService) GetCompany(id uint) (*models.Company, error) {
	company, err := s.companyRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	s.companyRepo.IncrementProfileViews(id)

	return company, nil
}

func (s *CompanyService) GetCompanyByUserID(userID uint) (*models.Company, error) {
	return s.companyRepo.GetByUserID(userID)
}

func (s *CompanyService) UpdateCompany(id uint, req *models.UpdateCompanyRequest) (*models.Company, error) {
	company, err := s.companyRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		company.Name = *req.Name
		company.Slug = slug.Make(*req.Name)
	}
	if req.Description != nil {
		company.Description = *req.Description
	}
	if req.Industry != nil {
		company.Industry = *req.Industry
	}
	if req.CompanySize != nil {
		company.CompanySize = *req.CompanySize
	}
	if req.FoundedYear != nil {
		company.FoundedYear = *req.FoundedYear
	}
	if req.Email != nil {
		company.Email = *req.Email
	}
	if req.Phone != nil {
		company.Phone = *req.Phone
	}
	if req.Website != nil {
		company.Website = *req.Website
	}
	if req.Address != nil {
		company.Address = *req.Address
	}
	if req.City != nil {
		company.City = *req.City
	}
	if req.State != nil {
		company.State = *req.State
	}
	if req.Country != nil {
		company.Country = *req.Country
	}
	if req.PostalCode != nil {
		company.PostalCode = *req.PostalCode
	}
	if req.LinkedinURL != nil {
		company.LinkedinURL = *req.LinkedinURL
	}
	if req.FacebookURL != nil {
		company.FacebookURL = *req.FacebookURL
	}
	if req.TwitterURL != nil {
		company.TwitterURL = *req.TwitterURL
	}
	if req.InstagramURL != nil {
		company.InstagramURL = *req.InstagramURL
	}

	if err := s.companyRepo.Update(company); err != nil {
		return nil, err
	}

	return company, nil
}

func (s *CompanyService) UploadFile(companyID uint, file *multipart.FileHeader, fileType string) (string, error) {
	company, err := s.companyRepo.GetByID(companyID)
	if err != nil {
		return "", err
	}

	if err := s.validateFile(file, fileType); err != nil {
		return "", err
	}

	filename := fmt.Sprintf("%d_%s_%d%s", companyID, fileType, time.Now().Unix(), filepath.Ext(file.Filename))
	url := fmt.Sprintf("/uploads/%s", filename)

	switch fileType {
	case "logo":
		company.LogoURL = url
	case "banner":
		company.BannerURL = url
	case "video":
		company.VideoURLs = append(company.VideoURLs, url)
	case "gallery":
		company.GalleryURLs = append(company.GalleryURLs, url)
	}

	if err := s.companyRepo.Update(company); err != nil {
		return "", err
	}

	return url, nil
}

func (s *CompanyService) validateFile(file *multipart.FileHeader, fileType string) error {
	maxSize := int64(10 * 1024 * 1024)
	if fileType == "video" {
		maxSize = int64(50 * 1024 * 1024)
	}

	if file.Size > maxSize {
		return fmt.Errorf("file size too large (max %dMB)", maxSize/(1024*1024))
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	var allowed []string

	switch fileType {
	case "logo", "banner", "gallery":
		allowed = []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	case "video":
		allowed = []string{".mp4", ".avi", ".mov", ".webm"}
	default:
		return errors.New("invalid file type")
	}

	for _, a := range allowed {
		if ext == a {
			return nil
		}
	}
	return errors.New("invalid file format")
}

func (s *CompanyService) GetAnalytics(companyID uint) (*models.CompanyAnalytics, error) {
	return s.companyRepo.GetAnalytics(companyID)
}

func (s *CompanyService) GetDashboardStats(companyID uint) (*models.DashboardStats, error) {
	analytics, err := s.companyRepo.GetAnalytics(companyID)
	if err != nil {
		return nil, err
	}

	totalApplicants, _ := s.applicationRepo.CountByCompanyID(companyID)
	totalJobs, _ := s.jobRepo.GetTotalJobsCount(companyID)
	activeJobs, _ := s.jobRepo.GetActiveJobsCount(companyID)

	shortlisted, _ := s.applicationRepo.CountByStatus(companyID, models.ApplicationStatusShortlisted)
	hired, _ := s.applicationRepo.CountByStatus(companyID, models.ApplicationStatusHired)
	interview, _ := s.applicationRepo.CountByStatus(companyID, models.ApplicationStatusInterview)
	rejected, _ := s.applicationRepo.CountByStatus(companyID, models.ApplicationStatusRejected)

	stats := &models.DashboardStats{
		TotalJobsPosted:  int(totalJobs),
		TotalApplicants:  int(totalApplicants),
		JobViews:         analytics.JobViews,
		TotalPositions:   int(activeJobs),
		JobsGrowth:       5.5,
		ApplicantsGrowth: 3.8,
		ViewsGrowth:      7.2,
		PositionsGrowth:  4.6,
		ApplicantsByStatus: map[string]int{
			"shortlisted": int(shortlisted),
			"hired":       int(hired),
			"interview":   int(interview),
			"rejected":    int(rejected),
		},
	}

	return stats, nil
}

func (s *CompanyService) ListCompanies(limit, offset int, filters map[string]interface{}) ([]*models.Company, int64, error) {
	return s.companyRepo.List(limit, offset, filters)
}
