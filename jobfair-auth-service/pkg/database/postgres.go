package database

import (
	"log"
	"time"

	// "jobfair-auth-service/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// func Connect(databaseURL string) (*gorm.DB, error) {
// 	config := &gorm.Config{
// 		Logger: logger.Default.LogMode(logger.Info),
// 	}

// 	db, err := gorm.Open(postgres.Open(databaseURL), config)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Configure connection pool
// 	sqlDB, err := db.DB()
// 	if err != nil {
// 		return nil, err
// 	}

// 	sqlDB.SetMaxIdleConns(10)
// 	sqlDB.SetMaxOpenConns(100)
// 	sqlDB.SetConnMaxLifetime(time.Hour)

// 	// Test connection
// 	if err := sqlDB.Ping(); err != nil {
// 		return nil, err
// 	}

// 	log.Println("✅ Database connected successfully")

// 	// Auto migrate
// 	if err := db.AutoMigrate(
// 		&models.User{},
// 		&models.JobSeekerProfile{},
// 		&models.OTPVerification{},
// 	); err != nil {
// 		return nil, err
// 	}

// 	log.Println("✅ Database migration completed")
// 	return db, nil
// }

func Connect(databaseURL string) (*gorm.DB, error) {
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	db, err := gorm.Open(postgres.Open(databaseURL), config)
	if err != nil {
		return nil, err
	}

	// Get underlying SQL DB
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	log.Println("✅ Database connected successfully")

	// ❌ NO AUTO MIGRATION - Use migration files instead
	// Migration files are located in: services/auth-service/migrations/

	return db, nil
}
