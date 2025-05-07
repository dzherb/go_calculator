BEGIN;

DROP TRIGGER set_updated_at_trigger ON users;
DROP FUNCTION set_updated_at();
DROP TABLE users;

COMMIT;