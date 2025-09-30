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
)

type CompanyService struct {
	companyRepo *repository.CompanyRepository
}

func NewCompanyService(repo *repository.CompanyRepository) *CompanyService {
	return &CompanyService{companyRepo: repo}
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
		Location:    req.Location,
		Website:     req.Website,
		Email:       req.Email,
		Phone:       req.Phone,
		IsVerified:  false,
	}

	return s.companyRepo.Create(company)
}

func (s *CompanyService) GetCompany(id uint) (*models.Company, error) {
	return s.companyRepo.GetByID(id)
}

func (s *CompanyService) UpdateCompany(id uint, req *models.UpdateCompanyRequest) (*models.Company, error) {
	company, err := s.companyRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		company.Name = req.Name
	}
	if req.Description != "" {
		company.Description = req.Description
	}
	if req.Industry != "" {
		company.Industry = req.Industry
	}
	if req.Location != "" {
		company.Location = req.Location
	}
	if req.Website != "" {
		company.Website = req.Website
	}
	if req.Email != "" {
		company.Email = req.Email
	}
	if req.Phone != "" {
		company.Phone = req.Phone
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
	}

	if err := s.companyRepo.Update(company); err != nil {
		return "", err
	}

	return url, nil
}

func (s *CompanyService) validateFile(file *multipart.FileHeader, fileType string) error {
	if file.Size > 5*1024*1024 {
		return errors.New("file size too large (max 5MB)")
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	var allowed []string

	switch fileType {
	case "logo", "banner":
		allowed = []string{".jpg", ".jpeg", ".png", ".gif"}
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
