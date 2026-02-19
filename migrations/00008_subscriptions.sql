-- +goose Up
CREATE TABLE subscriptions (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE subscription_privileges (
    subscription_id INT REFERENCES subscriptions(id) ON DELETE CASCADE,
    privilege_id INT REFERENCES privileges(id) ON DELETE CASCADE,
    PRIMARY KEY (subscription_id, privilege_id)
);

-- +goose Down
DROP TABLE subscription_privileges, subscriptions;
