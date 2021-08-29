DROP TABLE IF EXISTS public.table_with_bit_int64;

CREATE TABLE public.table_with_bit_int64
(
    id bigserial not null,
    type_id int NOT NULL,
    value varchar NOT NULL,
    source bigint NOT NULL,
    CONSTRAINT table_with_bit_int64_pkey PRIMARY KEY (id)
)
WITH (
    OIDS = FALSE
);

create unique index concurrently table_with_bit_int64_idx_1 on table_with_bit_int64(type_id, value);

COPY table_with_bit_int64 FROM '/home/admin/repo/bit_manipulation_research/script/result2-1.csv' DELIMITER ',' CSV HEADER;