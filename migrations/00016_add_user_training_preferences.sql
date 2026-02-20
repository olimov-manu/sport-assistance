-- +goose Up
CREATE TABLE user_preferred_locations (
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    location_name TEXT NOT NULL,
    PRIMARY KEY (user_id, location_name)
);

CREATE TABLE training_time_slots (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE user_training_time_slots (
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    training_time_slot_id INT NOT NULL REFERENCES training_time_slots(id) ON DELETE RESTRICT,
    PRIMARY KEY (user_id, training_time_slot_id)
);

INSERT INTO location_preference_types (name)
VALUES
    ('где угодно'),
    ('выбранные районы города')
ON CONFLICT (name) DO NOTHING;

INSERT INTO training_time_slots (name)
VALUES
    ('утро'),
    ('день'),
    ('вечер'),
    ('выходные')
ON CONFLICT (name) DO NOTHING;

-- +goose Down
DROP TABLE IF EXISTS user_training_time_slots;
DROP TABLE IF EXISTS training_time_slots;
DROP TABLE IF EXISTS user_preferred_locations;

UPDATE users
SET location_preference_type_id = NULL
WHERE location_preference_type_id IN (
    SELECT id
    FROM location_preference_types
    WHERE name IN ('где угодно', 'выбранные районы города')
);

DELETE FROM location_preference_types
WHERE name IN ('где угодно', 'выбранные районы города');
