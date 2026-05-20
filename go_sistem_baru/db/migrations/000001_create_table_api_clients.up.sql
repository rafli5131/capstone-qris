CREATE TYPE api_client_status AS ENUM ('ACTIVE', 'INACTIVE');

CREATE TABLE api_clients (
    client_id   VARCHAR(100) PRIMARY KEY,
    client_secret VARCHAR(255) NOT NULL,
    status      api_client_status NOT NULL DEFAULT 'ACTIVE',
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Seed a default test client
INSERT INTO api_clients (client_id, client_secret, status)
VALUES ('MK-9921-X', 'super-secret-key-change-in-production', 'ACTIVE');
