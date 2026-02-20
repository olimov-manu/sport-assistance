-- +goose Up
INSERT INTO sports (name)
VALUES
    ('теннис'),
    ('падел'),
    ('пилатес')
ON CONFLICT (name) DO NOTHING;

-- +goose Down
DELETE FROM user_sports
WHERE sport_id IN (
    SELECT id
    FROM sports
    WHERE name IN ('теннис', 'падел', 'пилатес')
);

DELETE FROM sports
WHERE name IN ('теннис', 'падел', 'пилатес');
