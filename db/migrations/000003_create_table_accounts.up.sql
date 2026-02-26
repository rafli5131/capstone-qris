CREATE TABLE accounts (
    account_id VARCHAR(100)    PRIMARY KEY,
    balance    DECIMAL(15, 2)  NOT NULL DEFAULT 0.00,
    currency   VARCHAR(10)     NOT NULL DEFAULT 'IDR',
    version    INT             NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Seed a test user account
INSERT INTO accounts (account_id, balance, currency, version)
VALUES ('user_123', 100000000.00, 'IDR', 0);
