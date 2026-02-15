-- Convert empty strings to NULL to avoid unique constraint violations
UPDATE users SET phone_number = NULL WHERE phone_number = '';
UPDATE users SET google_id = NULL WHERE google_id = '';
UPDATE users SET line_id = NULL WHERE line_id = '';
