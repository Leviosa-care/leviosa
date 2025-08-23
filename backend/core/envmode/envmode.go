package envmode

import "fmt"

type Mode uint

const (
	Dev Mode = iota
	Prod
	Staging
)

var modeStr = []string{
	"development",
	"production",
	"staging",
}

func (m *Mode) Set(value string) error {
	switch value {
	case "development":
		*m = Dev
	case "production":
		*m = Prod
	case "staging":
		*m = Staging
	default:
		return fmt.Errorf("mode value can only be 'development', 'production' or 'staging', got : %q", value)
	}
	return nil
}

func (m *Mode) String() string {
	return modeStr[*m]
}
