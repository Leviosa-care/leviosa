package pg

import "fmt"

func QualifiedTable(schema, table string) string {
	return fmt.Sprintf(`"%s"."%s"`, schema, table)
}
