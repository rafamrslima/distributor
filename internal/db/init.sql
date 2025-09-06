CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL,
    EmailCc TEXT NULL,
    content BYTEA NOT NULL,
    messageReceivedAt TIMESTAMPTZ DEFAULT now()
);