CREATE DATABASE orders;
GRANT ALL PRIVILEGES ON DATABASE orders TO citizix_user;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
create type order_status as enum ('pending', 'accepted', 'in_progress', 'delivered', 'canceled');
CREATE TABLE IF NOT EXISTS orders (
                                          id UUID DEFAULT gen_random_uuid(),
                                          courier_id UUID NULL,
                                          customer_phone_number char(15) NOT NULL,
                                          status order_status DEFAULT 'pending',
                                          created_at TIMESTAMPTZ NOT NULL,
                                          PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS order_validations (
    order_id UUID NOT NULL,
    courier TIMESTAMPTZ,
    courier_error VARCHAR(256),
    created_at TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (order_id)
);
