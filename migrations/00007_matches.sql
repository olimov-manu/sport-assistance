-- +goose Up
CREATE TABLE matches (
                         id SERIAL PRIMARY KEY,
                         match_type_id INT REFERENCES match_types(id) ON DELETE RESTRICT,
                         created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE user_matches (
                              user_id INT REFERENCES users(id) ON DELETE CASCADE,
                              match_id INT REFERENCES matches(id) ON DELETE CASCADE,
                              PRIMARY KEY (user_id, match_id)
);

-- +goose Down
DROP TABLE user_matches, matches;
