-- +goose Up
CREATE TABLE activities (
    id SERIAL PRIMARY KEY,
    service_id INT REFERENCES services(id) ON DELETE RESTRICT,
    name TEXT NOT NULL
);

CREATE TABLE activity_calendars (
    id SERIAL PRIMARY KEY,
    activity_id INT REFERENCES activities(id) ON DELETE CASCADE,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL
);

CREATE TABLE user_activity_calendars (
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    calendar_id INT REFERENCES activity_calendars(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, calendar_id)
);

CREATE INDEX idx_calendar_activity ON activity_calendars(activity_id);

-- +goose Down
DROP TABLE user_activity_calendars, activity_calendars, activities;
