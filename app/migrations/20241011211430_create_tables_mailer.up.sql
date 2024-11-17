-- --------------------------------------------------------------------------------------------------

-- sequence name = table_name + "_" + primary_key_name + "_seq"
CREATE SEQUENCE printshop_global.mailer_queue_message_id_seq START 1;

-- --------------------------------------------------------------------------------------------------

CREATE TABLE printshop_global.mailer_messages (
    message_id int8 NOT NULL CONSTRAINT pk_mailer_messages PRIMARY KEY,
    message_channel character varying(128) NOT NULL,
    message_data jsonb NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT NOW()
);

-- {"header":{"CorrelationID":"56a8ee4a-7fcf-44c5-849e-e9f6a453e380"},"email":{"content_type":"text/plain","from":{"name":"Ivan Ivanov","email":"ivan.ivanov@localhost"},"to":{"name":"Ivan Ivanov","email":"ivan.ivanov@localhost"},"reply_to":{"name":"Ivan Ivanov","email":""},"subject":"Test Subject","content":"Test Content"}}

-- --------------------------------------------------------------------------------------------------

CREATE TABLE printshop_global.mailer_queue (
    message_id int8 NOT NULL CONSTRAINT pk_mailer_queue PRIMARY KEY,
    remaining_attempts int2 NOT NULL, -- кол-во оставшихся попыток отправки сообщения
    item_status int2 NOT NULL, -- 1=READY, 2=PROCESSING, 3=RETRY
    updated_at timestamp with time zone NOT NULL DEFAULT NOW() -- item with status = READY and updated_at > NOW() = delayed
);

CREATE INDEX ix_mailer_queue_item_status ON printshop_global.mailer_queue  (item_status, updated_at);

-- --------------------------------------------------------------------------------------------------

CREATE TABLE printshop_global.mailer_queue_errors (
    message_id int8 NOT NULL,
    error_message text NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT NOW()
);

CREATE INDEX ix_mailer_queue_errors_message_id ON printshop_global.mailer_queue_errors (message_id);
CREATE INDEX ix_mailer_queue_errors_created_at ON printshop_global.mailer_queue_errors (created_at);

-- --------------------------------------------------------------------------------------------------

CREATE TABLE printshop_global.mailer_queue_completed (
    message_id int8 NOT NULL CONSTRAINT pk_mailer_queue_completed PRIMARY KEY,
    updated_at timestamp with time zone NOT NULL DEFAULT NOW()
);

CREATE INDEX ix_mailer_queue_completed_updated_at ON printshop_global.mailer_queue_completed (updated_at);
