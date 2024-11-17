-- --------------------------------------------------------------------------------------------------

-- sequence name = table_name + "_" + primary_key_name + "_seq"
CREATE SEQUENCE printshop_global.notifier_queue_notice_id_seq START 1;

-- --------------------------------------------------------------------------------------------------

CREATE TABLE printshop_global.notifier_notices (
    notice_id int8 NOT NULL CONSTRAINT pk_notifier_notices PRIMARY KEY,
    notice_key character varying(128) NOT NULL,
    notice_data jsonb NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT NOW()
);

-- --------------------------------------------------------------------------------------------------

CREATE TABLE printshop_global.notifier_queue (
    notice_id int8 NOT NULL CONSTRAINT pk_notifier_queue PRIMARY KEY,
    remaining_attempts int2 NOT NULL, -- кол-во оставшихся попыток отправки сообщения
    item_status int2 NOT NULL, -- 1=READY, 2=PROCESSING, 3=RETRY
    updated_at timestamp with time zone NOT NULL DEFAULT NOW() -- item with status = READY and updated_at > NOW() = delayed
);

CREATE INDEX ix_notifier_queue_item_status ON printshop_global.notifier_queue  (item_status, updated_at);

-- --------------------------------------------------------------------------------------------------

CREATE TABLE printshop_global.notifier_queue_errors (
    notice_id int8 NOT NULL,
    error_message text NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT NOW()
);

CREATE INDEX ix_notifier_queue_errors_notice_id ON printshop_global.notifier_queue_errors (notice_id);
CREATE INDEX ix_notifier_queue_errors_created_at ON printshop_global.notifier_queue_errors (created_at);

-- --------------------------------------------------------------------------------------------------

CREATE TABLE printshop_global.notifier_queue_completed (
    notice_id int8 NOT NULL CONSTRAINT pk_notifier_queue_completed PRIMARY KEY,
    updated_at timestamp with time zone NOT NULL DEFAULT NOW()
);

CREATE INDEX ix_notifier_queue_completed_updated_at ON printshop_global.notifier_queue_completed (updated_at);
