CREATE TABLE IF NOT EXISTS "Users"(
	id bigserial PRIMARY KEY,
    "displayName" text NOT NULL,
    "emailAddress" text UNIQUE NOT NULL,
    password bytea NOT NULL,
    "creationTime" timestamp NOT NULL
);