package auth

import "context"

// sessionContextKey is used to store session data in request context
type sessionContextKey struct{}

// SessionFromContext extracts session from request context
func SessionFromContext(ctx context.Context) (*Session, bool) {
	session, ok := ctx.Value(sessionContextKey{}).(*Session)
	if !ok || session == nil {
		return nil, false
	}
	return session, true
}
