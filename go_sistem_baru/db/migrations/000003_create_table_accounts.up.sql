CREATE TABLE accounts (
    account_id VARCHAR(100)     PRIMARY KEY,
    username   VARCHAR(100)     UNIQUE NOT NULL DEFAULT '',
    password_hash TEXT           NOT NULL DEFAULT '',
    role       VARCHAR(50)      NOT NULL DEFAULT 'USER',
    balance    DECIMAL(15, 2)   NOT NULL DEFAULT 0.00,
    currency   VARCHAR(10)      NOT NULL DEFAULT 'IDR',
    version    INT              NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Seed initial accounts with username/password
INSERT INTO accounts (account_id, username, password_hash, role, balance, currency, version)
VALUES
    ('admin_001', 'admin', '$2b$10$lKYvxUUHomV6EmBURTGvUOvjJ5dRm9cIGaYYVun6JfO4QS82fOKKO', 'ADMIN', 0.00, 'IDR', 0),
    ('user_001', 'user', '$2b$10$IiqRmJFjMJ6lTc5AvO4ShOs4epOXDyVxAEqRZX2MmXmVyumSfsvmS', 'USER', 100000000.00, 'IDR', 0),
    ('merchant_001', 'merchant', '$2b$10$IkgX3zIP7GmnV/Pokng2Vux0DmM6YzNsErUfNlAgJ47tEFQFUBzZ.', 'MERCHANT', 0.00, 'IDR', 0);
