CREATE TABLE printshop_controls.submit_form_versions (
    form_id uuid NOT NULL CONSTRAINT fk_submit_form_versions_form_id REFERENCES printshop_controls.submit_forms (form_id) ON DELETE CASCADE,
    version int4 NOT NULL CHECK(version > 0),
    rewrite_name character varying(32) NULL,
    form_caption character varying(128) NOT NULL,
    form_detailing int2 NOT NULL, -- 1=NORMAL, 2=EXTENDED
    compiled_body jsonb NOT NULL,
    activity_status int2 NOT NULL, -- 1=DRAFT, 2=TESTING, 3=PUBLISHED, 4=ARCHIVED
    created_at timestamp with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp with time zone NULL,
    CONSTRAINT pk_submit_form_versions_form_id_version PRIMARY KEY (form_id, version)
);

CREATE UNIQUE INDEX uk_submit_form_versions_rewrite_name_activity_status ON printshop_controls.submit_form_versions (rewrite_name, activity_status) WHERE activity_status IN (2/*=TESTING*/, 3/*=PUBLISHED*/);

-- --------------------------------------------------------------------------------------------------

INSERT INTO printshop_controls.submit_form_versions (form_id, version, rewrite_name, form_caption, form_detailing, compiled_body, activity_status, created_at, updated_at)
VALUES ('aa22434f-f09d-4b20-89c0-0785948cdc04', 1, 'flyers-001', 'Флаеры', 1/*NORMAL*/, '[]', 3/*PUBLISHED*/, '2023-07-03 19:33:28.945816', NULL),
       ('aa22434f-f09d-4b20-89c0-0785948cdc04', 2, 'flyers-001', 'Флаеры', 1/*NORMAL*/, '[]', 2/*TESTING*/, '2023-07-03 19:33:28.945816', NULL);