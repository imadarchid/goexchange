CREATE TYPE order_side_type AS ENUM ('BUY', 'SELL');
CREATE TYPE order_status_type AS ENUM ('PENDING', 'SUBMITTED', 'CANCELED', 'PARTIALLY_FILLED', 'FILLED');
CREATE TYPE order_type AS ENUM ('MARKET', 'LIMIT');

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(4),
    last_name VARCHAR(4),
    dob DATE NOT NULL, 
    balance DOUBLE PRECISION NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE assets (
    id SERIAL PRIMARY KEY,
    ticker VARCHAR(5) NOT NULL UNIQUE,
    asset_name VARCHAR(4) NOT NULL,
    is_tradable BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    price DOUBLE PRECISION NOT NULL,
    amount INT NOT NULL,
    side order_side_type NOT NULL,
    order_type order_type NOT NULL,
    asset VARCHAR(5) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by INT NOT NULL,
    order_status order_status_type NOT NULL DEFAULT 'PENDING',
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE RESTRICT,
    FOREIGN KEY (asset) REFERENCES assets(ticker)
);

CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    price DOUBLE PRECISION NOT NULL,
    amount INT NOT NULL,
    buyer_order INT NOT NULL,
    seller_order INT NOT NULL,
    asset VARCHAR(5) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (buyer_order) REFERENCES orders(id),
    FOREIGN KEY (seller_order) REFERENCES orders(id),
    FOREIGN KEY (asset) REFERENCES assets(ticker)
);

