-- +goose Up
INSERT INTO sport_targets (name)
VALUES
    ('похудение'),
    ('поддержание формы'),
    ('улучшение гибкости / тонуса'),
    ('повышение выносливости'),
    ('определю позже')
ON CONFLICT (name) DO NOTHING;

-- +goose Down
UPDATE users
SET sport_target_id = NULL
WHERE sport_target_id IN (
    SELECT id
    FROM sport_targets
    WHERE name IN (
        'похудение',
        'поддержание формы',
        'улучшение гибкости / тонуса',
        'повышение выносливости',
        'определю позже'
    )
);

DELETE FROM sport_targets
WHERE name IN (
    'похудение',
    'поддержание формы',
    'улучшение гибкости / тонуса',
    'повышение выносливости',
    'определю позже'
);
