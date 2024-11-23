CREATE TABLE IF NOT EXISTS users
(
    id       SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    balance  FLOAT DEFAULT 0.0
);

CREATE TABLE IF NOT EXISTS transactions
(
    id               SERIAL PRIMARY KEY,
    user_id          INT REFERENCES users (id) ON DELETE CASCADE,
    amount           FLOAT       NOT NULL,
    transaction_type VARCHAR(10) NOT NULL,
    created_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);