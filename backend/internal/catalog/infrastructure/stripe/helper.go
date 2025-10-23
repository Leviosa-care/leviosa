package payment

import (
	"fmt"
	"strings"
)

func generateIdempotencyKey(prefix, name string, priceInCents int, currency string) string {
	return fmt.Sprintf("%s_%s_%d_%s",
		prefix,
		strings.ReplaceAll(name, " ", "_"),
		priceInCents,
		currency)
}
