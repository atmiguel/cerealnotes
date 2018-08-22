-- Types
CREATE TYPE note_type AS ENUM ('predictions', 'marginalia', 'meta', 'questions');

-- Tables
CREATE TABLE IF NOT EXISTS app_user (
	id bigserial PRIMARY KEY,
	display_name text NOT NULL,
	email_address text UNIQUE NOT NULL,
	password bytea NOT NULL,
	creation_time timestamp NOT NULL
);

CREATE TABLE IF NOT EXISTS publication (
	id bigserial PRIMARY KEY,
	author_id bigint references app_user(id) NOT NULL,
	creation_time timestamp NOT NULL
);

CREATE TABLE IF NOT EXISTS note (
	id bigserial PRIMARY KEY,
	author_id bigint references app_user(id) NOT NULL,
	content text NOT NULL,
	creation_time timestamp NOT NULL
);

CREATE TABLE IF NOT EXISTS publication_to_note_relationship (
	id bigserial PRIMARY KEY,
	publication_id bigint references publication(id) NOT NULL,
	note_id bigint UNIQUE references note(id) NOT NULL
);

CREATE TABLE IF NOT EXISTS notetype_to_note_relationship (
	id bigserial PRIMARY KEY,
	type note_type,
	note_id bigint UNIQUE references note(id) NOT NULL
);