DROP TRIGGER IF EXISTS trigger_update_wallets_updated_at ON wallets;
DROP TRIGGER IF EXISTS trigger_update_transactions_updated_at ON transactions;
DROP TRIGGER IF EXISTS trigger_update_categories_updated_at ON categories;

DROP FUNCTION IF EXISTS update_updated_at_column();

ALTER TABLE wallets DROP COLUMN IF EXISTS updated_at;
ALTER TABLE transactions DROP COLUMN IF EXISTS updated_at;
ALTER TABLE categories DROP COLUMN IF EXISTS updated_at;
