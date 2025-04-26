-- +goose Up
CREATE TABLE users (
	id UUID PRIMARY KEY,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	email TEXT NOT NULL UNIQUE,
	password TEXT NOT NULL
);

CREATE INDEX idx_email ON users(email);

-- +goose Down
DROP INDEX IF EXISTS idx_email;
DROP TABLE users;

