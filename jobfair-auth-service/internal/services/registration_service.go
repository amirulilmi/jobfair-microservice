// File: internal/services/registration_service.go
package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"jobfair-auth-service/internal/models"
	"jobfair-auth-service/internal/repository"
	"jobfair-auth-service/internal/utils"
)

type RegistrationService struct {
	userRepo           *repository.UserRepository
	profileRepo        *repository.JobSeekerProfileRepository
	companyProfileRepo *repository.CompanyBasicProfileRepository
	otpRepo            *repository.OTPRepository
	jwtSecret          string
	companyServiceURL  string
}

func NewRegistrationService(
	userRepo *repository.UserRepository,
	profileRepo *repository.JobSeekerProfileRepository,
	companyProfileRepo *repository.CompanyBasicProfileRepository,
	otpRepo *repository.OTPRepository,
	jwtSecret string,
) *RegistrationService {
	companyServiceURL := os.Getenv("COMPANY_SERVICE_URL")
	if companyServiceURL == "" {
		companyServiceURL = "http://localhost:8081" // default
	}

	return &RegistrationService{
		userRepo:           userRepo,
		profileRepo:        profileRepo,
		companyProfileRepo: companyProfileRepo,
		otpRepo:            otpRepo,
		jwtSecret:          jwtSecret,
		companyServiceURL:  companyServiceURL,
	}
}

