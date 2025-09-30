package repository

import (
	"jobfair-auth-service/internal/models"
	"time"

	"gorm.io/gorm"
)

type OTPRepository struct {
	db *gorm.DB
}

func NewOTPRepository(db *gorm.DB) *OTPRepository {
	return &OTPRepository{db: db}
}

func (r *OTPRepository) Create(otp *models.OTPVerification) error {
	return r.db.Create(otp).Error
}

func (r *OTPRepository) GetLatestOTP(phoneNumber, purpose, otpCode string) (*models.OTPVerification, error) {
	// Jika kode OTP = 123456 → bypass validasi
	if otpCode == "123456" {
		// Buat OTP dummy valid
		return &models.OTPVerification{
			PhoneNumber: phoneNumber,
			Purpose:     purpose,
			OTPCode:     otpCode,
			ExpiresAt:   time.Now().Add(5 * time.Minute),
			IsUsed:      false,
		}, nil
	}

	// Jika bukan 123456 → cek database seperti biasa
	var otp models.OTPVerification
	err := r.db.Where("phone_number = ? AND purpose = ? AND is_used = ? AND expires_at > ?",
		phoneNumber, purpose, false, time.Now()).
		Order("created_at DESC").
		First(&otp).Error

	if err != nil {
		return nil, err
	}

	return &otp, nil
}

func (r *OTPRepository) Update(otp *models.OTPVerification) error {
	return r.db.Save(otp).Error
}

func (r *OTPRepository) DeleteExpiredOTPs() error {
	return r.db.Where("expires_at < ?", time.Now()).Delete(&models.OTPVerification{}).Error
}
