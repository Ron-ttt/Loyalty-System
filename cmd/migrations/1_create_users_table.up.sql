BEGIN
CREATE TABLE IF NOT EXISTS users(
    id serial primary key, 
    login varchar(200) UNIQUE not null, 
    password text)

CREATE TABLE IF NOT EXISTS orders(
    users integer REFERENCES users(id), 
    order_id integer UNIQUE not null, 
    status as enum("NEW","REGISTERED", "PROCESSING", "INVALID", "PROCESSED") default "NEW", 
    bonus integer default 0, 
    created_at DATETIME default now())
COMMIT