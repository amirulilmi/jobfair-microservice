package services

import (
	"errors"
	"time"

	"jobfair-company-service/internal/models"
	"jobfair-company-service/internal/repository"
	"github.com/gosimple/slug"
)

type JobService struct {
	jobRepo         *repository.JobRepository
	companyRepo     *repository.CompanyRepository
	applicationRepo *repository.ApplicationRepository
}

func NewJobService(jobRepo *repository.JobRepository, companyRepo *repository.CompanyRepository, applicationRepo *repository.ApplicationRepository) *JobService {
	return &JobService{
		jobRepo:         jobRepo,
		companyRepo:     companyRepo,
		applicationRepo: applicationRepo,
	}
}

func (s *JobService) CreateJob(companyID uint, req *models.CreateJobRequest) (*models.Job, error) {
	company, err := s.companyRepo.GetByID(companyID)
	if err != nil {
		return nil, errors.New("company not found")
	}

	job := &models.Job{
		CompanyID:        companyID,
		Title:            req.Title,
		Description:      req.Description,
		Requirements:     req.Requirements,
		Responsibilities: req.Responsibilities,
		JobType:          req.JobType,
		JobLevel:         req.JobLevel,
		Status:           models.JobStatusDraft,
		Location:         req.Location,
		City:             req.City,
		Country:          req.Country,
		IsRemote:         req.IsRemote,
		SalaryMin:        req.SalaryMin,
		SalaryMax:        req.SalaryMax,
		ShowSalary:       req.ShowSalary,
		Skills:           req.Skills,
		Benefits:         req.Benefits,
		Positions:        req.Positions,
		ExpiresAt:        req.ExpiresAt,
		Slug:             slug.Make(req.Title + "-" + company.Name),
	}

	return s.jobRepo.Create(job)
}

func (s *JobService) GetJob(id uint) (*models.Job, error) {
	job, err := s.jobRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	s.jobRepo.IncrementViewCount(id)

	return job, nil
}

func (s *JobService) UpdateJob(id uint, companyID uint, req *models.UpdateJobRequest) (*models.Job, error) {
	job, err := s.jobRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if job.CompanyID != companyID {
		return nil, errors.New("unauthorized to update this job")
	}

	if req.Title != nil {
		job.Title = *req.Title
		company, _ := s.companyRepo.GetByID(companyID)
		if company != nil {
			job.Slug = slug.Make(*req.Title + "-" + company.Name)
		}
	}
	if req.Description != nil {
		job.Description = *req.Description
	}
	if req.Requirements != nil {
		job.Requirements = *req.Requirements
	}
	if req.Responsibilities != nil {
		job.Responsibilities = *req.Responsibilities
	}
	if req.JobType != nil {
		job.JobType = *req.JobType
	}
	if req.JobLevel != nil {
		job.JobLevel = *req.JobLevel
	}
	if req.Status != nil {
		job.Status = *req.Status
		if *req.Status == models.JobStatusActive && job.PublishedAt == nil {
			now := time.Now()
			job.PublishedAt = &now
		}
	}
	if req.Location != nil {
		job.Location = *req.Location
	}
	if req.City != nil {
		job.City = *req.City
	}
	if req.Country != nil {
		job.Country = *req.Country
	}
	if req.IsRemote != nil {
		job.IsRemote = *req.IsRemote
	}
	if req.SalaryMin != nil {
		job.SalaryMin = *req.SalaryMin
	}
	if req.SalaryMax != nil {
		job.SalaryMax = *req.SalaryMax
	}
	if req.ShowSalary != nil {
		job.ShowSalary = *req.ShowSalary
	}
	if req.Skills != nil {
		job.Skills = req.Skills
	}
	if req.Benefits != nil {
		job.Benefits = req.Benefits
	}
	if req.Positions != nil {
		job.Positions = *req.Positions
	}
	if req.ExpiresAt != nil {
		job.ExpiresAt = req.ExpiresAt
	}

	if err := s.jobRepo.Update(job); err != nil {
		return nil, err
	}

	return job, nil
}

func (s *JobService) DeleteJob(id uint, companyID uint) error {
	job, err := s.jobRepo.GetByID(id)
	if err != nil {
		return err
	}

	if job.CompanyID != companyID {
		return errors.New("unauthorized to delete this job")
	}

	return s.jobRepo.Delete(id)
}

func (s *JobService) ListJobs(companyID uint, limit, offset int, filters map[string]interface{}) ([]*models.Job, int64, error) {
	return s.jobRepo.List(companyID, limit, offset, filters)
}

func (s *JobService) PublishJob(id uint, companyID uint) error {
	job, err := s.jobRepo.GetByID(id)
	if err != nil {
		return err
	}

	if job.CompanyID != companyID {
		return errors.New("unauthorized to publish this job")
	}

	return s.jobRepo.UpdateStatus(id, models.JobStatusActive)
}

func (s *JobService) CloseJob(id uint, companyID uint) error {
	job, err := s.jobRepo.GetByID(id)
	if err != nil {
		return err
	}

	if job.CompanyID != companyID {
		return errors.New("unauthorized to close this job")
	}

	return s.jobRepo.UpdateStatus(id, models.JobStatusClosed)
}