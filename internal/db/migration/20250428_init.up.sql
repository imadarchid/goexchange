CREATE TYPE order_side_type AS ENUM ('BUY', 'SELL');
CREATE TYPE order_status_type AS ENUM ('PENDING', 'SUBMITTED', 'CANCELED', 'PARTIALLY_FILLED', 'FILLED');
CREATE TYPE order_type AS ENUM ('MARKET', 'LIMIT');

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    dob DATE NOT NULL, 
    balance DOUBLE PRECISION NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE assets (
    id SERIAL PRIMARY KEY,
    ticker VARCHAR(5) NOT NULL UNIQUE,
    asset_name VARCHAR(255) NOT NULL,
    is_tradable BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    price DOUBLE PRECISION NOT NULL,
    amount INT NOT NULL,
    side order_side_type NOT NULL,
    order_type order_type NOT NULL,
    asset VARCHAR(5) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by UUID NOT NULL,
    order_status order_status_type NOT NULL DEFAULT 'SUBMITTED',
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE RESTRICT,
    FOREIGN KEY (asset) REFERENCES assets(ticker)
);

CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    price DOUBLE PRECISION NOT NULL,
    amount INT NOT NULL,
    buyer_order UUID NOT NULL,
    seller_order UUID NOT NULL,
    asset VARCHAR(5) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (buyer_order) REFERENCES orders(id),
    FOREIGN KEY (seller_order) REFERENCES orders(id),
    FOREIGN KEY (asset) REFERENCES assets(ticker)
);