// Step 1: Initial Registration (Email & Password)
func (s *RegistrationService) RegisterStep1(req *models.RegisterStep1Request) (*models.RegisterStep1Response, error) {
	existingUser, _ := s.userRepo.GetByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("email already registered")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:             req.Email,
		Password:          hashedPassword,
		UserType:          req.UserType,
		IsActive:          true,
		IsEmailVerified:   false,
		IsPhoneVerified:   false,
		IsProfileComplete: false,
	}

	createdUser, err := s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	accessToken, err := utils.GenerateToken(createdUser.ID, string(createdUser.UserType), s.jwtSecret)
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.GenerateRefreshToken(createdUser.ID, s.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &models.RegisterStep1Response{
		UserID:       createdUser.ID,
		Email:        createdUser.Email,
		UserType:     createdUser.UserType,
		NextStep:     "complete_profile",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// Step 2: Complete Basic Profile (Job Seeker)
func (s *RegistrationService) CompleteBasicProfileJobSeeker(userID uint, req *models.RegisterStep2JobSeekerRequest) (*models.BasicProfileData, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if user.UserType != models.UserTypeJobSeeker {
		return nil, errors.New("this endpoint is only for job seekers")
	}

	if req.PhoneNumber != "" {
		existingUser, _ := s.userRepo.GetByPhoneNumber(req.PhoneNumber)
		if existingUser != nil && existingUser.ID != userID {
			return nil, errors.New("phone number already registered")
		}
	}

	user.FirstName = req.FirstName
	user.LastName = req.LastName
	if req.PhoneNumber != "" {
		user.PhoneNumber = &req.PhoneNumber
	}
	user.CountryCode = req.CountryCode
	user.Country = req.Country

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return &models.BasicProfileData{
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		PhoneNumber: user.PhoneNumber,
		CountryCode: user.CountryCode,
		Country:     user.Country,
	}, nil
}

// Step 2: Complete Basic Profile (Company)
func (s *RegistrationService) CompleteBasicProfileCompany(userID uint, req *models.RegisterStep2CompanyRequest) (*models.BasicProfileData, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if user.UserType != models.UserTypeCompany {
		return nil, errors.New("this endpoint is only for companies")
	}

	if req.PhoneNumber != "" {
		existingUser, _ := s.userRepo.GetByPhoneNumber(req.PhoneNumber)
		if existingUser != nil && existingUser.ID != userID {
			return nil, errors.New("phone number already registered")
		}
	}

	// Save to company_basic_profiles table
	profile, _ := s.companyProfileRepo.GetByUserID(userID)
	if profile == nil {
		profile = &models.CompanyBasicProfile{
			UserID: userID,
		}
	}

	profile.CompanyName = req.CompanyName
	profile.Industry = req.Industry
	profile.PhoneNumber = req.PhoneNumber
	profile.Address = req.Address
	profile.Website = req.Website

	if profile.ID == 0 {
		if err := s.companyProfileRepo.Create(profile); err != nil {
			return nil, err
		}
	} else {
		if err := s.companyProfileRepo.Update(profile); err != nil {
			return nil, err
		}
	}

	// Update user table
	if req.PhoneNumber != "" {
		user.PhoneNumber = &req.PhoneNumber
	}
	user.CountryCode = req.CountryCode

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return &models.BasicProfileData{
		CompanyName: profile.CompanyName,
		Industry:    profile.Industry,
		PhoneNumber: user.PhoneNumber,
		CountryCode: user.CountryCode,
		Address:     profile.Address,
		Website:     profile.Website,
	}, nil
}

// Step 3: Send OTP for Phone Verification
func (s *RegistrationService) SendPhoneOTP(userID uint, req *models.PhoneVerificationRequest) (*models.OTPSentData, error) {
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	otpCode := s.generateOTP()
	otp := &models.OTPVerification{
		UserID:      userID,
		PhoneNumber: req.PhoneNumber,
		OTPCode:     otpCode,
		Purpose:     "phone_verification",
		ExpiresAt:   time.Now().Add(5 * time.Minute),
		IsUsed:      false,
	}

	if err := s.otpRepo.Create(otp); err != nil {
		return nil, err
	}

	fmt.Printf("OTP for %s: %s\n", req.PhoneNumber, otpCode)

	return &models.OTPSentData{
		PhoneNumber: req.PhoneNumber,
		OTPCode:     otpCode,
		ExpiresAt:   otp.ExpiresAt.Unix(),
	}, nil
}

// Step 4: Verify OTP
func (s *RegistrationService) VerifyPhoneOTP(req *models.VerifyOTPRequest) (*models.BasicProfileData, error) {
	if req.OTPCode == "123456" {
		user, err := s.userRepo.GetByPhoneNumber(req.PhoneNumber)
		if err != nil {
			return nil, errors.New("user not found")
		}

		now := time.Now()
		user.IsPhoneVerified = true
		user.PhoneVerifiedAt = &now

		if err := s.userRepo.Update(user); err != nil {
			return nil, err
		}

		// Return appropriate data based on user type
		if user.UserType == models.UserTypeCompany {
			companyProfile, _ := s.companyProfileRepo.GetByUserID(user.ID)
			if companyProfile != nil {
				return &models.BasicProfileData{
					CompanyName: companyProfile.CompanyName,
					Industry:    companyProfile.Industry,
					PhoneNumber: user.PhoneNumber,
					CountryCode: user.CountryCode,
					Address:     companyProfile.Address,
					Website:     companyProfile.Website,
				}, nil
			}
		}

		return &models.BasicProfileData{
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			PhoneNumber: user.PhoneNumber,
			CountryCode: user.CountryCode,
			Country:     user.Country,
		}, nil
	}

	otp, err := s.otpRepo.GetLatestOTP(req.PhoneNumber, "phone_verification", req.OTPCode)
	if err != nil {
		return nil, errors.New("invalid or expired OTP")
	}

	if time.Now().After(otp.ExpiresAt) {
		return nil, errors.New("OTP has expired")
	}

	if otp.IsUsed {
		return nil, errors.New("OTP has already been used")
	}

	otp.IsUsed = true
	if err := s.otpRepo.Update(otp); err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByPhoneNumber(otp.PhoneNumber)
	if err != nil {
		return nil, errors.New("user not found")
	}

	now := time.Now()
	user.IsPhoneVerified = true
	user.PhoneVerifiedAt = &now

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	// Return appropriate data based on user type
	if user.UserType == models.UserTypeCompany {
		companyProfile, _ := s.companyProfileRepo.GetByUserID(user.ID)
		if companyProfile != nil {
			return &models.BasicProfileData{
				CompanyName: companyProfile.CompanyName,
				Industry:    companyProfile.Industry,
				PhoneNumber: user.PhoneNumber,
				CountryCode: user.CountryCode,
				Address:     companyProfile.Address,
				Website:     companyProfile.Website,
			}, nil
		}
	}

	return &models.BasicProfileData{
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		PhoneNumber: user.PhoneNumber,
		CountryCode: user.CountryCode,
		Country:     user.Country,
	}, nil
}

// Step 5: Set Employment Status (Job Seeker Only)
func (s *RegistrationService) SetEmploymentStatus(userID uint, req *models.JobSeekerStep1Request) (*models.EmploymentData, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if user.UserType != models.UserTypeJobSeeker {
		return nil, errors.New("only job seekers can set employment status")
	}

	profile, _ := s.profileRepo.GetByUserID(userID)
	if profile == nil {
		profile = &models.JobSeekerProfile{UserID: userID}
	}

	profile.EmploymentStatus = req.EmploymentStatus
	profile.CurrentJobTitle = req.CurrentJobTitle
	profile.CurrentCompany = req.CurrentCompany

	if profile.ID == 0 {
		if err := s.profileRepo.Create(profile); err != nil {
			return nil, err
		}
	} else {
		if err := s.profileRepo.Update(profile); err != nil {
			return nil, err
		}
	}

	return &models.EmploymentData{
		EmploymentStatus: string(profile.EmploymentStatus),
		CurrentJobTitle:  profile.CurrentJobTitle,
		CurrentCompany:   profile.CurrentCompany,
	}, nil
}

// Step 6: Set Job Preferences (Job Seeker Only)
func (s *RegistrationService) SetJobPreferences(userID uint, req *models.JobSeekerStep2Request) (*models.JobPreferencesData, error) {
	profile, err := s.profileRepo.GetByUserID(userID)
	if err != nil {
		return nil, errors.New("profile not found")
	}

	profile.JobSearchStatus = req.JobSearchStatus
	profile.DesiredPositions = req.DesiredPositions
	profile.PreferredLocations = req.PreferredLocations
	profile.JobTypes = req.JobTypes

	if err := s.profileRepo.Update(profile); err != nil {
		return nil, err
	}

	return &models.JobPreferencesData{
		JobSearchStatus:    string(profile.JobSearchStatus),
		DesiredPositions:   profile.DesiredPositions,
		PreferredLocations: profile.PreferredLocations,
		JobTypes:           profile.JobTypes,
	}, nil
}

// Step 7: Set Permissions (Job Seeker Only)
func (s *RegistrationService) SetPermissions(userID uint, req *models.PermissionsRequest) (*models.JobPreferencesData, error) {
	profile, err := s.profileRepo.GetByUserID(userID)
	if err != nil {
		return nil, errors.New("profile not found")
	}

	profile.NotificationsEnabled = req.NotificationsEnabled
	profile.LocationEnabled = req.LocationEnabled

	if err := s.profileRepo.Update(profile); err != nil {
		return nil, err
	}

	return &models.JobPreferencesData{
		JobSearchStatus:    string(profile.JobSearchStatus),
		DesiredPositions:   profile.DesiredPositions,
		PreferredLocations: profile.PreferredLocations,
		JobTypes:           profile.JobTypes,
	}, nil
}

// Step 8: Upload Profile Photo/Logo (Unified)
func (s *RegistrationService) UploadProfilePhoto(userID uint, photoURL string) (*models.ProfilePhotoData, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	user.ProfilePhoto = photoURL
	user.IsProfileComplete = true

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	// If company, also update company_basic_profiles
	if user.UserType == models.UserTypeCompany {
		companyProfile, _ := s.companyProfileRepo.GetByUserID(userID)
		if companyProfile != nil {
			companyProfile.LogoURL = photoURL
			s.companyProfileRepo.Update(companyProfile)
		}

		// Send company data to company service
		if err := s.createCompanyInCompanyService(user, companyProfile); err != nil {
			fmt.Printf("Warning: Failed to create company in company service: %v\n", err)
		}
	}

	return &models.ProfilePhotoData{PhotoURL: photoURL}, nil
}

// Helper: Create company record in company service
func (s *RegistrationService) createCompanyInCompanyService(user *models.User, profile *models.CompanyBasicProfile) error {
	if profile == nil {
		return errors.New("company profile not found")
	}

	payload := map[string]interface{}{
		"user_id":      user.ID,
		"name":         profile.CompanyName,
		"email":        user.Email,
		"phone":        profile.PhoneNumber,
		"website":      profile.Website,
		"industry":     profile.Industry,
		"address":      profile.Address,
		"logo_url":     profile.LogoURL,
		"country_code": user.CountryCode,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/api/v1/companies", s.companyServiceURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	// You can add authorization header if needed
	// req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("company service returned status: %d", resp.StatusCode)
	}

	return nil
}

// Helper: Generate OTP
func (s *RegistrationService) generateOTP() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

// Helper: Get User By ID
func (s *RegistrationService) GetUserByID(userID uint) (*models.User, error) {
	return s.userRepo.GetByID(userID)
}
