package postgres_test

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/adapters/postgres"
	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/common/migrations"
	tu "github.com/Leviosa-care/leviosa/backend/internal/common/testutils"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	pgContainer *tu.PostgresContainer
	testPool    *pgxpool.Pool
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	// Setup PostgreSQL testcontainer
	var err error
	pgContainer, err = tu.SetupPostgres(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to setup postgres: %v", err)
	}
	defer tu.TeardownPostgres(ctx, nil, pgContainer)

	// Create pool
	testPool, err = pgxpool.New(ctx, pgContainer.ConnectionString)
	if err != nil {
		log.Fatalf("Failed to create pool: %v", err)
	}
	defer testPool.Close()

	// Run migrations
	goose.SetBaseFS(migrations.FS)
	if err = goose.SetDialect("pgx"); err != nil {
		log.Fatalf("Setting dialect for migrations: %v", err)
	}

	gooseDB, err := sql.Open("pgx", testPool.Config().ConnString())
	if err != nil {
		log.Fatalf("Failed to open temp *sql.DB for goose migrations: %v", err)
	}
	defer gooseDB.Close()

	if err = goose.UpContext(ctx, gooseDB, "."); err != nil {
		log.Fatalf("Migrations failed: %v", err)
	}

	// Run tests
	os.Exit(m.Run())
}

func TestGetRoomHoursForDate(t *testing.T) {
	ctx := context.Background()
	repo := postgres.NewRoomScheduleRepository(testPool)

	// Setup test room
	roomID := uuid.New()
	insertTestRoom(t, ctx, roomID)
	defer cleanupRooms(t, ctx)

	t.Run("returns recurring schedule for matching day of week", func(t *testing.T) {
		// Clean previous schedules
		cleanupSchedules(t, ctx)

		// Saturday = 6
		saturday := 6
		schedule := &domain.RoomAvailabilitySchedule{
			RoomID:    roomID,
			DayOfWeek: &saturday,
			OpenTime:  parseTime(t, "09:00:00"),
			CloseTime: parseTime(t, "17:00:00"),
			Priority:  1,
			IsActive:  true,
		}

		err := repo.Create(ctx, schedule)
		require.NoError(t, err)

		// Query for a Saturday (December 6, 2025 is a Saturday)
		testDate := time.Date(2025, 12, 6, 0, 0, 0, 0, time.UTC)

		result, err := repo.GetRoomHoursForDate(ctx, roomID, testDate)
		require.NoError(t, err)
		assert.Equal(t, roomID, result.RoomID)
		assert.Equal(t, "09:00", result.OpenTime.Format("15:04"))
		assert.Equal(t, "17:00", result.CloseTime.Format("15:04"))
	})

	t.Run("returns specific date schedule over recurring pattern", func(t *testing.T) {
		// Clean previous schedules
		cleanupSchedules(t, ctx)

		// Setup recurring Saturday schedule (priority 1)
		saturday := 6
		recurringSchedule := &domain.RoomAvailabilitySchedule{
			RoomID:    roomID,
			DayOfWeek: &saturday,
			OpenTime:  parseTime(t, "09:00:00"),
			CloseTime: parseTime(t, "17:00:00"),
			Priority:  1,
			IsActive:  true,
		}
		err := repo.Create(ctx, recurringSchedule)
		require.NoError(t, err)

		// Setup specific date override (higher priority)
		specificDate := time.Date(2025, 12, 13, 0, 0, 0, 0, time.UTC) // Also a Saturday
		specificSchedule := &domain.RoomAvailabilitySchedule{
			RoomID:       roomID,
			SpecificDate: &specificDate,
			OpenTime:     parseTime(t, "10:00:00"),
			CloseTime:    parseTime(t, "21:00:00"),
			Priority:     10,
			IsActive:     true,
		}
		err = repo.Create(ctx, specificSchedule)
		require.NoError(t, err)

		// Query for the specific date
		result, err := repo.GetRoomHoursForDate(ctx, roomID, specificDate)
		require.NoError(t, err)

		// Should return specific date hours, not recurring
		assert.Equal(t, "10:00", result.OpenTime.Format("15:04"))
		assert.Equal(t, "21:00", result.CloseTime.Format("15:04"))
		assert.NotNil(t, result.SpecificDate)
	})

	t.Run("returns ErrRepositoryNotFound when no schedule exists", func(t *testing.T) {
		nonexistentRoomID := uuid.New()
		testDate := time.Date(2025, 12, 15, 0, 0, 0, 0, time.UTC)

		_, err := repo.GetRoomHoursForDate(ctx, nonexistentRoomID, testDate)
		assert.ErrorIs(t, err, errs.ErrRepositoryNotFound)
	})
}

