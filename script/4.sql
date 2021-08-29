DROP TABLE IF EXISTS public.table_with_bit_bit;

CREATE TABLE public.table_with_bit_bit
(
    id bigserial not null,
    type_id int NOT NULL,
    value varchar NOT NULL,
    source bit(1000) NOT NULL,
    CONSTRAINT table_with_bit_bit_pkey PRIMARY KEY (id)
)
WITH (
    OIDS = FALSE
);

create unique index concurrently table_with_bit_bit_idx_1 on table_with_bit_bit(type_id, value);

COPY table_with_bit_bit FROM '' DELIMITER ',' CSV HEADER;