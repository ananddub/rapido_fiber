-- name: CreateUser :exec
INSERT INTO users (name, phone) VALUES ($1, $2);

-- name: GetUser :one
SELECT * FROM users WHERE id = $1 and deleted_at is null;

-- name: GetUserById :one
SELECT * FROM users WHERE id = $1 and deleted_at is null;

-- name: GetUserByPhone :one
SELECT * FROM users WHERE phone = $1 and deleted_at is null;

-- name: ListUsers :many
SELECT * FROM users where deleted_at is null;

-- name: UpdateUser :exec
UPDATE users SET name = $1, phone = $2 WHERE id = $3;

-- name: DeleteUser :exec
UPDATE users SET deleted_at = $1 WHERE id = $2;

-- name: RestoreUser :exec
UPDATE users SET deleted_at = $1 WHERE id = $2;

-- name: CreateCaptain :exec
INSERT INTO captains (name, phone) VALUES ($1, $2);

-- name: GetCaptain :one
SELECT * FROM captains WHERE id = $1 and deleted_at is null;

-- name: GetCaptainById :one
SELECT * FROM captains WHERE id = $1 and deleted_at is null;

-- name: ListCaptains :many
SELECT * FROM captains where deleted_at is null;

-- name: UpdateCaptain :exec
UPDATE captains SET name = $1, phone = $2 WHERE id = $3;

-- name: DeleteCaptain :exec
UPDATE captains SET deleted_at = $1 WHERE id = $2;

-- name: RestoreCaptain :exec
UPDATE captains SET deleted_at = $1 WHERE id = $2;

-- name: GetCaptainByPhone :one
SELECT * FROM captains
WHERE phone = $1 AND deleted_at IS NULL;

-- name: UpdateUserStatus :exec
UPDATE users
SET
    status = $2,
    current_booking_id = $3,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: UpdateCaptainStatus :exec
UPDATE captains
SET
    status = $2,
    current_booking_id = $3,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: LockUser :exec
SELECT * FROM users WHERE id = $1 FOR UPDATE;

-- name: LockCaptain :exec
SELECT * FROM captains WHERE id = $1 FOR UPDATE;
