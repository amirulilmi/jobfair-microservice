package handlers

import (
	"net/http"
	"strconv"

	"jobfair-company-service/internal/models"
	"jobfair-company-service/internal/services"

	"github.com/gin-gonic/gin"
)

type JobHandler struct {
	companyService *services.CompanyService
	jobService     *services.JobService
}

func NewJobHandler(companyService *services.CompanyService, jobService *services.JobService) *JobHandler {
	return &JobHandler{
		companyService: companyService,
		jobService:     jobService,
	}
}

func (h *JobHandler) CreateJob(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse("Unauthorized", "UNAUTHORIZED", nil))
		return
	}

	company, err := h.companyService.GetCompanyByUserID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse("Company not found", "NOT_FOUND", nil))
		return
	}

	var req models.CreateJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("Invalid request", "VALIDATION_ERROR", err.Error()))
		return
	}

	job, err := h.jobService.CreateJob(company.ID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(err.Error(), "CREATE_FAILED", nil))
		return
	}

	c.JSON(http.StatusCreated, models.SuccessResponse("Job created successfully", job))
}

func (h *JobHandler) GetJob(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("Invalid job ID", "INVALID_ID", nil))
		return
	}

	job, err := h.jobService.GetJob(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse("Job not found", "NOT_FOUND", nil))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Job retrieved successfully", job))
}

func (h *JobHandler) UpdateJob(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse("Unauthorized", "UNAUTHORIZED", nil))
		return
	}

	company, err := h.companyService.GetCompanyByUserID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse("Company not found", "NOT_FOUND", nil))
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("Invalid job ID", "INVALID_ID", nil))
		return
	}

	var req models.UpdateJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("Invalid request", "VALIDATION_ERROR", err.Error()))
		return
	}

	job, err := h.jobService.UpdateJob(uint(id), company.ID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(err.Error(), "UPDATE_FAILED", nil))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Job updated successfully", job))
}

func (h *JobHandler) DeleteJob(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse("Unauthorized", "UNAUTHORIZED", nil))
		return
	}

	company, err := h.companyService.GetCompanyByUserID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse("Company not found", "NOT_FOUND", nil))
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("Invalid job ID", "INVALID_ID", nil))
		return
	}

	if err := h.jobService.DeleteJob(uint(id), company.ID); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(err.Error(), "DELETE_FAILED", nil))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Job deleted successfully", nil))
}

func (h *JobHandler) ListJobs(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse("Unauthorized", "UNAUTHORIZED", nil))
		return
	}

	company, err := h.companyService.GetCompanyByUserID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse("Company not found", "NOT_FOUND", nil))
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	filters := make(map[string]interface{})
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}
	if jobType := c.Query("job_type"); jobType != "" {
		filters["job_type"] = jobType
	}

	jobs, total, err := h.jobService.ListJobs(company.ID, limit, offset, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse("Failed to retrieve jobs", "SERVER_ERROR", nil))
		return
	}

	pagination := models.PaginationMeta{
		Page:       page,
		Limit:      limit,
		TotalItems: total,
		TotalPages: int((total + int64(limit) - 1) / int64(limit)),
	}

	c.JSON(http.StatusOK, models.PaginatedSuccessResponse("Jobs retrieved successfully", jobs, pagination))
}

func (h *JobHandler) PublishJob(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse("Unauthorized", "UNAUTHORIZED", nil))
		return
	}

	company, err := h.companyService.GetCompanyByUserID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse("Company not found", "NOT_FOUND", nil))
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("Invalid job ID", "INVALID_ID", nil))
		return
	}

	if err := h.jobService.PublishJob(uint(id), company.ID); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(err.Error(), "PUBLISH_FAILED", nil))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Job published successfully", nil))
}

func (h *JobHandler) CloseJob(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse("Unauthorized", "UNAUTHORIZED", nil))
		return
	}

	company, err := h.companyService.GetCompanyByUserID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse("Company not found", "NOT_FOUND", nil))
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("Invalid job ID", "INVALID_ID", nil))
		return
	}

	if err := h.jobService.CloseJob(uint(id), company.ID); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(err.Error(), "CLOSE_FAILED", nil))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Job closed successfully", nil))
}