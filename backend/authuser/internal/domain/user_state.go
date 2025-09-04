package domain

type UserState string

const (
	Unverified UserState = "unverified"
	Pending    UserState = "pending"
	Active     UserState = "active"
)
