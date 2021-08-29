DROP TABLE IF EXISTS public.table_without_bit;

CREATE TABLE public.table_without_bit
(
    id bigserial not null,
    type_id int NOT NULL,
    value varchar NOT NULL,
    source int NOT NULL,
    status boolean not null,
    CONSTRAINT table_without_bit_pkey PRIMARY KEY (id)
)
WITH (
    OIDS = FALSE
);

create unique index concurrently table_without_bit_idx_1 on table_without_bit(type_id, source, VALUE);
create index concurrently table_without_bit_idx_2 on table_without_bit(type_id, source, VALUE, status);

COPY table_without_bit FROM '/home/admin/repo/bit_manipulation_research/script/result1-1.csv' DELIMITER ',' CSV HEADER;