-- --------------------------------------------------------------------------------------------------

CREATE TABLE printshop_catalog.laminates (
    laminate_id int8 NOT NULL GENERATED BY DEFAULT AS IDENTITY CONSTRAINT pk_laminates PRIMARY KEY,
    tag_version int4 NOT NULL DEFAULT 1 CHECK(tag_version > 0),
    laminate_article character varying(32) NULL,
    laminate_caption character varying(64) NOT NULL,
    type_id int8 NOT NULL CHECK(type_id > 0),
    laminate_length double precision NOT NULL, -- meter
    laminate_width double precision NOT NULL, -- meter
    laminate_thickness double precision NOT NULL, -- meter
    laminate_weight_m2 double precision NOT NULL, -- kilogram
    laminate_status int2 NOT NULL, -- 1=DRAFT, 2=ENABLED, 3=DISABLED
    created_at timestamp with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp with time zone NOT NULL DEFAULT NOW(),
    deleted_at timestamp with time zone NULL
);

CREATE UNIQUE INDEX uk_laminates_laminate_article ON printshop_catalog.laminates (laminate_article) WHERE deleted_at IS NULL;
CREATE INDEX ix_laminates_type_id ON printshop_catalog.laminates (type_id) WHERE deleted_at IS NULL;
CREATE INDEX ix_laminates_laminate_thickness ON printshop_catalog.laminates (laminate_thickness) WHERE deleted_at IS NULL;

-- --------------------------------------------------------------------------------------------------

INSERT INTO printshop_catalog.laminates (laminate_id, tag_version, laminate_article, laminate_caption, type_id, laminate_length, laminate_width, laminate_thickness, laminate_weight_m2, laminate_status, created_at, updated_at, deleted_at)
VALUES
    (1, 1, 'lam-1', 'Глянцевый 450', 8, 10.0, 0.450, 0.00003, 1.0, 2/*ENABLED*/, '2023-07-26 21:08:40.908057', '2023-07-26 21:08:40.908057', NULL);

ALTER SEQUENCE printshop_catalog.laminates_laminate_id_seq RESTART WITH 2;