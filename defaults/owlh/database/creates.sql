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

CREATE TABLE knownports (
    kp_id integer PRIMARY KEY AUTOINCREMENT,
    kp_uniqueid text NOT NULL,
    kp_param text NOT NULL,
    kp_value text NOT NULL
);

CREATE TABLE plugins (
    plugin_id integer PRIMARY KEY AUTOINCREMENT,
    plugin_uniqueid text NOT NULL,
    plugin_param text NOT NULL,
    plugin_value text NOT NULL
);

CREATE TABLE analyzer (
    analyzer_id integer PRIMARY KEY AUTOINCREMENT,
    analyzer_uniqueid text NOT NULL,
    analyzer_param text NOT NULL,
    analyzer_value text NOT NULL
);
CREATE TABLE nodeconfig (
    config_id integer PRIMARY KEY AUTOINCREMENT,
    config_uniqueid text NOT NULL,
    config_param text NOT NULL,
    config_value text NOT NULL
);
CREATE TABLE dataflow (
    flow_id integer PRIMARY KEY AUTOINCREMENT,
    flow_uniqueid text NOT NULL,
    flow_param text NOT NULL,
    flow_value text NOT NULL
);
CREATE TABLE mainconf (
    main_id integer PRIMARY KEY AUTOINCREMENT,
    main_uniqueid text NOT NULL,
    main_param text NOT NULL,
    main_value text NOT NULL
);
CREATE TABLE files (
    file_id integer PRIMARY KEY AUTOINCREMENT,
    file_uniqueid text NOT NULL,
    file_param text NOT NULL,
    file_value text NOT NULL
);