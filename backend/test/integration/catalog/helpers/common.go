package helpers

// Helper to create a pointer to a string
func StrPtr(s string) *string { return &s }

// Helper to create a pointer to a string
func IntPtr(i int) *int { return &i }

// Helper to create a pointer to a string
func BoolPtr(b bool) *bool { return &b }

// Helper to create a pointer to a PublishedStatus
func StatusStrPtr(s string) *string { return &s }
