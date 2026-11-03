package models

type NavigationMenu struct {
	Label string `json:"label"`

	// For Anyone Wondering Where To Get Icon, use Lucide Icons https://lucide.dev/icons/.
	// Attach icon name with i-lucide- prefix. Example: lucide-home, lucide-user, lucide-settings etc.
	Icon       string             `json:"icon"`
	To         string             `json:"to"`
	Children   *[]NavigationMenu  `json:"children,omitempty,omitzero"`
	Permission PermissionCodeEnum `json:"-"`
}

var NavigationMenus = []NavigationMenu{
	{
		Label:      "User Management",
		Icon:       "i-lucide-user-pen",
		To:         "/users",
		Permission: PERMISSION_VIEW_USER,
	},
	{
		Label:      "Role Management",
		Icon:       "i-lucide-shield-check",
		To:         "/roles",
		Permission: PERMISSION_VIEW_ROLE,
	},
}
