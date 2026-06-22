-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    login VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);


CREATE TABLE IF NOT EXISTS orders (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) NOT NULL,
    order_number VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL,
    accrual DECIMAL(10,2) DEFAULT 0,
    uploaded_at TIMESTAMP DEFAULT NOW(),
    UNIQUE (order_number)
);


CREATE TABLE IF NOT EXISTS balances (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) UNIQUE NOT NULL,
    current_balance DECIMAL(12,2) DEFAULT 0, 
    withdrawn_balance DECIMAL(12,2) DEFAULT 0
);


CREATE TABLE IF NOT EXISTS withdrawals (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) NOT NULL,
    order_number VARCHAR(255) NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    processed_at TIMESTAMP DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS withdrawals;
DROP TABLE IF EXISTS balances;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
