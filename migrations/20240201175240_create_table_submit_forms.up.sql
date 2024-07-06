-- --------------------------------------------------------------------------------------------------

CREATE TABLE printshop_controls.submit_forms (
    form_id uuid NOT NULL CONSTRAINT pk_submit_forms PRIMARY KEY,
    tag_version int4 NOT NULL DEFAULT 1 CHECK(tag_version > 0),
    rewrite_name character varying(32) NULL,
    param_name character varying(32) NULL,
    form_caption character varying(128) NOT NULL,
    form_detailing int2 NOT NULL, -- 1=NORMAL, 2=EXTENDED
    form_status int2 NOT NULL, -- 1=DRAFT, 2=ENABLED, 3=DISABLED
    created_at timestamp with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp with time zone NULL,
    deleted_at timestamp with time zone NULL
);

CREATE UNIQUE INDEX uk_submit_forms_rewrite_name ON printshop_controls.submit_forms (rewrite_name) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX uk_submit_forms_param_name ON printshop_controls.submit_forms (param_name) WHERE deleted_at IS NULL;

-- --------------------------------------------------------------------------------------------------

INSERT INTO printshop_controls.submit_forms (form_id, tag_version, rewrite_name, param_name, form_caption, form_detailing, form_status, created_at, updated_at, deleted_at)
VALUES ('aa22434f-f09d-4b20-89c0-0785948cdc04', 1, 'flyers-001', 'Flyers', 'Флаеры', 1/*NORMAL*/, 2/*ENABLED*/, '2023-07-03 19:33:28.945816', NULL, NULL);