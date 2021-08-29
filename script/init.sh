#!/usr/bin/env bash

cd $1;
location=$(pwd);

psql "postgresql://$2:$3@$4/$5" -c "
DROP TABLE IF EXISTS public.table_without_bit;
CREATE TABLE public.table_without_bit
(
    id bigserial not null,
    type_id int NOT NULL,
    value varchar NOT NULL,
    source int NOT NULL,
    status boolean not null
)
WITH (
    OIDS = FALSE
);
DROP TABLE IF EXISTS public.table_with_bit_int64;
CREATE TABLE public.table_with_bit_int64
(
    id bigserial not null,
    type_id int NOT NULL,
    value varchar NOT NULL,
    source bigint NOT NULL
)
WITH (
    OIDS = FALSE
);
DROP TABLE IF EXISTS public.table_with_bit_string;
CREATE TABLE public.table_with_bit_string
(
    id bigserial not null,
    type_id int NOT NULL,
    value varchar NOT NULL,
    source varchar NOT NULL
)
WITH (
    OIDS = FALSE
);
DROP TABLE IF EXISTS public.table_with_bit_bit;
CREATE TABLE public.table_with_bit_bit
(
    id bigserial not null,
    type_id int NOT NULL,
    value varchar NOT NULL,
    source bit(1000) NOT NULL
)
WITH (
    OIDS = FALSE
);
";

for file in result1*.csv
do
    echo $file;
    echo "execute $location/$file";
    echo "COPY table_without_bit FROM '$location/$file' DELIMITER ',' CSV HEADER;";
    psql "postgresql://$2:$3@$4/$5" -c "COPY table_without_bit FROM '$location/$file' DELIMITER ',' CSV HEADER;";
done

for file in result2*.csv
do
    echo $file;
    echo "execute $location/$file";
    echo "COPY table_with_bit_int64 FROM '$location/$file' DELIMITER ',' CSV HEADER;";
    psql "postgresql://$2:$3@$4/$5" -c "COPY table_with_bit_int64 FROM '$location/$file' DELIMITER ',' CSV HEADER;";
done

for file in result3*.csv
do
    echo $file;
    echo "execute $location/$file";
    echo "COPY table_with_bit_string FROM '$location/$file' DELIMITER ',' CSV HEADER;";
    psql "postgresql://$2:$3@$4/$5" -c "COPY table_with_bit_string FROM '$location/$file' DELIMITER ',' CSV HEADER;";
done

for file in result3*.csv
do
    echo $file;
    echo "execute $location/$file";
    echo "COPY table_with_bit_bit FROM '$location/$file' DELIMITER ',' CSV HEADER;";
    psql "postgresql://$2:$3@$4/$5" -c "COPY table_with_bit_bit FROM '$location/$file' DELIMITER ',' CSV HEADER;";
done

psql "postgresql://$2:$3@$4/$5" -c "
create unique index concurrently table_without_bit_idx_1 on table_without_bit(type_id, source, VALUE);
";

psql "postgresql://$2:$3@$4/$5" -c "
create index concurrently table_without_bit_idx_2 on table_without_bit(type_id, source, VALUE, status);
";

psql "postgresql://$2:$3@$4/$5" -c "
create unique index concurrently table_with_bit_int64_idx_1 on table_with_bit_int64(type_id, value);
";

psql "postgresql://$2:$3@$4/$5" -c "
CREATE UNIQUE INDEX concurrently table_with_bit_string_idx_1 on table_with_bit_string(type_id, value);
";

psql "postgresql://$2:$3@$4/$5" -c "
create unique index concurrently table_with_bit_bit_idx_1 on table_with_bit_bit(type_id, value);
";

psql "postgresql://$2:$3@$4/$5" -c "
ALTER TABLE table_without_bit ADD PRIMARY KEY (id);
";

psql "postgresql://$2:$3@$4/$5" -c "
ALTER TABLE table_with_bit_int64 ADD PRIMARY KEY (id);
";

psql "postgresql://$2:$3@$4/$5" -c "
ALTER TABLE table_with_bit_string ADD PRIMARY KEY (id);
";

psql "postgresql://$2:$3@$4/$5" -c "
ALTER TABLE table_with_bit_bit ADD PRIMARY KEY (id);
";

psql "postgresql://$2:$3@$4/$5" -c "VACUUM FULL ANALYZE;";

psql "postgresql://$2:$3@$4/$5" -c "SELECT pg_size_pretty( pg_total_relation_size('table_without_bit') ) as without_bit,
pg_size_pretty( pg_total_relation_size('table_with_bit_int64') ) as with_bit_int64,
pg_size_pretty( pg_total_relation_size('table_with_bit_string') ) as with_bit_string,
pg_size_pretty( pg_total_relation_size('table_with_bit_bit') ) as with_bit_bit;";