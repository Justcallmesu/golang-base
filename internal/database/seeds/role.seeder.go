package seeds

import (
	"fmt"

	"github.com/charmbracelet/log"
	"gorm.io/gorm"
	"justcallmesu.com/golang-base/internal/app/models"
)

var Roles = []*models.Role{
	{
		Name:        "Admin",
		Description: "Administrator role with all permissions",
		Permissions: Permissions,
	},
}

func RoleSeeder(database *gorm.DB, logger *log.Logger) {
	logger.Info("Seeding Roles...")

	transactionError := database.Transaction(func(tx *gorm.DB) error {
		for _, role := range Roles {
			logger.Warnf("Seeding Role: %s", role.Name)

			if err := tx.Omit("Permissions").Where("name = ?", role.Name).FirstOrCreate(role).Error; err != nil {
				return fmt.Errorf("error seeding role %s: %v", role.Name, err)
			}
		}
		return nil
	})

	if transactionError != nil {
		logger.Errorf("Error Seeding Roles: %v", transactionError)
		return
	}

	logger.Info("Role Seed Finished...")
}

func RolePermissionSeeder(database *gorm.DB, logger *log.Logger) {
	logger.Info("Seeding Role Permissions...")

	transactionError := database.Transaction(func(tx *gorm.DB) error {
		for _, role := range Roles {
			var dbRole models.Role

			logger.Debugf("Associating Permissions for Role: %s", role.Name)

			if err := tx.Where("name = ?", role.Name).First(&dbRole).Error; err != nil && err != gorm.ErrRecordNotFound {
				logger.Errorf("Error Seeding %s Role Permissions: %v", role.Name, err)
				continue
			}

			for _, permission := range role.Permissions {
				var dbPermission models.Permission

				if err := tx.Where("permission_code = ?", permission.PermissionCode).First(&dbPermission).Error; err != nil {
					logger.Errorf("Error Finding Permission %s for Role %s: %v", permission.PermissionCode, role.Name, err)
					continue
				}

				err := tx.Model(&dbRole).Association("Permissions").Append(&dbPermission)
				if err != nil {
					logger.Errorf("Error Associating Permission %s to Role %s: %v", permission.PermissionCode, role.Name, err)
				}
			}

			logger.Infof("Seeded %s Role Permissions", role.Name)
		}

		return nil
	})

	if transactionError != nil {
		logger.Errorf("Error Seeding Role Permissions: %v", transactionError)
		return
	}

	logger.Info("Role Permission Seed Finished...")
}
