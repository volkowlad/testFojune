CREATE TABLE wallet (
    id SERIAL,
    walletid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    balance INTEGER
);