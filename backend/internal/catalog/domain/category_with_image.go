package domain

type CategoryWithImage struct {
	Category *Category `json:"category"`
	Image    *Image    `json:"image,omitempty"`
}
