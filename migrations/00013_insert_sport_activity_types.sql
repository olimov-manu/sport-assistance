-- +goose Up
INSERT INTO sport_activity_levels (name)
VALUES
    ('низкий'),
    ('средний'),
    ('высокий')
ON CONFLICT (name) DO NOTHING;

-- +goose Down
UPDATE users
SET sport_activity_level_id = NULL
WHERE sport_activity_level_id IN (
    SELECT id
    FROM sport_activity_levels
    WHERE name IN ('низкий', 'средний', 'высокий')
);

DELETE FROM sport_activity_levels
WHERE name IN ('низкий', 'средний', 'высокий');
