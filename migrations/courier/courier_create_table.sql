CREATE DATABASE IF NOT EXISTS couriers;
GRANT ALL PRIVILEGES ON DATABASE courier_location TO citizix_user;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS couriers (
                                          id UUID NOT NULL,
                                          first_name char(30) NOT NULL,
                                          is_available BOOLEAN NOT NULL,
                                          PRIMARY KEY (id)
);
