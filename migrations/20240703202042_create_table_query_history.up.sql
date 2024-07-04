-- --------------------------------------------------------------------------------------------------

CREATE TABLE printshop_calculation.query_history (
    query_id uuid NOT NULL CONSTRAINT pk_query_history PRIMARY KEY,
    query_caption character varying(64) NULL,
    query_params jsonb NOT NULL,
    query_result jsonb NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT NOW(),
    download_count int4 NOT NULL DEFAULT 0,
    downloaded_at timestamp with time zone NULL,
    deleted_at timestamp with time zone NULL
);

CREATE INDEX ix_query_history_created_at ON printshop_calculation.query_history (created_at) WHERE deleted_at IS NULL;
