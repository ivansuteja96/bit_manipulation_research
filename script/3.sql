DROP TABLE IF EXISTS public.table_with_bit_string;

CREATE TABLE public.table_with_bit_string
(
    id bigserial not null,
    type_id int NOT NULL,
    value varchar NOT NULL,
    source varchar NOT NULL,
    CONSTRAINT table_with_bit_string_pkey PRIMARY KEY (id)
)
WITH (
    OIDS = FALSE
);

CREATE UNIQUE INDEX concurrently table_with_bit_string_idx_1 on table_with_bit_string(type_id, value);

COPY table_with_bit_string FROM '' DELIMITER ',' CSV HEADER;