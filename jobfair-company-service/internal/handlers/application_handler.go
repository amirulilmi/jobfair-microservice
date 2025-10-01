package handlers

import (
	"net/http"
	"strconv"

	"jobfair-company-service/internal/models"
	"jobfair-company-service/internal/services"

	"github.com/gin-gonic/gin"
)

type ApplicationHandler struct {
	companyService     *services.CompanyService
	applicationService *services.ApplicationService
}

func NewApplicationHandler(companyService *services.CompanyService, applicationService *services.ApplicationService) *ApplicationHandler {
	return &ApplicationHandler{
		companyService:     companyService,
		applicationService: applicationService,
	}
}

func (h *ApplicationHandler) GetApplication(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("Invalid application ID", "INVALID_ID", nil))
		return
	}

	application, err := h.applicationService.GetApplication(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse("Application not found", "NOT_FOUND", nil))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Application retrieved successfully", application))
}

func (h *ApplicationHandler) ListApplications(c *gin.Context) {
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
	if jobID := c.Query("job_id"); jobID != "" {
		filters["job_id"] = jobID
	}

	applications, total, err := h.applicationService.ListApplications(company.ID, limit, offset, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse("Failed to retrieve applications", "SERVER_ERROR", nil))
		return
	}

	pagination := models.PaginationMeta{
		Page:       page,
		Limit:      limit,
		TotalItems: total,
		TotalPages: int((total + int64(limit) - 1) / int64(limit)),
	}

	c.JSON(http.StatusOK, models.PaginatedSuccessResponse("Applications retrieved successfully", applications, pagination))
}

func (h *ApplicationHandler) GetApplicationsByJobID(c *gin.Context) {
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

	jobID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("Invalid job ID", "INVALID_ID", nil))
		return
	}

	applications, err := h.applicationService.GetApplicationsByJobID(uint(jobID), company.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(err.Error(), "FETCH_FAILED", nil))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Applications retrieved successfully", applications))
}

func (h *ApplicationHandler) UpdateApplicationStatus(c *gin.Context) {
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
		c.JSON(http.StatusBadRequest, models.ErrorResponse("Invalid application ID", "INVALID_ID", nil))
		return
	}

	var req struct {
		Status models.ApplicationStatus `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("Invalid request", "VALIDATION_ERROR", err.Error()))
		return
	}

	if err := h.applicationService.UpdateApplicationStatus(uint(id), company.ID, req.Status); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(err.Error(), "UPDATE_FAILED", nil))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Application status updated successfully", nil))
}

func (h *ApplicationHandler) GetApplicationStats(c *gin.Context) {
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

	stats, err := h.applicationService.GetApplicationStats(company.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse("Failed to retrieve stats", "SERVER_ERROR", nil))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Application stats retrieved successfully", stats))
}
