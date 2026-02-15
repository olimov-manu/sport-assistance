-- +goose Up
CREATE TABLE chats (
                       id SERIAL PRIMARY KEY,
                       chat_type_id INT REFERENCES chat_types(id) ON DELETE RESTRICT,
                       created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE chat_members (
                              chat_id INT REFERENCES chats(id) ON DELETE CASCADE,
                              user_id INT REFERENCES users(id) ON DELETE CASCADE,
                              PRIMARY KEY (chat_id, user_id)
);

CREATE TABLE messages (
                          id SERIAL PRIMARY KEY,
                          chat_id INT REFERENCES chats(id) ON DELETE CASCADE,
                          sender_id INT REFERENCES users(id) ON DELETE SET NULL,
                          body TEXT NOT NULL,
                          created_at TIMESTAMP DEFAULT now()
);

CREATE INDEX idx_messages_chat ON messages(chat_id);
CREATE INDEX idx_messages_sender ON messages(sender_id);

-- +goose Down
DROP TABLE messages, chat_members, chats;
