CREATE DATABASE IF NOT EXISTS courier_location;
GRANT ALL PRIVILEGES ON DATABASE courier_location TO citizix_user;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS courier_latest_cord (
                                          courier_id UUID,
                                          latitude double precision,
                                          longitude double precision,
                                          created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                          PRIMARY KEY (courier_id, created_at)
);
