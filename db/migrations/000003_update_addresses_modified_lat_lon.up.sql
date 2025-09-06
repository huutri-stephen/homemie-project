-- migrate_up_addresses.sql
BEGIN;

ALTER TABLE addresses
    ALTER COLUMN latitude TYPE DECIMAL(12,8) USING latitude::DECIMAL(12,8);

ALTER TABLE addresses
    ALTER COLUMN longitude TYPE DECIMAL(12,8) USING longitude::DECIMAL(12,8);

COMMIT;
