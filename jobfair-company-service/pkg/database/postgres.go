package database

import (
    "jobfair-company-service/internal/models"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

func Connect(databaseURL string) (*gorm.DB, error) {
    db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })
    if err != nil {
        return nil, err
    }

    // Auto migrate tables
    err = db.AutoMigrate(&models.Company{}, &models.CompanyAnalytics{})
    if err != nil {
        return nil, err
    }

    return db, nil
}