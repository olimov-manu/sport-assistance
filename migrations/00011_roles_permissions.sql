-- +goose Up
CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE permissions (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE role_permissions (
    role_id INT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id INT NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    PRIMARY KEY (role_id, permission_id)
);

ALTER TABLE users
    ADD COLUMN role_id INT REFERENCES roles(id) ON DELETE RESTRICT;

CREATE INDEX idx_users_role ON users(role_id);

-- +goose Down
DROP INDEX IF EXISTS idx_users_role;

ALTER TABLE users
    DROP COLUMN IF EXISTS role_id;

DROP TABLE IF EXISTS role_permissions, permissions, roles;
