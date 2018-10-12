\c cerealnotes_test;

-- Types
CREATE TYPE category_type AS ENUM ('predictions', 'marginalia', 'meta', 'questions');

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
	-- we need to have some sort of foreign key assurance that all note to publication relationships refer to the same author
	author_id bigint references app_user(id) NOT NULL,
	creation_time timestamp NOT NULL
);

CREATE TABLE IF NOT EXISTS note (
	id bigserial PRIMARY KEY,
	author_id bigint references app_user(id) ON DELETE CASCADE NOT NULL,
	content text NOT NULL,
	creation_time timestamp NOT NULL
);

CREATE TABLE IF NOT EXISTS note_to_publication_relationship (
	note_id bigint PRIMARY KEY references note(id) ON DELETE CASCADE,
	publication_id bigint references publication(id) ON DELETE CASCADE NOT NULL
);

CREATE TABLE IF NOT EXISTS note_to_category_relationship (
	note_id bigint PRIMARY KEY references note(id) ON DELETE CASCADE,
	category category_type NOT NULL
);