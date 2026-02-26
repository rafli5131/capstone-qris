CREATE TYPE transaction_status AS ENUM ('PENDING', 'SUCCESS', 'FAILED', 'REVERSED');

CREATE TABLE transactions (
    transaction_id VARCHAR(36)        PRIMARY KEY,
    trace_id       VARCHAR(100)       NOT NULL DEFAULT '',
    account_id     VARCHAR(100)       NOT NULL REFERENCES accounts(account_id),
    merchant_id    VARCHAR(100)       NOT NULL REFERENCES merchants(merchant_id),
    amount         DECIMAL(15, 2)     NOT NULL,
    status         transaction_status NOT NULL DEFAULT 'PENDING',
    created_at     TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_transactions_trace_id   ON transactions(trace_id);
CREATE INDEX idx_transactions_account_id ON transactions(account_id);
CREATE INDEX idx_transactions_merchant_id ON transactions(merchant_id);
CREATE INDEX idx_transactions_status     ON transactions(status);
