-- +goose Up
-- +goose StatementBegin
CREATE TYPE verification_status_enum AS ENUM (
    'PENDING',
    'PARTIAL',
    'APPROVED',
    'REJECTED'
);

CREATE TYPE document_status_enum AS ENUM (
    'PENDING',
    'APPROVED',
    'REJECTED'
);

CREATE TYPE verification_stage_enum AS ENUM (
    'BASIC',
    'DOCUMENTS',
    'BACKGROUND',
    'FINAL'
);

CREATE TABLE captain_verifications (
    id SERIAL PRIMARY KEY,
    captain_id INT NOT NULL REFERENCES captains(id),

    overall_status verification_status_enum DEFAULT 'PENDING',
    current_stage verification_stage_enum DEFAULT 'BASIC',

    is_blacklisted BOOLEAN DEFAULT FALSE,
    blacklist_reason TEXT,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_captain_verifications_captain_id
ON captain_verifications(captain_id);

CREATE TABLE captain_aadhar_details (
    id SERIAL PRIMARY KEY,
    verification_id INT NOT NULL REFERENCES captain_verifications(id) ON DELETE CASCADE,

    aadhar_number VARCHAR(20) NOT NULL,
    aadhar_name VARCHAR(255) NOT NULL,
    front_url TEXT,
    back_url TEXT,

    status document_status_enum DEFAULT 'PENDING',
    admin_comment TEXT,
    verified_by INT,
    verified_at TIMESTAMP,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_captain_aadhar_verification_id
ON captain_aadhar_details(verification_id);

CREATE TABLE captain_license_details (
    id SERIAL PRIMARY KEY,
    verification_id INT NOT NULL REFERENCES captain_verifications(id) ON DELETE CASCADE,

    license_number VARCHAR(50) NOT NULL,
    expiry_date DATE NOT NULL,
    front_url TEXT,
    back_url TEXT,

    status document_status_enum DEFAULT 'PENDING',
    admin_comment TEXT,
    verified_by INT,
    verified_at TIMESTAMP,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_captain_license_verification_id
ON captain_license_details(verification_id);

CREATE TABLE captain_vehicles (
    id SERIAL PRIMARY KEY,
    verification_id INT NOT NULL REFERENCES captain_verifications(id) ON DELETE CASCADE,

    vehicle_number VARCHAR(20) NOT NULL,
    vehicle_type VARCHAR(50),
    rc_book_url TEXT,
    insurance_url TEXT,

    status document_status_enum DEFAULT 'PENDING',
    admin_comment TEXT,
    verified_by INT,
    verified_at TIMESTAMP,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_captain_vehicles_verification_id
ON captain_vehicles(verification_id);


CREATE TABLE captain_bank_accounts (
    id SERIAL PRIMARY KEY,
    verification_id INT NOT NULL REFERENCES captain_verifications(id) ON DELETE CASCADE,

    account_number VARCHAR(50) NOT NULL,
    ifsc_code VARCHAR(20) NOT NULL,
    account_holder_name VARCHAR(255),

    status document_status_enum DEFAULT 'PENDING',
    admin_comment TEXT,
    verified_by INT,
    verified_at TIMESTAMP,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_captain_bank_verification_id
ON captain_bank_accounts(verification_id);

CREATE TABLE captain_background_checks (
    id SERIAL PRIMARY KEY,
    verification_id INT NOT NULL REFERENCES captain_verifications(id) ON DELETE CASCADE,

    police_verification_id VARCHAR(100),
    report_url TEXT,

    status document_status_enum DEFAULT 'PENDING',
    admin_comment TEXT,
    verified_by INT,
    verified_at TIMESTAMP,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_captain_background_verification_id
ON captain_background_checks(verification_id);

-- +goose StatementEnd



-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS captain_background_checks;
DROP TABLE IF EXISTS captain_bank_accounts;
DROP TABLE IF EXISTS captain_vehicles;
DROP TABLE IF EXISTS captain_license_details;
DROP TABLE IF EXISTS captain_aadhar_details;
DROP TABLE IF EXISTS captain_verifications;

DROP TYPE IF EXISTS verification_stage_enum;
DROP TYPE IF EXISTS document_status_enum;
DROP TYPE IF EXISTS verification_status_enum;

-- +goose StatementEnd
