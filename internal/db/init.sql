CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    client_name TEXT NOT NULL,
    report_name TEXT NOT NULL,
    email TEXT NOT NULL,
    content BYTEA NOT NULL,
    messageReceivedAt TIMESTAMPTZ DEFAULT now()
);