CREATE DATABASE  wallet;
-- Wallets table
CREATE TABLE wallets (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL DEFAULT 0,
    balance DECIMAL(15, 8) NOT NULL DEFAULT 0.00000000,
    created_at INTEGER NOT NULL DEFAULT 0,
    updated_at INTEGER NOT NULL DEFAULT 0
);
COMMENT ON TABLE wallets IS 'user wallets table';
COMMENT ON COLUMN wallets.user_id IS 'user id';
COMMENT ON COLUMN wallets.balance IS 'user wallet balance amount';
CREATE INDEX idx_wallets_user_id ON wallets(user_id);

-- transactions table
CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    order_id VARCHAR(64) NOT NULL DEFAULT '',
    user_id INTEGER NOT NULL DEFAULT 0,
    tx_type INTEGER NOT NULL DEFAULT 0,
    amount DECIMAL(15, 8) NOT NULL DEFAULT 0.00000000,
    related_user_id INTEGER NOT NULL DEFAULT 0,
    created_at INTEGER NOT NULL DEFAULT 0,
    updated_at INTEGER NOT NULL DEFAULT 0
);
COMMENT ON TABLE transactions IS 'user transactions log table';
COMMENT ON COLUMN transactions.order_id IS 'order id';
COMMENT ON COLUMN transactions.user_id IS 'user id';
COMMENT ON COLUMN transactions.tx_type IS '0: undefine 1: deposit, 2: withdrawal, 3: transfer in, 4: transfer out';
COMMENT ON COLUMN transactions.amount IS 'transaction amount';
COMMENT ON COLUMN transactions.related_user_id IS 'related user id';
CREATE INDEX idx_transactions_order_id ON transactions(order_id);
CREATE INDEX idx_transactions_user_id ON transactions(user_id);
CREATE INDEX idx_transactions_related_user_id ON transactions(related_user_id);