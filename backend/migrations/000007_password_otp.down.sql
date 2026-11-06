DROP TABLE IF EXISTS crm_schema.auth_security_events;
DROP TABLE IF EXISTS crm_schema.password_reset_token;
DROP TABLE IF EXISTS crm_schema.otp_challenge;
DROP INDEX IF EXISTS crm_schema.uq_users_phone_number_present;
ALTER TABLE crm_schema.users DROP COLUMN IF EXISTS auth_version;

-- The old NOT NULL/unique constraint was intentionally replaced with a partial unique index.
-- Restore its previous behavior only if every user currently has a phone number.
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM crm_schema.users WHERE phone_number IS NULL) THEN
        RAISE EXCEPTION 'Cannot roll back 000007 while users.phone_number contains NULL; restore phone numbers first';
    END IF;
END $$;
ALTER TABLE crm_schema.users ALTER COLUMN phone_number SET NOT NULL;
ALTER TABLE crm_schema.users ADD CONSTRAINT users_phone_number_key UNIQUE (phone_number);
