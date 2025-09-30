package repository

import (
	"jobfair-auth-service/internal/models"

	"gorm.io/gorm"
)

type JobSeekerProfileRepository struct {
	db *gorm.DB
}

func NewJobSeekerProfileRepository(db *gorm.DB) *JobSeekerProfileRepository {
	return &JobSeekerProfileRepository{db: db}
}

// ✅ Tambahkan helper untuk membersihkan string kosong
func sanitizeProfile(profile *models.JobSeekerProfile) {
	if profile.EmploymentStatus == "" {
		profile.EmploymentStatus = "" // biarkan kosong jika kolom tidak ada constraint ketat
	}

	// ✅ Kalau kosong, ubah ke nilai default NULL agar tidak melanggar constraint
	if profile.JobSearchStatus == "" {
		profile.JobSearchStatus = "" // bisa juga nil kalau kamu pakai pointer *string di model
	}
}

func (r *JobSeekerProfileRepository) Create(profile *models.JobSeekerProfile) error {
	sanitizeProfile(profile) // ✅ pastikan sebelum insert, sudah bersih
	return r.db.Create(profile).Error
}

func (r *JobSeekerProfileRepository) GetByUserID(userID uint) (*models.JobSeekerProfile, error) {
	var profile models.JobSeekerProfile
	if err := r.db.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *JobSeekerProfileRepository) Update(profile *models.JobSeekerProfile) error {
	sanitizeProfile(profile) // ✅ pastikan sebelum update juga bersih
	return r.db.Save(profile).Error
}

func (r *JobSeekerProfileRepository) Delete(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.JobSeekerProfile{}).Error
}
