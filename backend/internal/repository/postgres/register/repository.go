package registerRepository

import (
	"context"
	"database/sql"

	"github.com/hengadev/leviosa/internal/domain/register"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type repository struct {
	DB     *sql.DB
	schema string
}

func (r *repository) GetDB() *sql.DB {
	return r.DB
}

func New(ctx context.Context, db *sql.DB) (registerService.ReadWriter, error) {
	return &repository{DB: db, schema: "register"}, nil
}

// NOTE: old API

//	func (s *Store) CreateRegistration(registration *types.Registration) error {
//		_, err := s.DB.Exec("INSERT INTO ? (userid, beginAt) VALUES (?, ?);", registration.EventId, registration.UserId, registration.BeginAt.Format(time.RFC3339))
//		if err != nil {
//			return err
//		}
//		return nil
//	}
//
//	func (s *Store) GetBeginAtByEventId(eventId string) time.Time {
//		var dateTemp string
//		err := s.DB.QueryRow("SELECT beginAt FROM events WHERE eventid=?;", eventId).Scan(&dateTemp)
//		if err != nil {
//			log.Fatalf("Failed to get the beginAt field from the row corresponding to the event %s - %s", eventId, err)
//		}
//		res, err := time.Parse(time.RFC3339, dateTemp)
//		if err != nil {
//			log.Fatal("Failed to parse the time from the database - ", err)
//		}
//		return res
//	}
