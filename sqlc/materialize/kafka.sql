
CREATE CONNECTION kafka_conn TO KAFKA (
    BROKER 'redpanda:29092',
    SECURITY PROTOCOL PLAINTEXT
);

CREATE SOURCE test_raw FROM KAFKA CONNECTION kafka_conn (
    TOPIC 'test'
)
FORMAT JSON
ENVELOPE NONE;

CREATE CONNECTION pg_conn
TO POSTGRES (
    HOST 'timescaledb',
    PORT 5432,
    USER 'postgres',
    PASSWORD 'postgres',
    DATABASE 'rapido'
);


CREATE SOURCE pg_source
FROM POSTGRES CONNECTION pg_conn (
    PUBLICATION 'mz_publication'
);
