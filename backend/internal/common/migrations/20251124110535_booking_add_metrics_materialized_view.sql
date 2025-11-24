-- +goose Up
-- +goose StatementBegin

-- Create materialized view for room utilization metrics
CREATE MATERIALIZED VIEW booking.room_daily_metrics AS
WITH room_operating_hours AS (
    SELECT
        r.id as room_id,
        EXTRACT(EPOCH FROM (r.operating_end_time - r.operating_start_time)) / 60 as minutes_open
    FROM booking.rooms r
    WHERE r.is_active = true
),
daily_availabilities AS (
    SELECT
        a.room_id,
        DATE(a.start_time) as date,
        SUM(EXTRACT(EPOCH FROM (a.end_time - a.start_time)) / 60) as total_available_minutes,
        SUM(
            CASE WHEN a.status = 'booked'
            THEN EXTRACT(EPOCH FROM (a.end_time - a.start_time)) / 60
            ELSE 0
            END
        ) as total_booked_minutes
    FROM booking.availabilities a
    WHERE a.status IN ('available', 'booked')
    GROUP BY a.room_id, DATE(a.start_time)
),
daily_gaps AS (
    SELECT
        room_id,
        date,
        COUNT(*) as gap_count,
        SUM(gap_minutes) as total_gap_minutes,
        COALESCE(AVG(gap_minutes), 0) as average_gap_minutes,
        SUM(CASE WHEN gap_minutes < 30 THEN 1 ELSE 0 END) as fragmentation_count,
        SUM(CASE WHEN gap_minutes < 30 THEN gap_minutes ELSE 0 END) as idle_minutes
    FROM (
        SELECT
            room_id,
            DATE(start_time) as date,
            EXTRACT(EPOCH FROM (
                LEAD(start_time) OVER (PARTITION BY room_id, DATE(start_time) ORDER BY start_time) -
                end_time
            )) / 60 as gap_minutes
        FROM booking.availabilities
        WHERE status IN ('available', 'booked')
    ) gaps
    WHERE gap_minutes > 0
    GROUP BY room_id, date
)
SELECT
    da.room_id,
    da.date,
    roh.minutes_open as total_minutes_open,
    COALESCE(da.total_available_minutes, 0) as total_minutes_available,
    COALESCE(da.total_booked_minutes, 0) as total_minutes_booked,
    CASE
        WHEN roh.minutes_open > 0
        THEN (COALESCE(da.total_booked_minutes, 0) / roh.minutes_open) * 100
        ELSE 0
    END as utilization_percent,
    COALESCE(dg.fragmentation_count, 0) as fragmentation_count,
    COALESCE(dg.idle_minutes, 0)::integer as idle_minutes,
    COALESCE(dg.average_gap_minutes, 0)::integer as average_gap_minutes,
    NOW() as created_at,
    NOW() as updated_at
FROM daily_availabilities da
JOIN room_operating_hours roh ON da.room_id = roh.room_id
LEFT JOIN daily_gaps dg ON da.room_id = dg.room_id AND da.date = dg.date;

-- Create indexes for fast queries
CREATE UNIQUE INDEX idx_room_daily_metrics_room_date
    ON booking.room_daily_metrics(room_id, date);

CREATE INDEX idx_room_daily_metrics_date
    ON booking.room_daily_metrics(date);

CREATE INDEX idx_room_daily_metrics_room_id
    ON booking.room_daily_metrics(room_id);

-- Create function to refresh materialized view
CREATE OR REPLACE FUNCTION booking.refresh_room_metrics()
RETURNS void AS $$
BEGIN
    REFRESH MATERIALIZED VIEW CONCURRENTLY booking.room_daily_metrics;
END;
$$ LANGUAGE plpgsql;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Drop function
DROP FUNCTION IF EXISTS booking.refresh_room_metrics();

-- Drop materialized view (indexes will be dropped automatically)
DROP MATERIALIZED VIEW IF EXISTS booking.room_daily_metrics;

-- +goose StatementEnd
