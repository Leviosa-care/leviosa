package session

import "context"

// SessionInfoFromContext extracts session info from request context
func SessionInfoFromContext(ctx context.Context) (*SessionInfo, bool) {
	sessionInfo, ok := ctx.Value(sessionContextKey{}).(*SessionInfo)
	if !ok || sessionInfo == nil {
		return nil, false
	}
	return sessionInfo, true
}
