package database

import (
	"log"

	"habitflow/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	demoEmail    = "demo@habitflow.app"
	demoPassword = "demo1234"
	demoName     = "Demo User"
	bcryptCost   = 12
)

// SeedDemoAccount creates or updates the demo account in the database.
// This function is idempotent and safe to call on every startup.
func SeedDemoAccount() {
	var user models.User
	result := DB.Where("email = ?", demoEmail).First(&user)

	// Hash the demo password
	hash, err := bcrypt.GenerateFromPassword([]byte(demoPassword), bcryptCost)
	if err != nil {
		log.Printf("Warning: Failed to hash demo password: %v", err)
		return
	}

	if result.Error == gorm.ErrRecordNotFound {
		// Demo account doesn't exist, create it
		user = models.User{
			Name:         demoName,
			Email:        demoEmail,
			PasswordHash: string(hash),
		}
		if err := DB.Create(&user).Error; err != nil {
			log.Printf("Warning: Failed to create demo account: %v", err)
			return
		}
		log.Println("Demo account created successfully")
	} else if result.Error != nil {
		// Some other database error
		log.Printf("Warning: Failed to check for demo account: %v", result.Error)
		return
	} else {
		// Demo account exists, update password to ensure it's correct
		if err := DB.Model(&user).Update("password_hash", string(hash)).Error; err != nil {
			log.Printf("Warning: Failed to update demo account password: %v", err)
			return
		}
		log.Println("Demo account verified and password reset")
	}
}
