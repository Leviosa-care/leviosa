package identity

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
	if r < 0 || int(r) >= len(roles) {
		return "unknown"
	}
	return roles[r]
}

func ConvertToRole(role string) (Role, bool) {
	switch role {
	case "administrator":
		return ADMINISTRATOR, true
	case "guest":
		return GUEST, true
	case "standard":
		return STANDARD, true
	case "premium":
		return PREMIUM, true
	case "partner":
		return PARTNER, true
	case "visitor":
		return VISITOR, true
	default:
		return VISITOR, false
	}
}

func (r Role) IsValid() bool {
	return r >= VISITOR && r <= ADMINISTRATOR
}

func (r Role) IsAdmin() bool {
	return r == ADMINISTRATOR
}
