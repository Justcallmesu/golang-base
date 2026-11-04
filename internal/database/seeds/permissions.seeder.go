package seeds

import (
	"github.com/charmbracelet/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"justcallmesu.com/golang-base/internal/app/models"
)

var Permissions = []models.Permission{
	{
		PermissionCode: models.PERMISSION_CREATE_USER,
		Group:          "User Management",
		Description:    "Permission to create a new user",
	},
	{
		PermissionCode: models.PERMISSION_UPDATE_USER,
		Group:          "User Management",
		Description:    "Permission to edit an existing user",
	},
	{
		PermissionCode: models.PERMISSION_DELETE_USER,
		Group:          "User Management",
		Description:    "Permission to delete a user",
	},
	{
		PermissionCode: models.PERMISSION_VIEW_USER,
		Group:          "User Management",
		Description:    "Permission to view user details",
	},
	{
		PermissionCode: models.PERMISSION_CREATE_ROLE,
		Group:          "Role Management",
		Description:    "Permission to create a new role",
	},
	{
		PermissionCode: models.PERMISSION_UPDATE_ROLE,
		Group:          "Role Management",
		Description:    "Permission to edit an existing role",
	},
	{
		PermissionCode: models.PERMISSION_DELETE_ROLE,
		Group:          "Role Management",
		Description:    "Permission to delete a role",
	},
	{
		PermissionCode: models.PERMISSION_VIEW_ROLE,
		Group:          "Role Management",
		Description:    "Permission to view role details",
	},
}

func PermissionSeeder(database *gorm.DB, logger *log.Logger) {

	logger.Info("Running Permission Seeders...")

	transactionError := database.Transaction(func(tx *gorm.DB) error {

		return tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "permission_code"}},
			DoUpdates: clause.AssignmentColumns([]string{"group", "description"}),
		}).Create(&Permissions).Error

	})

	if transactionError != nil {
		logger.Errorf("Error Seeding Permissions: %v", transactionError)
		return
	}

	logger.Info("Permission Seed Finished...")

}
