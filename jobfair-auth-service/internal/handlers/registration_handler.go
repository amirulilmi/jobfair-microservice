package handlers

import (
	"fmt"
	"net/http"
	"time"

	"jobfair-auth-service/internal/models"
	"jobfair-auth-service/internal/services"

	"github.com/gin-gonic/gin"
)

type RegistrationHandler struct {
	registrationService *services.RegistrationService
}

func NewRegistrationHandler(registrationService *services.RegistrationService) *RegistrationHandler {
	return &RegistrationHandler{registrationService: registrationService}
}

func (h *RegistrationHandler) RegisterStep1(c *gin.Context) {
	var req models.RegisterStep1Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Message: err.Error()})
		return
	}

	data, err := h.registrationService.RegisterStep1(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, models.APIResponse{Data: data, Message: "Registration step 1 completed", Success: true})
}

func (h *RegistrationHandler) CompleteBasicProfile(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req models.RegisterStep2Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Message: err.Error()})
		return
	}

	data, err := h.registrationService.CompleteBasicProfile(userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{Data: data, Message: "Basic profile completed successfully", Success: true})
}

func (h *RegistrationHandler) SendPhoneOTP(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req models.PhoneVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Message: err.Error()})
		return
	}

	data, err := h.registrationService.SendPhoneOTP(userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{Data: data, Message: "OTP sent successfully", Success: true})
}

func (h *RegistrationHandler) VerifyPhoneOTP(c *gin.Context) {
	var req models.VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Message: err.Error()})
		return
	}

	data, err := h.registrationService.VerifyPhoneOTP(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{Data: data, Message: "Phone verified successfully", Success: true})
}

func (h *RegistrationHandler) SetEmploymentStatus(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req models.JobSeekerStep1Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Message: err.Error()})
		return
	}

	data, err := h.registrationService.SetEmploymentStatus(userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{Data: data, Message: "Employment status saved", Success: true})
}

func (h *RegistrationHandler) SetJobPreferences(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req models.JobSeekerStep2Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Message: err.Error()})
		return
	}

	data, err := h.registrationService.SetJobPreferences(userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{Data: data, Message: "Job preferences saved", Success: true})
}

func (h *RegistrationHandler) SetPermissions(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req models.PermissionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Message: err.Error()})
		return
	}

	data, err := h.registrationService.SetPermissions(userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{Data: data, Message: "Permissions updated", Success: true})
}

func (h *RegistrationHandler) UploadProfilePhoto(c *gin.Context) {
	userID := c.GetUint("user_id")

	file, err := c.FormFile("photo")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Message: "photo file is required"})
		return
	}

	if file.Size > 5*1024*1024 {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Message: "file size too large (max 5MB)"})
		return
	}

	photoURL := fmt.Sprintf("/uploads/profiles/%d_%d.jpg", userID, time.Now().Unix())

	data, err := h.registrationService.UploadProfilePhoto(userID, photoURL)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{Data: data, Message: "Profile photo uploaded", Success: true})
}
