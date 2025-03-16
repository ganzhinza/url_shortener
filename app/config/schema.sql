CREATE TABLE IF NOT EXISTS urls(
    original VARCHAR(4096) PRIMARY KEY,
    short CHAR(10)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_short ON urls (short);
CREATE UNIQUE INDEX IF NOT EXISTS idx_original ON urls (original);