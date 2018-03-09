CREATE TABLE IF NOT EXISTS "Users"(
	id bigserial PRIMARY KEY,
    "displayName" text NOT NULL,
    "emailAddress" text UNIQUE NOT NULL,
    password bytea NOT NULL,
    "creationTime" timestamp NOT NULL
);

CREATE TABLE IF NOT EXISTS "Publications"(
	id bigserial PRIMARY KEY,
    "authorId" bigint references "Users"(id) NOT NULL,
    "sequenceNumber" serial NOT NULL,
	"creationTime" timestamp NOT NULL

);

CREATE TABLE IF NOT EXISTS "Notes"(
	id bigserial PRIMARY KEY,
    "authorId" bigint references "Users"(id) NOT NULL,
    "type" text,
    "content" text NOT NULL,
    "publicationId" bigint references "Publications"(id),
    "creationTime" timestamp NOT NULL
);