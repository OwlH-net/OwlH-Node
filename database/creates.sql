CREATE TABLE servers (
    server_id integer PRIMARY KEY AUTOINCREMENT,
    server_uniqueid text NOT NULL,
    server_param text NOT NULL,
    server_value text NOT NULL
);