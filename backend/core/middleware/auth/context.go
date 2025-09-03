package auth

import "context"

// sessionContextKey is used to store session data in request context
type sessionContextKey struct{}

// SessionInfoFromContext extracts session info from request context
func SessionInfoFromContext(ctx context.Context) (*SessionInfo, bool) {
	sessionInfo, ok := ctx.Value(sessionContextKey{}).(*SessionInfo)
	if !ok || sessionInfo == nil {
		return nil, false
	}
	return sessionInfo, true
}

// SessionFromContext extracts session from request context
// Deprecated: Use SessionInfoFromContext instead
func SessionFromContext(ctx context.Context) (*Session, bool) {
	// Check if it's the new SessionInfo first
	if sessionInfo, ok := ctx.Value(sessionContextKey{}).(*SessionInfo); ok && sessionInfo != nil {
		// Convert SessionInfo back to Session for backward compatibility
		session := &Session{
			ID:     sessionInfo.ID,
			UserID: sessionInfo.UserID,
			Role:   sessionInfo.Role,
			State:  sessionInfo.State,
		}
		return session, true
	}

	// Fallback to old Session type
	session, ok := ctx.Value(sessionContextKey{}).(*Session)
	if !ok || session == nil {
		return nil, false
	}
	return session, true
}
