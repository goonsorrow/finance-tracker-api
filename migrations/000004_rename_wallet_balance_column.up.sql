BEGIN;

ALTER TABLE wallets 
RENAME COLUMN balance TO computed_balance;

ALTER TABLE wallets
ADD COLUMN last_computed_at TIMESTAMP WITH TIME ZONE DEFAULT NOW();

UPDATE wallets 
SET last_computed_at = NOW();

COMMIT;