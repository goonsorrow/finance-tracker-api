BEGIN;

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,  
    name VARCHAR(100) NOT NULL,
    type VARCHAR(10) NOT NULL CHECK (type IN ('income', 'expense')),
    icon VARCHAR(255),
    usage_count BIGINT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, name, type)
);


INSERT INTO categories (name, type, user_id, usage_count) VALUES
('Зарплата', 'income', NULL, 0),('Фриланс', 'income', NULL, 0),
('Инвестиции', 'income', NULL, 0),('Аренда', 'income', NULL, 0),
('Продукты', 'expense', NULL, 0),('Кафе', 'expense', NULL, 0),
('Транспорт', 'expense', NULL, 0),('Коммуналка', 'expense', NULL, 0),
('Связь', 'expense', NULL, 0),('Развлечения', 'expense', NULL, 0),
('Здоровье', 'expense', NULL, 0),('Одежда', 'expense', NULL, 0);

CREATE TABLE wallets (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'RUB',
    balance DECIMAL(15,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    wallet_id INT NOT NULL REFERENCES wallets(id) ON DELETE CASCADE,
    user_id INT NOT NULL REFERENCES users(id),
    type VARCHAR(10) NOT NULL CHECK (type IN ('income','expense','initial')),
    amount DECIMAL(15,2) NOT NULL,
    category_id INT REFERENCES categories(id),  
    description TEXT,
    date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_wallets_user ON wallets(user_id);
CREATE INDEX idx_transactions_user ON transactions(user_id);
CREATE INDEX idx_transactions_date ON transactions(date);
CREATE INDEX idx_categories_user ON categories(user_id);
CREATE INDEX idx_categories_usage ON categories(usage_count DESC);

COMMIT;
