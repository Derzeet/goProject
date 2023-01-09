CREATE TABLE IF NOT EXISTS directors (
    -- id column is a 64-bit auto-incrementing integer & primary key (defines the row)
    id bigserial PRIMARY KEY,
    name text not null,
    surname text not null,
    -- genres column is array of zero-or-more text values. 
    awards text[]
);
