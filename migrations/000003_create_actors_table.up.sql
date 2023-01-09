CREATE TABLE IF NOT EXISTS actors (
    -- id column is a 64-bit auto-incrementing integer & primary key (defines the row)
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone not null default NOW(), 
    first_name text not null,
    last_name text not null,
    age integer not null,
    -- genres column is array of zero-or-more text values. 
    movies integer[] not NULL,
    version integer NOT NULL DEFAULT 1
);
