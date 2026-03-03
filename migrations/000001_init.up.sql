BEGIN;

-- Users
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    base_currency VARCHAR(3) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    total_computed_balance DECIMAL(15,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Categories
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(10) NOT NULL CHECK (type IN ('income', 'expense')),
    icon VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, name, type)
);

-- Default categories template
INSERT INTO categories (name, type) VALUES
('Зарплата', 'income'),
('Фриланс', 'income'),
('Инвестиции', 'income'),
('Аренда', 'income'),
('Продукты', 'expense'),
('Кафе', 'expense'),
('Транспорт', 'expense'),
('Коммуналка', 'expense'),
('Связь', 'expense'),
('Развлечения', 'expense'),
('Здоровье', 'expense'),
('Одежда', 'expense');

-- Wallets
CREATE TABLE wallets (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'RUB',
    balance BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Mos
CREATE TABLE movements (
    id SERIAL PRIMARY KEY,
    wallet_id INT NOT NULL REFERENCES wallets(id) ON DELETE CASCADE,
    user_id INT NOT NULL REFERENCES users(id),
    type VARCHAR(10) NOT NULL CHECK (type IN ('income','expense','initial')),
    amount BIGINT NOT NULL,
    category_id INT REFERENCES categories(id),
    description TEXT,
    date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Refresh tokens
CREATE TABLE refresh_tokens (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, token)
);

-- Indexes
CREATE INDEX idx_wallets_user ON wallets(user_id);
CREATE INDEX idx_movements_user ON movements(user_id);
CREATE INDEX idx_movements_date ON movements(date);
CREATE INDEX idx_categories_user ON categories(user_id);
CREATE INDEX idx_categories_usage ON categories(usage_count DESC);
CREATE INDEX idx_refresh_token ON refresh_tokens(token);

COMMIT;
