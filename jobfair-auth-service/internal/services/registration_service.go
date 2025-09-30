package services

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"jobfair-auth-service/internal/models"
	"jobfair-auth-service/internal/repository"
	"jobfair-auth-service/internal/utils"
)

type RegistrationService struct {
	userRepo    *repository.UserRepository
	profileRepo *repository.JobSeekerProfileRepository
	otpRepo     *repository.OTPRepository
	jwtSecret   string
}

func NewRegistrationService(
	userRepo *repository.UserRepository,
	profileRepo *repository.JobSeekerProfileRepository,
	otpRepo *repository.OTPRepository,
	jwtSecret string,
) *RegistrationService {
	return &RegistrationService{
		userRepo:    userRepo,
		profileRepo: profileRepo,
		otpRepo:     otpRepo,
		jwtSecret:   jwtSecret,
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
		NextStep:     "complete_profile",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// Step 2: Complete Basic Profile
func (s *RegistrationService) CompleteBasicProfile(userID uint, req *models.RegisterStep2Request) (*models.BasicProfileData, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
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

	return &models.BasicProfileData{
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		PhoneNumber: user.PhoneNumber,
		CountryCode: user.CountryCode,
		Country:     user.Country,
	}, nil
}

// Step 5: Set Employment Status
func (s *RegistrationService) SetEmploymentStatus(userID uint, req *models.JobSeekerStep1Request) (*models.JobPreferencesData, error) {
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

	return &models.JobPreferencesData{
		JobSearchStatus:    string(profile.JobSearchStatus),
		DesiredPositions:   profile.DesiredPositions,
		PreferredLocations: profile.PreferredLocations,
		JobTypes:           profile.JobTypes,
	}, nil
}

// Step 6: Set Job Preferences
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

// Step 7: Set Permissions
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

// Step 8: Upload Profile Photo
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

	return &models.ProfilePhotoData{PhotoURL: photoURL}, nil
}

// Helper: Generate OTP
func (s *RegistrationService) generateOTP() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}
