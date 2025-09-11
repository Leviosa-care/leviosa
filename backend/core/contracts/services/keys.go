package services

import "fmt"

// Vault path generators for consistent service key management

// ServicePepperPath generates the Vault path for a service's pepper secret
func ServicePepperPath(serviceName string) string {
	return fmt.Sprintf("secret/data/%s/pepper", serviceName)
}

// ServiceKEKPath generates the Vault path for a service's Key Encryption Key
func ServiceKEKPath(serviceName string) string {
	return fmt.Sprintf("transit/keys/%s-kek", serviceName)
}

// ServiceAPIKeyPath generates the Vault path for a service's API key
func ServiceAPIKeyPath(serviceName string) string {
	return fmt.Sprintf("secret/data/services/%s/api-key", serviceName)
}

// ServiceVaultPaths returns all Vault paths for a given service
func ServiceVaultPaths(serviceName string) map[string]string {
	return map[string]string{
		"pepper":  ServicePepperPath(serviceName),
		"kek":     ServiceKEKPath(serviceName),
		"api_key": ServiceAPIKeyPath(serviceName),
	}
}