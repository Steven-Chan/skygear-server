BEGIN;

ALTER TABLE _auth_principal DROP CONSTRAINT _auth_principal_user_id_provider_key;

END;