func TestCreate(t *testing.T) {
	ctx := context.Background()
	repo := postgres.NewRoomScheduleRepository(testPool)

	// Setup test room
	roomID := uuid.New()
	insertTestRoom(t, ctx, roomID)
	defer cleanupRooms(t, ctx)

	t.Run("successfully creates recurring schedule", func(t *testing.T) {
		sunday := 0
		schedule := &domain.RoomAvailabilitySchedule{
			RoomID:    roomID,
			DayOfWeek: &sunday,
			OpenTime:  parseTime(t, "08:00:00"),
			CloseTime: parseTime(t, "20:00:00"),
			Priority:  1,
			IsActive:  true,
		}

		err := repo.Create(ctx, schedule)
		require.NoError(t, err)

		assert.NotEqual(t, uuid.Nil, schedule.ID)
		assert.False(t, schedule.CreatedAt.IsZero())
		assert.False(t, schedule.UpdatedAt.IsZero())
	})

	t.Run("successfully creates specific date schedule", func(t *testing.T) {
		specificDate := time.Date(2025, 12, 25, 0, 0, 0, 0, time.UTC)
		schedule := &domain.RoomAvailabilitySchedule{
			RoomID:       roomID,
			SpecificDate: &specificDate,
			OpenTime:     parseTime(t, "12:00:00"),
			CloseTime:    parseTime(t, "16:00:00"),
			Priority:     5,
			IsActive:     true,
		}

		err := repo.Create(ctx, schedule)
		require.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, schedule.ID)
	})
}

func TestUpdate(t *testing.T) {
	ctx := context.Background()
	repo := postgres.NewRoomScheduleRepository(testPool)

	// Setup test room
	roomID := uuid.New()
	insertTestRoom(t, ctx, roomID)
	defer cleanupRooms(t, ctx)

	t.Run("successfully updates schedule", func(t *testing.T) {
		cleanupSchedules(t, ctx)

		// Create initial schedule
		sunday := 0
		schedule := &domain.RoomAvailabilitySchedule{
			RoomID:    roomID,
			DayOfWeek: &sunday,
			OpenTime:  parseTime(t, "08:00:00"),
			CloseTime: parseTime(t, "17:00:00"),
			Priority:  1,
			IsActive:  true,
		}

		err := repo.Create(ctx, schedule)
		require.NoError(t, err)

		// Update the schedule
		schedule.OpenTime = parseTime(t, "09:00:00")
		schedule.CloseTime = parseTime(t, "18:00:00")
		schedule.Priority = 2

		err = repo.Update(ctx, schedule)
		require.NoError(t, err)

		// Verify update
		updated, err := repo.GetByID(ctx, schedule.ID)
		require.NoError(t, err)
		assert.Equal(t, "09:00", updated.OpenTime.Format("15:04"))
		assert.Equal(t, "18:00", updated.CloseTime.Format("15:04"))
		assert.Equal(t, 2, updated.Priority)
	})
}

