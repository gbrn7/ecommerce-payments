package helpers

import (
	"ecommerce-payments/internal/models"
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func SetupPostgreSQL() {
	var err error
	dbuser := GetEnv("DB_USER", "")
	dbpass := GetEnv("DB_PASSWORD", "")
	dbhost := GetEnv("DB_HOST", "127.0.0.1")
	dbport := GetEnv("DB_PORT", "5432")
	dbname := GetEnv("DB_NAME", "")

	dsn := fmt.Sprintf("host=%s  port=%s user=%s password=%s dbname=%s sslmode=disable", dbhost, dbport, dbuser, dbpass, dbname)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Fatal("failed to connect to database", err)
	}

	logrus.Info("successfully connect to database")

	DB.AutoMigrate(&models.PaymentTransaction{}, &models.PaymentMethod{})
}
