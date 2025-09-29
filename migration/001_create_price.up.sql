CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;

CREATE TABLE coin_prices (
    symbol VARCHAR(16) NOT NULL,
    price NUMERIC(30,10) NOT NULL,
    time BIGINT NOT NULL,
    PRIMARY KEY (symbol, time)
);
