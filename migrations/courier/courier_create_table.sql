CREATE DATABASE IF NOT EXISTS courier_location;
GRANT ALL PRIVILEGES ON DATABASE courier_location TO citizix_user;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS couriers (
                                          ID UUID NOT NULL,
                                          FirstName char(30) NOT NULL,
                                          IsAvailable BOOLEAN NOT NULL,
                                          PRIMARY KEY (ID)
);
