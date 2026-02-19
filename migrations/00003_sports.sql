-- +goose Up
CREATE TABLE user_sports (
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    sport_id INT REFERENCES sports(id) ON DELETE RESTRICT,
    PRIMARY KEY (user_id, sport_id)
);

-- +goose Down
DROP TABLE user_sports;
