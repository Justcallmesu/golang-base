package navigation_menu_service

import (
	"justcallmesu.com/golang-base/internal/app/models"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (service *Service) GetNavigationMenuBasedOnUserPermission(user *models.User) *[]models.NavigationMenu {
	var filteredMenu []models.NavigationMenu

	for _, navigationMenu := range models.NavigationMenus {
		if user.Role.HasPermission(navigationMenu.Permission) {
			filteredMenu = append(filteredMenu, navigationMenu)
		}
	}

	return &filteredMenu
}
