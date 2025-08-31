package identity

// Level returns the hierarchical level of the role (higher = more privileges)
func (r Role) Level() int {
	if level, ok := roleLevels[r]; ok {
		return level
	}
	return -1 // Invalid role
}

// IsAtLeast checks if the current role has at least the privilege level of the target role
func (r Role) IsAtLeast(target Role) bool {
	return r.Level() >= target.Level()
}

// HasPermission checks if the current role has permission to perform actions of the target role
func (r Role) HasPermission(target Role) bool {
	return r.IsAtLeast(target)
}