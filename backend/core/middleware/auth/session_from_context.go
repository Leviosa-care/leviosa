package auth

import "context"

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