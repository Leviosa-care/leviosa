package models

type Role int8

const (
	UNKNOWN Role = iota
	ANONYMOUS
	BASIC
	PREMIUM
	GUEST
	FREELANCE
	ADMINISTRATOR
)

var roles = [7]string{
	"unknown",
	"anonymous",
	"basic",
	"premium",
	"guest",
	"freelance",
	"admin",
}

func (r Role) String() string {
	return roles[r]
}

func ConvertToRole(role string) Role {
	switch role {
	case "admin":
		return ADMINISTRATOR
	case "guest":
		return GUEST
	case "anonymous":
		return ANONYMOUS
	case "basic":
		return BASIC
	case "premium":
		return PREMIUM
	case "freelance":
		return FREELANCE
	default:
		return UNKNOWN
	}
}

// Function qui retourne si un role est superieur (ou egal a un autre role).
func (r Role) IsSuperior(role Role) bool {
	switch r {
	case ADMINISTRATOR:
		return role == ADMINISTRATOR || role == GUEST || role == BASIC || role == PREMIUM || role == ANONYMOUS
	case PREMIUM:
		return role == PREMIUM || role == BASIC || role == ANONYMOUS
	case GUEST:
		return role == GUEST
	case BASIC:
		return role == BASIC
	case ANONYMOUS:
		return role == BASIC || role == ANONYMOUS
	default:
		return false
	}
}
