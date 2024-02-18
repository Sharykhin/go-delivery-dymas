CREATE DATABASE IF NOT EXISTS couriers;
GRANT ALL PRIVILEGES ON DATABASE couriers TO citizix_user;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS courier (
          id UUID DEFAULT gen_random_uuid(),
          first_name char(30) NOT NULL,
          is_available BOOLEAN DEFAULT TRUE,
          PRIMARY KEY (id)
);
