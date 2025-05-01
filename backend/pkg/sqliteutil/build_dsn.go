package sqliteutil

import (
	"fmt"

	"github.com/hengadev/leviosa/pkg/flags"
)

func BuildDSN(env mode.EnvMode, name string) string {
	// here is how it is supposed to be brother
	dsn := fmt.Sprintf("file:%s.db?_journal_mode=WAL&_busy_timeout=5000&_foreign_keys=ON", name)
	password := "somepassword"
	if env == mode.ModeStaging || env == mode.ModeProd {
		// TODO: here just add what is need for the encryption using sqlcypher brother
		return fmt.Sprintf("%s&_pragma_key=%s&_pragma_cipher_page_size=4096", dsn, password)
	}
	return dsn
}
