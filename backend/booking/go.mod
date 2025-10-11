module github.com/Leviosa-care/booking

go 1.24.2

require (
	github.com/Leviosa-care/core v0.0.0
	github.com/google/uuid v1.6.0
	github.com/hengadev/encx v0.5.2
	github.com/hengadev/errsx v1.2.0
	github.com/jackc/pgx/v5 v5.7.5
	github.com/pressly/goose/v3 v3.25.0
	github.com/rabbitmq/amqp091-go v1.10.0
	github.com/redis/go-redis/v9 v9.12.1
	github.com/stretchr/testify v1.11.0
)

replace github.com/Leviosa-care/core => ../core