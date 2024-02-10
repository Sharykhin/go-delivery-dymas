CREATE DATABASE IF NOT EXISTS couriers;
GRANT ALL PRIVILEGES ON DATABASE couriers TO citizix_user;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS courier (
          id UUID DEFAULT gen_random_uuid(),
          first_name char(30) NOT NULL,
          is_available BOOLEAN DEFAULT TRUE,
          PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS order_assignments (
        courier_id UUID NOT NULL,
        order_id UUID NOT NULL
        created_at TIMESTAMPTZ NOT NULL,
        PRIMARY KEY (order_id)
    );