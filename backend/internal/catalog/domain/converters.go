package domain

// In a separate helper file (e.g., converters.go) or in your domain package
// toProductRes is a helper function that converts a Product and a Category
// into the public-facing ProductRes DTO.
func ToProductRes(p *Product, c *Category) *ProductRes {
	// You might also need to fetch image keys and populate them here

	return &ProductRes{
		ID:                p.ID,
		Name:              p.Name,
		Description:       p.Description,
		Category:          *c, // Dereference the Category pointer
		Duration:          p.Duration,
		CreatedAt:         p.CreatedAt,
		UpdatedAt:         p.UpdatedAt,
		Status:            p.Status,
		Availability:      p.Availability,
		BufferTime:        p.BufferTime,
		CancellationHours: p.CancellationHours,
		Metadata:          p.Metadata,
	}
}
