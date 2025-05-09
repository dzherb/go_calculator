BEGIN;

CREATE TYPE expression_status AS ENUM (
    'new',
    'processing',
    'succeed',
    'aborted',
    'failed'
);

CREATE TABLE expressions
(
    id         SERIAL PRIMARY KEY,
    user_id    INTEGER REFERENCES users (id) ON DELETE CASCADE,
    status     expression_status         NOT NULL,
    expression VARCHAR(256)              NOT NULL,
    result     DOUBLE PRECISION,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL
);

CREATE TRIGGER set_updated_at_trigger
    BEFORE UPDATE
    ON expressions
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();

COMMIT;
