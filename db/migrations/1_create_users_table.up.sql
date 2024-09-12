BEGIN;
DROP TYPE IF EXISTS status_enum;
CREATE TYPE status_enum AS ENUM ('NEW', 'REGISTERED', 'PROCESSING', 'INVALID', 'PROCESSED');
CREATE TABLE IF NOT EXISTS users(
    id serial primary key, 
    login varchar(200) UNIQUE not null, 
    password text);

CREATE TABLE IF NOT EXISTS orders(
    id serial primary key,
    users_id integer REFERENCES users(id), 
    order_id bigint UNIQUE not null, 
    status status_enum DEFAULT 'NEW',
    bonus integer default 0, 
    created_at TIMESTAMP default now());
COMMIT;