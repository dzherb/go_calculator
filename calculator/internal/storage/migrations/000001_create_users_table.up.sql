BEGIN;

CREATE TABLE users
(
    id            SERIAL PRIMARY KEY,
    username      VARCHAR(32) UNIQUE        NOT NULL,
    password_hash CHAR(60)                  NOT NULL,
    created_at    TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at    TIMESTAMPTZ DEFAULT now() NOT NULL
);

CREATE
OR REPLACE FUNCTION set_updated_at()
    RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at
= now();
RETURN NEW;
END;
$$
LANGUAGE plpgsql;

CREATE TRIGGER set_updated_at_trigger
    BEFORE UPDATE
    ON users
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();

COMMIT;
