package models

type Role int8

const (
	VISITOR Role = iota
	STANDARD
	PREMIUM
	GUEST
	PARTNER
	ADMINISTRATOR
)

var roles = [6]string{
	"visitor",
	"standard",
	"premium",
	"guest",
	"partner",
	"administrator",
}

func (r Role) String() string {
	return roles[r]
}

func ConvertToRole(role string) Role {
	switch role {
	case "administrator":
		return ADMINISTRATOR
	case "guest":
		return GUEST
	case "standard":
		return STANDARD
	case "premium":
		return PREMIUM
	case "partner":
		return PARTNER
	default:
		return VISITOR
	}
}