func TestDelete(t *testing.T) {
	ctx := context.Background()
	repo := postgres.NewRoomScheduleRepository(testPool)

	// Setup test room
	roomID := uuid.New()
	insertTestRoom(t, ctx, roomID)
	defer cleanupRooms(t, ctx)

	t.Run("successfully soft-deletes schedule", func(t *testing.T) {
		cleanupSchedules(t, ctx)

		// Create schedule
		sunday := 0
		schedule := &domain.RoomAvailabilitySchedule{
			RoomID:    roomID,
			DayOfWeek: &sunday,
			OpenTime:  parseTime(t, "08:00:00"),
			CloseTime: parseTime(t, "17:00:00"),
			Priority:  1,
			IsActive:  true,
		}

		err := repo.Create(ctx, schedule)
		require.NoError(t, err)

		// Delete the schedule
		err = repo.Delete(ctx, schedule.ID)
		require.NoError(t, err)

		// Verify it's soft-deleted (GetByID should fail because it filters is_active=true)
		_, err = repo.GetByID(ctx, schedule.ID)
		assert.ErrorIs(t, err, errs.ErrRepositoryNotFound)
	})
}

func TestListByRoom(t *testing.T) {
	ctx := context.Background()
	repo := postgres.NewRoomScheduleRepository(testPool)

	// Setup test room
	roomID := uuid.New()
	insertTestRoom(t, ctx, roomID)
	defer cleanupRooms(t, ctx)

	t.Run("lists all active schedules for a room", func(t *testing.T) {
		cleanupSchedules(t, ctx)

		// Create multiple schedules
		sunday := 0
		schedule1 := &domain.RoomAvailabilitySchedule{
			RoomID:    roomID,
			DayOfWeek: &sunday,
			OpenTime:  parseTime(t, "08:00:00"),
			CloseTime: parseTime(t, "17:00:00"),
			Priority:  1,
			IsActive:  true,
		}
		err := repo.Create(ctx, schedule1)
		require.NoError(t, err)

		monday := 1
		schedule2 := &domain.RoomAvailabilitySchedule{
			RoomID:    roomID,
			DayOfWeek: &monday,
			OpenTime:  parseTime(t, "09:00:00"),
			CloseTime: parseTime(t, "18:00:00"),
			Priority:  1,
			IsActive:  true,
		}
		err = repo.Create(ctx, schedule2)
		require.NoError(t, err)

		// List schedules
		schedules, err := repo.ListByRoom(ctx, roomID)
		require.NoError(t, err)
		assert.Len(t, schedules, 2)
	})
}

// Test helpers
func insertTestRoom(t *testing.T, ctx context.Context, roomID uuid.UUID) {
	// Insert minimal room data for FK constraints
	query := `
		INSERT INTO booking.buildings (id, name_encrypted, name_hash, address_encrypted, address_hash,
			city_encrypted, city_hash, postal_code_encrypted, country_encrypted, country_hash,
			dek_encrypted, key_version)
		VALUES ($1, 'encrypted', 'hash', 'encrypted', 'hash', 'encrypted', 'hash', 'encrypted', 'encrypted', 'hash', 'encrypted', 1)
	`
	buildingID := uuid.New()
	_, err := testPool.Exec(ctx, query, buildingID)
	require.NoError(t, err)

	query = `
		INSERT INTO booking.rooms (id, building_id, name_encrypted, name_hash, capacity, dek_encrypted, key_version, is_active)
		VALUES ($1, $2, 'encrypted', 'hash', 1, 'encrypted', 1, true)
	`
	_, err = testPool.Exec(ctx, query, roomID, buildingID)
	require.NoError(t, err)
}

func cleanupRooms(t *testing.T, ctx context.Context) {
	_, _ = testPool.Exec(ctx, "DELETE FROM booking.room_availability_schedules")
	_, _ = testPool.Exec(ctx, "DELETE FROM booking.rooms")
	_, _ = testPool.Exec(ctx, "DELETE FROM booking.buildings")
}

func cleanupSchedules(t *testing.T, ctx context.Context) {
	_, _ = testPool.Exec(ctx, "DELETE FROM booking.room_availability_schedules")
}

func parseTime(t *testing.T, timeStr string) time.Time {
	parsed, err := time.Parse("15:04:05", timeStr)
	require.NoError(t, err)
	return parsed
}
