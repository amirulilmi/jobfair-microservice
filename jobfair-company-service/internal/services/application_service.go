package services

import (
	"errors"

	"jobfair-company-service/internal/models"
	"jobfair-company-service/internal/repository"
)

type ApplicationService struct {
	applicationRepo *repository.ApplicationRepository
	jobRepo         *repository.JobRepository
}

func NewApplicationService(applicationRepo *repository.ApplicationRepository, jobRepo *repository.JobRepository) *ApplicationService {
	return &ApplicationService{
		applicationRepo: applicationRepo,
		jobRepo:         jobRepo,
	}
}

func (s *ApplicationService) GetApplication(id uint) (*models.JobApplication, error) {
	return s.applicationRepo.GetByID(id)
}

func (s *ApplicationService) ListApplications(companyID uint, limit, offset int, filters map[string]interface{}) ([]*models.JobApplication, int64, error) {
	return s.applicationRepo.List(companyID, limit, offset, filters)
}

func (s *ApplicationService) GetApplicationsByJobID(jobID uint, companyID uint) ([]*models.JobApplication, error) {
	job, err := s.jobRepo.GetByID(jobID)
	if err != nil {
		return nil, err
	}

	if job.CompanyID != companyID {
		return nil, errors.New("unauthorized to view applications for this job")
	}

	return s.applicationRepo.GetByJobID(jobID)
}

func (s *ApplicationService) UpdateApplicationStatus(id uint, companyID uint, status models.ApplicationStatus) error {
	application, err := s.applicationRepo.GetByID(id)
	if err != nil {
		return err
	}

	if application.CompanyID != companyID {
		return errors.New("unauthorized to update this application")
	}

	return s.applicationRepo.UpdateStatus(id, status)
}

func (s *ApplicationService) GetApplicationStats(companyID uint) (map[string]interface{}, error) {
	total, _ := s.applicationRepo.CountByCompanyID(companyID)
	shortlisted, _ := s.applicationRepo.CountByStatus(companyID, models.ApplicationStatusShortlisted)
	hired, _ := s.applicationRepo.CountByStatus(companyID, models.ApplicationStatusHired)
	interview, _ := s.applicationRepo.CountByStatus(companyID, models.ApplicationStatusInterview)
	rejected, _ := s.applicationRepo.CountByStatus(companyID, models.ApplicationStatusRejected)

	stats := map[string]interface{}{
		"total":       total,
		"shortlisted": shortlisted,
		"hired":       hired,
		"interview":   interview,
		"rejected":    rejected,
	}

	return stats, nil
}