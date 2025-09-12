package auth

import (
	"context"
	"net/http"

	"github.com/Leviosa-care/core/contracts/services"
	"github.com/Leviosa-care/core/ctxutil"
	"github.com/Leviosa-care/core/errs"
	"github.com/Leviosa-care/core/httpx"
	mw "github.com/Leviosa-care/core/middleware"
)

// ServiceInfo contains service authentication details for request context
type ServiceInfo struct {
	Name string `json:"name"`
}

// ServiceContextKey is the key used to store service info in request context
const ServiceContextKey = "service_info"

// RequireServiceAuth validates service authentication headers and makes service info available in context
func (m *SessionAuthMiddleware) RequireServiceAuth(next mw.Handler) mw.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger, err := ctxutil.GetLoggerFromContext(ctx)
		if err != nil {
			httpx.RespondWithError(w, err, http.StatusInternalServerError)
			return
		}

		// Extract service authentication headers
		serviceName := r.Header.Get(services.ServiceNameHeader)
		serviceKey := r.Header.Get(services.ServiceKeyHeader)

		if serviceName == "" {
			logger.WarnContext(ctx, "Service auth middleware: Missing service name header",
				"operation", "require_service_auth",
				"method", r.Method,
				"path", r.URL.Path,
				"header", services.ServiceNameHeader)
			httpx.RespondWithError(w, errs.ErrUnauthorized, http.StatusUnauthorized)
			return
		}

		if serviceKey == "" {
			logger.WarnContext(ctx, "Service auth middleware: Missing service key header",
				"operation", "require_service_auth",
				"method", r.Method,
				"path", r.URL.Path,
				"service_name", serviceName,
				"header", services.ServiceKeyHeader)
			httpx.RespondWithError(w, errs.ErrUnauthorized, http.StatusUnauthorized)
			return
		}

		// Validate service name
		if !services.IsValidService(serviceName) {
			logger.WarnContext(ctx, "Service auth middleware: Invalid service name",
				"operation", "require_service_auth",
				"method", r.Method,
				"path", r.URL.Path,
				"service_name", serviceName,
				"valid_services", services.AllServices())
			httpx.RespondWithError(w, errs.ErrInvalidValue, http.StatusBadRequest)
			return
		}

		// Validate service key against Vault
		// TODO: Implement Vault service key validation once Vault client is added to middleware
		// For now, we'll use a placeholder validation
		if !m.validateServiceKey(ctx, serviceName, serviceKey) {
			logger.WarnContext(ctx, "Service auth middleware: Invalid service key",
				"operation", "require_service_auth",
				"method", r.Method,
				"path", r.URL.Path,
				"service_name", serviceName)
			httpx.RespondWithError(w, errs.ErrUnauthorized, http.StatusUnauthorized)
			return
		}

		// Create service info for context
		serviceInfo := &ServiceInfo{
			Name: serviceName,
		}

		// Add service info to context
		ctx = context.WithValue(ctx, ServiceContextKey, serviceInfo)
		r = r.WithContext(ctx)

		logger.InfoContext(ctx, "Service auth middleware: Service authentication successful",
			"operation", "require_service_auth",
			"method", r.Method,
			"path", r.URL.Path,
			"service_name", serviceName)

		// Continue to next handler
		next(w, r)
	}
}

// validateServiceKey validates the service key against stored service credentials in Vault
func (m *SessionAuthMiddleware) validateServiceKey(ctx context.Context, serviceName, serviceKey string) bool {
	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		// If we can't get logger, continue with validation but log to stderr
		return false
	}

	// Get the Vault path for this service's API key
	vaultPath := services.ServiceAPIKeyPath(serviceName)

	logger.InfoContext(ctx, "Service auth middleware: Validating service key with Vault",
		"operation", "validate_service_key",
		"service_name", serviceName,
		"vault_path", vaultPath)

	// Read the stored API key from Vault
	secret, err := m.vaultClient.Logical().Read(vaultPath)
	if err != nil {
		logger.ErrorContext(ctx, "Service auth middleware: Failed to read service key from Vault",
			"operation", "validate_service_key",
			"service_name", serviceName,
			"vault_path", vaultPath,
			"error", err)
		return false
	}

	if secret == nil || secret.Data == nil {
		logger.WarnContext(ctx, "Service auth middleware: Service key not found in Vault",
			"operation", "validate_service_key",
			"service_name", serviceName,
			"vault_path", vaultPath)
		return false
	}

	// Extract the stored key hash from Vault response
	// Vault KV v2 nests data under "data" key
	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		logger.ErrorContext(ctx, "Service auth middleware: Invalid Vault response format",
			"operation", "validate_service_key",
			"service_name", serviceName,
			"vault_path", vaultPath)
		return false
	}

	storedKeyHash, ok := data["key_hash"].(string)
	if !ok || storedKeyHash == "" {
		logger.WarnContext(ctx, "Service auth middleware: Missing or invalid key_hash in Vault",
			"operation", "validate_service_key",
			"service_name", serviceName,
			"vault_path", vaultPath)
		return false
	}

	// Hash the provided key and compare with stored hash
	providedKeyHash := m.crypto.HashBasic(ctx, []byte(serviceKey))

	isValid := storedKeyHash == providedKeyHash

	if isValid {
		logger.InfoContext(ctx, "Service auth middleware: Service key validation successful",
			"operation", "validate_service_key",
			"service_name", serviceName)
	} else {
		logger.WarnContext(ctx, "Service auth middleware: Service key validation failed",
			"operation", "validate_service_key",
			"service_name", serviceName)
	}

	return isValid
}

// GetServiceInfoFromContext extracts service info from request context
func GetServiceInfoFromContext(ctx context.Context) (*ServiceInfo, error) {
	serviceInfo, ok := ctx.Value(ServiceContextKey).(*ServiceInfo)
	if !ok {
		return nil, errs.NewInvalidValueErr("service info not found in context")
	}
	return serviceInfo, nil
}

