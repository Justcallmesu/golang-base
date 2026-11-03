package seeds

import (
	"errors"
	"os"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"justcallmesu.com/golang-base/internal/app/models"
)

type UserSeederModel struct {
	*models.User
	TargetRole string
}

var Users = []*UserSeederModel{
	{
		User: &models.User{
			Username: "admin",
			Password: "youradminpasswordmaybe",
			FullName: "Administrator",
		},
		TargetRole: "Admin",
	},
}

func UserSeeder(database *gorm.DB, logger *log.Logger) {

	logger.Info("Running User Seeders...")

	transactionError := database.Transaction(func(tx *gorm.DB) error {

		for _, userSeeder := range Users {
			logger.Warnf("Seeding user %s", userSeeder.Username)

			var user models.User = *userSeeder.User

			if user.Username == "admin" {
				password, err := HandleAdminPassword(user.Password)
				if err != nil {
					logger.Errorf("Error handling admin password: %v", err)
					continue
				}
				user.Password = password
			}

			if err := user.SetPassword(user.Password); err != nil {
				logger.Errorf("Error hashing password for %s: %v", user.Username, err)
				continue
			}

			if err := tx.Model(&models.User{}).Where("username = ?", user.Username).First(&user).Take(&user).Error; err != nil && err != gorm.ErrRecordNotFound {
				logger.Errorf("Error finding user %s: %v", userSeeder.Username, err.Error())
			}

			if user.ID.UUID != uuid.Nil {
				logger.Warnf("User %s already exists, skipping seeding", userSeeder.Username)
				continue
			} else {
				logger.Warnf("Creating user %s", userSeeder.Username)
				if err := tx.Create(&user).Error; err != nil {
					logger.Errorf("Error creating user %s: %v", userSeeder.Username, err)
					continue
				}
			}

			if userSeeder.TargetRole == "" {
				logger.Warnf("No target role specified for user %s, skipping role assignment", userSeeder.Username)
				continue
			}

			var role models.Role
			if err := tx.Where("name = ?", userSeeder.TargetRole).First(&role).Error; err != nil {
				logger.Errorf("Error finding role %s: %v", userSeeder.TargetRole, err)
				continue
			}

			if err := tx.Model(&user).Session(&gorm.Session{SkipHooks: true}).Association("Role").Replace(&role); err != nil {
				logger.Errorf("Error assigning role %s to user %s: %v", userSeeder.TargetRole, userSeeder.Username, err)
				continue
			}
		}

		return nil
	})

	if transactionError != nil {
		logger.Errorf("Error Seeding Users: %v", transactionError)
		return
	}

	logger.Info("User Seed Finished...")

}

func HandleAdminPassword(defaultPassword string) (string, error) {
	password := os.Getenv("SEED_ADMIN_PASSWORD")

	if password == "" && os.Getenv("APP_ENV") == "production" {
		return "", errors.New("SEED_ADMIN_PASSWORD is required in production")
	} else if password == "" {
		log.Warn("SEED_ADMIN_PASSWORD is not set, using default password for admin user")
		password = defaultPassword
	}

	return password, nil
}
