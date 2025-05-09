BEGIN;

DROP TRIGGER set_updated_at_trigger ON expressions;
DROP TABLE expressions;
DROP TYPE expression_status;

COMMIT;