package handlers

import (
	"net/http"
	"strconv"

	"jobfair-company-service/internal/models"
	"jobfair-company-service/internal/services"

	"github.com/gin-gonic/gin"
)

type CompanyHandler struct {
	service *services.CompanyService
}

func NewCompanyHandler(service *services.CompanyService) *CompanyHandler {
	return &CompanyHandler{service: service}
}

func (h *CompanyHandler) CreateCompany(c *gin.Context) {
	var req models.CreateCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("Invalid request", "VALIDATION_ERROR", err.Error()))
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse("Unauthorized", "UNAUTHORIZED", nil))
		return
	}

	company, err := h.service.CreateCompany(userID.(uint), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(err.Error(), "CREATE_FAILED", nil))
		return
	}

	c.JSON(http.StatusCreated, models.SuccessResponse("Company created successfully", company))
}

func (h *CompanyHandler) GetCompany(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("Invalid company ID", "INVALID_ID", nil))
		return
	}

	company, err := h.service.GetCompany(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse("Company not found", "NOT_FOUND", nil))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Company retrieved successfully", company))
}

func (h *CompanyHandler) GetMyCompany(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse("Unauthorized", "UNAUTHORIZED", nil))
		return
	}

	company, err := h.service.GetCompanyByUserID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse("Company not found", "NOT_FOUND", nil))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Company retrieved successfully", company))
}

func (h *CompanyHandler) UpdateCompany(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("Invalid company ID", "INVALID_ID", nil))
		return
	}

	var req models.UpdateCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("Invalid request", "VALIDATION_ERROR", err.Error()))
		return
	}

	company, err := h.service.UpdateCompany(uint(id), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(err.Error(), "UPDATE_FAILED", nil))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Company updated successfully", company))
}

func (h *CompanyHandler) UploadFile(c *gin.Context, fileType string) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("Invalid company ID", "INVALID_ID", nil))
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("File upload failed", "UPLOAD_FAILED", err.Error()))
		return
	}

	url, err := h.service.UploadFile(uint(id), file, fileType)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(err.Error(), "UPLOAD_FAILED", nil))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("File uploaded successfully", gin.H{"url": url}))
}

func (h *CompanyHandler) UploadLogo(c *gin.Context)    { h.UploadFile(c, "logo") }
func (h *CompanyHandler) UploadBanner(c *gin.Context)  { h.UploadFile(c, "banner") }
func (h *CompanyHandler) UploadVideo(c *gin.Context)   { h.UploadFile(c, "video") }
func (h *CompanyHandler) UploadGallery(c *gin.Context) { h.UploadFile(c, "gallery") }

func (h *CompanyHandler) GetAnalytics(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse("Invalid company ID", "INVALID_ID", nil))
		return
	}

	analytics, err := h.service.GetAnalytics(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse("Analytics not found", "NOT_FOUND", nil))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Analytics retrieved successfully", analytics))
}

func (h *CompanyHandler) GetDashboard(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse("Unauthorized", "UNAUTHORIZED", nil))
		return
	}

	company, err := h.service.GetCompanyByUserID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse("Company not found", "NOT_FOUND", nil))
		return
	}

	stats, err := h.service.GetDashboardStats(company.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse("Failed to retrieve dashboard stats", "SERVER_ERROR", nil))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("Dashboard stats retrieved successfully", stats))
}

func (h *CompanyHandler) ListCompanies(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	filters := make(map[string]interface{})
	if industry := c.Query("industry"); industry != "" {
		filters["industry"] = industry
	}
	if city := c.Query("city"); city != "" {
		filters["city"] = city
	}

	companies, total, err := h.service.ListCompanies(limit, offset, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse("Failed to retrieve companies", "SERVER_ERROR", nil))
		return
	}

	pagination := models.PaginationMeta{
		Page:       page,
		Limit:      limit,
		TotalItems: total,
		TotalPages: int((total + int64(limit) - 1) / int64(limit)),
	}

	c.JSON(http.StatusOK, models.PaginatedSuccessResponse("Companies retrieved successfully", companies, pagination))
}
