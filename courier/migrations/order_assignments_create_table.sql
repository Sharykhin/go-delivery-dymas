CREATE TABLE IF NOT EXISTS order_assignments (
                                                 courier_id UUID NOT NULL,
                                                 order_id UUID NOT NULL
                                                 created_at TIMESTAMPTZ NOT NULL,
                                                 PRIMARY KEY (order_id)
    );