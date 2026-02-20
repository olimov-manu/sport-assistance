-- +goose Up
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    surname VARCHAR(100) NOT NULL,
    gender VARCHAR(20) NOT NULL,
    birth_date DATE NOT NULL,

    height_cm INT,
    weight_kg INT,

    sport_activity_level_id INT REFERENCES sport_activity_levels(id) ON DELETE RESTRICT,
    sport_target_id INT REFERENCES sport_targets(id) ON DELETE RESTRICT,
    location_preference_type_id INT REFERENCES location_preference_types(id) ON DELETE RESTRICT,
    town_id INT REFERENCES towns(id) ON DELETE RESTRICT,
    role_id INT REFERENCES roles(id) ON DELETE RESTRICT,

    phone_number VARCHAR(30) NOT NULL UNIQUE,
    is_phone_verified BOOLEAN NOT NULL DEFAULT false,
    email VARCHAR(255) NOT NULL UNIQUE,
    is_email_verified BOOLEAN NOT NULL DEFAULT false,
    password TEXT,

    is_have_injury BOOLEAN DEFAULT false,
    injury_description TEXT,
    photo TEXT,

    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_users_town ON users(town_id);
CREATE INDEX idx_users_activity_level ON users(sport_activity_level_id);
CREATE INDEX idx_users_sport_target ON users(sport_target_id);
CREATE INDEX idx_users_location_preference_type ON users(location_preference_type_id);
CREATE INDEX idx_users_role ON users(role_id);

-- +goose Down
DROP TABLE users;
