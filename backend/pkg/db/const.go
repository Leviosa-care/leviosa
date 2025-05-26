package db

// constants
const DefaultDB = "leviosa.db"

// To write custom update queries based existence or non existence of zero values
type SQLMappable interface {
	GetSQLColumnMapping() map[string]string
	GetProhibitedFields() []string
}
