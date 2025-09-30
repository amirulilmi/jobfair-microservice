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
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error(), "data": nil})
		return
	}

	userID := uint(1) // TODO: ambil dari JWT
	company, err := h.service.CreateCompany(userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error(), "data": nil})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Company created successfully", "data": company})
}

func (h *CompanyHandler) GetCompany(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	company, err := h.service.GetCompany(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Company not found", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Company retrieved successfully", "data": company})
}

func (h *CompanyHandler) UpdateCompany(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var req models.UpdateCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error(), "data": nil})
		return
	}

	company, err := h.service.UpdateCompany(uint(id), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error(), "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Company updated successfully", "data": company})
}

func (h *CompanyHandler) UploadFile(c *gin.Context, fileType string) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "File upload failed", "data": nil})
		return
	}

	url, err := h.service.UploadFile(uint(id), file, fileType)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error(), "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "File uploaded successfully", "data": url})
}

func (h *CompanyHandler) UploadLogo(c *gin.Context)   { h.UploadFile(c, "logo") }
func (h *CompanyHandler) UploadBanner(c *gin.Context) { h.UploadFile(c, "banner") }
func (h *CompanyHandler) UploadVideo(c *gin.Context)  { h.UploadFile(c, "video") }

func (h *CompanyHandler) GetAnalytics(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	analytics, err := h.service.GetAnalytics(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Analytics not found", "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Analytics retrieved successfully", "data": analytics})
}
