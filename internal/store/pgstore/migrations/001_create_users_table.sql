-- Write your migrate up statements here

CREATE TABLE users (
	id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
	user_name VARCHAR(50) UNIQUE NOT NULL,
	email TEXT UNIQUE NOT NULL,
	password_hash BYTEA NOT NULL,
	bio TEXT NOT NULL DEFAULT '',
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

---- create above / drop below ----

DROP TABLE users;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
