package services

// Service name constants for consistent identification across the system
const (
	AuthUser     = "authuser"
	Catalog      = "catalog"
	Settings     = "settings"
	Notification = "notification"
)

// AllServices returns a slice of all service names for iteration
func AllServices() []string {
	return []string{
		AuthUser,
		Catalog,
		Settings,
		Notification,
	}
}

// IsValidService checks if a service name is valid
func IsValidService(name string) bool {
	for _, service := range AllServices() {
		if service == name {
			return true
		}
	}
	return false
}