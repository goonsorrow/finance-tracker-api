BEGIN;

ALTER TABLE wallets 
RENAME COLUMN computed_balance TO balance;

ALTER TABLE wallets
DROP COLUMN IF EXISTS last_time_updated_at;

COMMIT;


