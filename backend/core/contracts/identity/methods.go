package identity

// String returns the string representation of the role
func (r Role) String() string {
	if str, ok := roleStrings[r]; ok {
		return str
	}
	return "unknown"
}

// IsValid checks if the role is a valid defined role
func (r Role) IsValid() bool {
	_, ok := roleStrings[r]
	return ok
}

// IsAdmin checks if the role is Administrator
func (r Role) IsAdmin() bool {
	return r == Administrator
}

