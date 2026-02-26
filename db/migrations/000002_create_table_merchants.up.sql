CREATE TABLE merchants (
    merchant_id   VARCHAR(100) PRIMARY KEY,
    merchant_name VARCHAR(255) NOT NULL,
    mcc           VARCHAR(10)  NOT NULL DEFAULT '',
    city          VARCHAR(100) NOT NULL DEFAULT '',
    is_active     BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Seed sample merchant matching the QRIS sample payload
INSERT INTO merchants (merchant_id, merchant_name, mcc, city, is_active)
VALUES
    ('ID10232756067300303', 'M Ivan Store',  '2741', 'Jakarta Timur', TRUE),
    ('MICH-001',            'Toko Berkah Mandiri', '5411', 'Malang', TRUE);
