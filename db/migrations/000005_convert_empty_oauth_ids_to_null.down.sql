-- Revert NULLs to empty strings (not recommended as it will violate UNIQUE constraint)
UPDATE users SET phone_number = '' WHERE phone_number IS NULL;
UPDATE users SET google_id = '' WHERE google_id IS NULL;
UPDATE users SET line_id = '' WHERE line_id IS NULL;
