-- Types
CREATE TYPE note_type AS ENUM ('uncategorized', 'predictions', 'marginalia', 'meta', 'questions');

-- Tables
CREATE TABLE IF NOT EXISTS user(
	id bigserial PRIMARY KEY,
	display_name text NOT NULL,
	email_address text UNIQUE NOT NULL,
	password bytea NOT NULL,
	creation_time timestamp NOT NULL
);

CREATE TABLE IF NOT EXISTS publication(
	id bigserial PRIMARY KEY,
	author_id bigint references user(id) NOT NULL,
	creation_time timestamp NOT NULL
);

CREATE TABLE IF NOT EXISTS note(
	id bigserial PRIMARY KEY,
	author_id bigint references user(id) NOT NULL,
	type note_type,
	content text NOT NULL,
	creation_time timestamp NOT NULL
);

CREATE TABLE IF NOT EXISTS publication_to_note_linkage(
	id bigserial PRIMARY KEY,
	publication_id bigint references publication(id) NOT NULL,
	note_id bigint references note(id) NOT NULL
);