package session

// sessionContextKey is used to store session data in request context
type sessionContextKey struct{}

// GetSessionContextKey returns the context key for session data
func GetSessionContextKey() sessionContextKey {
	return sessionContextKey{}
}

