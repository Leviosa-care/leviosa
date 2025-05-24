package registerRepository

import "fmt"

func getTablename(day, year int, month string) string {
	return fmt.Sprintf("registration_%d_%s_%d", day, month, year)
}
