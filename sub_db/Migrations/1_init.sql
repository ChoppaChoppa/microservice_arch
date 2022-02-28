CREATE TABLE OrderInfo (
                           id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                           track_number VARCHAR NOT NULL,
                           entry VARCHAR NOT NULL,
                           delivery UUID NOT NULL REFERENCES deliver(id),
                           payments UUID NOT NULL REFERENCES payment(id),
                           items UUID NOT NULL REFERENCES item(id),
                           locale VARCHAR NOT NULL,
                           internal_signature VARCHAR,
                           customer_id VARCHAR NOT NULL,
                           delivery_service VARCHAR NOT NULL,
                           shard_key VARCHAR NOT NULL,
                           sm_id INTEGER NOT NULL,
                           date_created TIMESTAMPTZ NOT NULL DEFAULT now(),
                           oof_shard VARCHAR NOT NULL
);

CREATE TABLE deliver (
                         id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                         name VARCHAR NOT NULL,
                         phone VARCHAR NOT NULL,
                         zip VARCHAR NOT NULL,
                         city VARCHAR NOT NULL,
                         address VARCHAR NOT NULL,
                         region VARCHAR NOT NULL,
                         email VARCHAR NOT NULL
);

CREATE TABLE payment (
                         id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                         transaction VARCHAR NOT NULL,
                         request_id VARCHAR,
                         currency VARCHAR,
                         provider VARCHAR,
                         amount INTEGER NOT NULL,
                         payment_dt INTEGER NOT NULL,
                         bank VARCHAR NOT NULL,
                         delivery_cost INTEGER NOT NULL,
                         goods_total INTEGER NOT NULL,
                         custom_fee INTEGER NOT NULL
);

CREATE TABLE item (
                      id UUID PRIMARY KEY DEFAULT  uuid_generate_v4(),
                      chrt_id INTEGER NOT NULL,
                      track_number VARCHAR NOT NULL,
                      price INTEGER NOT NULL,
                      rid VARCHAR NOT NULL,
                      name VARCHAR NOT NULL,
                      sale INTEGER NOT NULL,
                      size VARCHAR NOT NULL,
                      total_price INTEGER NOT NULL,
                      nm_id INTEGER NOT NULL,
                      brand VARCHAR NOT NULL,
                      status INTEGER NOT NULL
);