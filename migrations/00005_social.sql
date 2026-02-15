-- +goose Up
CREATE TABLE friends (
                         user_id INT REFERENCES users(id) ON DELETE CASCADE,
                         friend_id INT REFERENCES users(id) ON DELETE CASCADE,
                         created_at TIMESTAMP DEFAULT now(),
                         CHECK (user_id < friend_id),
                         PRIMARY KEY (user_id, friend_id)
);

-- +goose Down
DROP TABLE friends;
