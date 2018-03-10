-- Types
CREATE TYPE note_type AS ENUM ('predictions', 'marginalia', 'meta', 'questions');

-- Tables
CREATE TABLE IF NOT EXISTS users(
	id bigserial PRIMARY KEY,
	display_name text NOT NULL,
	email_address text UNIQUE NOT NULL,
	password bytea NOT NULL,
	creation_time timestamp NOT NULL
);

CREATE TABLE IF NOT EXISTS publications(
	id bigserial PRIMARY KEY,
	author_id bigint references users(id) NOT NULL,
	creation_time timestamp NOT NULL
);

CREATE TABLE IF NOT EXISTS notes(
	id bigserial PRIMARY KEY,
	author_id bigint references users(id) NOT NULL,
	type note_type,
	content text NOT NULL,
	publication_id bigint references publications(id),
	creation_time timestamp NOT NULL
);