
-- name: CreateBooking :exec
INSERT INTO bookings (id,user_id, captain_id, pickup_location, drop_location, status)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: GetBooking :one
SELECT * FROM bookings WHERE id = $1 and deleted_at is null;

-- name: ListBookings :many
SELECT * FROM bookings where deleted_at is null;

-- name: UpdateBookingStatus :exec
UPDATE bookings SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2;

-- name: CancelBooking :exec
UPDATE bookings SET status = 'CANCELLED', cancelled_by = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2;

-- name: UpdateBookingPayment :exec
UPDATE bookings SET paid_price = $1, is_paid = $2, payment_method = $3, updated_at = CURRENT_TIMESTAMP WHERE id = $4;

-- name: DeleteBooking :exec
UPDATE bookings SET deleted_at = $1 WHERE id = $2;

-- name: LockBooking :exec
SELECT * FROM bookings WHERE id = $1 FOR UPDATE;


-- name: UpdateBookingSucess :one
UPDATE bookings SET status = 'COMPLETED', is_paid = true, is_successful = true , updated_at = CURRENT_TIMESTAMP WHERE id = $1 RETURNING *;

-- name: UpdateBookingCancel :one
UPDATE bookings SET status = 'CANCELLED',is_cannelled = true, cancelled_by = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2 and is_verified = false RETURNING *;

-- name: UpdateBookingVerify :one
UPDATE bookings SET is_verified = true, updated_at = CURRENT_TIMESTAMP WHERE id = $1 AND status = 'PENDING' AND is_cannelled = false RETURNING *;

-- name: GetUserBookingsByUserId :many
SELECT * FROM bookings WHERE user_id = $1 and deleted_at is null;

-- name: GetCaptainBookingsByCaptainId :many
SELECT * FROM bookings WHERE captain_id = $1 and deleted_at is null;

-- name: GetCurrentBookingByUserId :one
SELECT b.* FROM bookings b left join user u on b.user_id = u.id where b.user_id = $1
and b.id = u.current_booking_id and b.deleted_at is null;

-- name: GetCurrentBookingByCaptainId :one
SELECT b.* FROM bookings b left join captains c on b.captain_id = c.id where b.captain_id = $1
and b.id = c.current_booking_id and b.deleted_at is null;
