-- +goose Up
-- +goose StatementBegin

CREATE TYPE booking_status AS ENUM (
    'PENDING',
    'CONFIRMED',
    'COMPLETED',
    'CANCELLED'
);

CREATE TABLE IF NOT EXISTS bookings (
    id UUID PRIMARY KEY,
    user_id INTEGER NOT NULL,
    captain_id INTEGER NOT NULL,
    pickup_location VARCHAR(255) NOT NULL,
    drop_location VARCHAR(255) NOT NULL,
    actual_price DECIMAL(10, 2) NOT NULL,
    paid_price DECIMAL(10, 2) NOT NULL,
    is_paid BOOLEAN DEFAULT FALSE,
    is_successful BOOLEAN DEFAULT FALSE,
    is_verified BOOLEAN DEFAULT FALSE,
    payment_method VARCHAR(50) DEFAULT NULL,
    status booking_status DEFAULT 'PENDING',
    is_cancelled BOOLEAN DEFAULT FALSE,
    cancelled_by VARCHAR(50) DEFAULT NULL,
    is_deleted BOOLEAN DEFAULT FALSE,
    deleted_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS bookings;
DROP TYPE IF EXISTS booking_status;
-- +goose StatementEnd
