-- name: CreateCaptainVerification :one
INSERT INTO captain_verifications (
    captain_id,
    overall_status,
    current_stage,
    is_blacklisted,
    blacklist_reason
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetCaptainVerification :one
SELECT * FROM captain_verifications
WHERE id = $1 LIMIT 1;

-- name: GetCaptainVerificationByCaptainID :one
SELECT * FROM captain_verifications
WHERE captain_id = $1 LIMIT 1;

-- name: UpdateCaptainVerificationStatus :one
UPDATE captain_verifications
SET overall_status = $2,
    current_stage = $3,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: UpdateCaptainBlacklist :one
UPDATE captain_verifications
SET is_blacklisted = $2,
    blacklist_reason = $3,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: ListPendingVerifications :many
SELECT * FROM captain_verifications
WHERE overall_status = 'PENDING'
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: DeleteCaptainVerification :exec
DELETE FROM captain_verifications
WHERE id = $1;

-- ========== AADHAR DETAILS ==========

-- name: CreateAadharDetails :one
INSERT INTO captain_aadhar_details (
    verification_id,
    aadhar_number,
    aadhar_name,
    front_url,
    back_url,
    status
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetAadharDetails :one
SELECT * FROM captain_aadhar_details
WHERE verification_id = $1 LIMIT 1;

-- name: UpdateAadharStatus :one
UPDATE captain_aadhar_details
SET status = $2,
    admin_comment = $3,
    verified_by = $4,
    verified_at = CURRENT_TIMESTAMP
WHERE verification_id = $1
RETURNING *;

-- name: UpdateAadharDocuments :one
UPDATE captain_aadhar_details
SET front_url = $2,
    back_url = $3
WHERE verification_id = $1
RETURNING *;

-- ========== LICENSE DETAILS ==========

-- name: CreateLicenseDetails :one
INSERT INTO captain_license_details (
    verification_id,
    license_number,
    expiry_date,
    front_url,
    back_url,
    status
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetLicenseDetails :one
SELECT * FROM captain_license_details
WHERE verification_id = $1 LIMIT 1;

-- name: UpdateLicenseStatus :one
UPDATE captain_license_details
SET status = $2,
    admin_comment = $3,
    verified_by = $4,
    verified_at = CURRENT_TIMESTAMP
WHERE verification_id = $1
RETURNING *;

-- name: UpdateLicenseDocuments :one
UPDATE captain_license_details
SET front_url = $2,
    back_url = $3
WHERE verification_id = $1
RETURNING *;

-- name: GetExpiringLicenses :many
SELECT * FROM captain_license_details
WHERE expiry_date <= $1
AND status = 'APPROVED'
ORDER BY expiry_date ASC;

-- ========== VEHICLES ==========

-- name: CreateVehicle :one
INSERT INTO captain_vehicles (
    verification_id,
    vehicle_number,
    vehicle_type,
    rc_book_url,
    insurance_url,
    status
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetVehiclesByVerificationID :many
SELECT * FROM captain_vehicles
WHERE verification_id = $1
ORDER BY created_at DESC;

-- name: GetVehicle :one
SELECT * FROM captain_vehicles
WHERE id = $1 LIMIT 1;

-- name: UpdateVehicleStatus :one
UPDATE captain_vehicles
SET status = $2,
    admin_comment = $3,
    verified_by = $4,
    verified_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: UpdateVehicleDocuments :one
UPDATE captain_vehicles
SET rc_book_url = $2,
    insurance_url = $3
WHERE id = $1
RETURNING *;

-- name: DeleteVehicle :exec
DELETE FROM captain_vehicles
WHERE id = $1;

-- ========== BANK ACCOUNTS ==========

-- name: CreateBankAccount :one
INSERT INTO captain_bank_accounts (
    verification_id,
    account_number,
    ifsc_code,
    account_holder_name,
    status
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetBankAccountsByVerificationID :many
SELECT * FROM captain_bank_accounts
WHERE verification_id = $1
ORDER BY created_at DESC;

-- name: GetBankAccount :one
SELECT * FROM captain_bank_accounts
WHERE id = $1 LIMIT 1;

-- name: UpdateBankAccountStatus :one
UPDATE captain_bank_accounts
SET status = $2,
    admin_comment = $3,
    verified_by = $4,
    verified_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteBankAccount :exec
DELETE FROM captain_bank_accounts
WHERE id = $1;

-- ========== BACKGROUND CHECKS ==========

-- name: CreateBackgroundCheck :one
INSERT INTO captain_background_checks (
    verification_id,
    police_verification_id,
    report_url,
    status
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetBackgroundCheck :one
SELECT * FROM captain_background_checks
WHERE verification_id = $1 LIMIT 1;

-- name: UpdateBackgroundCheckStatus :one
UPDATE captain_background_checks
SET status = $2,
    admin_comment = $3,
    verified_by = $4,
    verified_at = CURRENT_TIMESTAMP
WHERE verification_id = $1
RETURNING *;

-- name: UpdateBackgroundCheckReport :one
UPDATE captain_background_checks
SET report_url = $2
WHERE verification_id = $1
RETURNING *;

-- ========== COMPLETE VERIFICATION DATA ==========

-- name: GetCompleteVerificationData :one
SELECT 
    cv.*,
    json_build_object(
        'id', cad.id,
        'aadhar_number', cad.aadhar_number,
        'aadhar_name', cad.aadhar_name,
        'status', cad.status
    ) as aadhar_details,
    json_build_object(
        'id', cld.id,
        'license_number', cld.license_number,
        'expiry_date', cld.expiry_date,
        'status', cld.status
    ) as license_details
FROM captain_verifications cv
LEFT JOIN captain_aadhar_details cad ON cv.id = cad.verification_id
LEFT JOIN captain_license_details cld ON cv.id = cld.verification_id
WHERE cv.id = $1;

-- name: ListVerificationsByStatus :many
SELECT * FROM captain_verifications
WHERE overall_status = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListVerificationsByStage :many
SELECT * FROM captain_verifications
WHERE current_stage = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetBlacklistedCaptains :many
SELECT * FROM captain_verifications
WHERE is_blacklisted = true
ORDER BY updated_at DESC
LIMIT $1 OFFSET $2;

-- name: CountVerificationsByStatus :one
SELECT COUNT(*) FROM captain_verifications
WHERE overall_status = $1;
