CREATE TABLE servers (
    server_id integer PRIMARY KEY AUTOINCREMENT,
    server_uniqueid text NOT NULL,
    server_param text NOT NULL,
    server_value text NOT NULL
);
CREATE TABLE stap (
    stap_id integer PRIMARY KEY AUTOINCREMENT,
    stap_uniqueid text NOT NULL,
    stap_param text NOT NULL,
    stap_value text NOT NULL
);