package productService

import (
	"context"

	"github.com/hengadev/errsx"
)

// (massage, coaching mental etc...)
type Product struct {
	ID          string
	Name        string
	Description string
}

func (p Product) Valid(ctx context.Context) error {
	var errs errsx.Map
	if p.Name != "" {
		errs.Set("name", "cannot have an empty name")
	}
	if p.Description != "" {
		errs.Set("name", "cannot have an empty description")
	}
	return errs.AsError()
}

func (p Product) AssertComparable() {}

func (p Product) GetSQLColumnMapping() map[string]string {
	return map[string]string{
		"ID":          "id",
		"Name":        "name",
		"Description": "description",
	}
}

func (p Product) GetProhibitedFields() []string {
	return []string{
		"ID",
	}
}
