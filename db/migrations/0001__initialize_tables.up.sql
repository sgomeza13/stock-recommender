START TRANSACTION;

CREATE TABLE IF NOT EXISTS stock(
    id SERIAL PRIMARY KEY,
    ticker TEXT NOT NULL,
    target_from DECIMAL(10,2),
    target_to DECIMAL(10,2),
    company TEXT NOT NULL,
    action TEXT NOT NULL,
    brokerage TEXT NOT NULL,
    rating_from TEXT NOT NULL,
    rating_to TEXT NOT NULL,
    time TIMESTAMP WITH TIME ZONE NOT NULL
);

COMMIT;