package identity

// Role represents a user role in the system
type Role int8

// Role enumeration values
const (
	Visitor Role = iota
	Standard
	Premium
	Guest
	Partner
	Administrator
)

// Role string constants for external use
const (
	VisitorStr       = "visitor"
	StandardStr      = "standard"
	PremiumStr       = "premium"
	GuestStr         = "guest"
	PartnerStr       = "partner"
	AdministratorStr = "administrator"
)

// roleStrings maps roles to their string representations
var roleStrings = map[Role]string{
	Visitor:       VisitorStr,
	Standard:      StandardStr,
	Premium:       PremiumStr,
	Guest:         GuestStr,
	Partner:       PartnerStr,
	Administrator: AdministratorStr,
}

// stringRoles maps string representations to roles (case-insensitive)
var stringRoles = map[string]Role{
	VisitorStr:       Visitor,
	StandardStr:      Standard,
	PremiumStr:       Premium,
	GuestStr:         Guest,
	PartnerStr:       Partner,
	AdministratorStr: Administrator,
}

// roleLevels defines the hierarchy of roles (higher number = more privileges)
var roleLevels = map[Role]int{
	Visitor:       0,
	Guest:         1,
	Standard:      2,
	Premium:       3,
	Partner:       4,
	Administrator: 5,
}

