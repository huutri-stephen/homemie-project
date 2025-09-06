-- migrate_down_addresses.sql
BEGIN;

ALTER TABLE addresses
    ALTER COLUMN latitude TYPE DECIMAL(10,8) USING latitude::DECIMAL(10,8);

ALTER TABLE addresses
    ALTER COLUMN longitude TYPE DECIMAL(11,8) USING longitude::DECIMAL(11,8);

COMMIT;
