-- --------------------------------------------------------------------------------------------------

CREATE TABLE printshop_catalog.boxes (
    box_id int4 NOT NULL GENERATED BY DEFAULT AS IDENTITY CONSTRAINT pk_boxes PRIMARY KEY,
    tag_version int4 NOT NULL DEFAULT 1 CHECK(tag_version > 0),
    box_article character varying(32) NULL,
    box_caption character varying(64) NOT NULL,
    box_length int4 NOT NULL CHECK(box_length > 0),-- mkm (mm * 1000)
    box_width int4 NOT NULL CHECK(box_width > 0), -- mkm (mm * 1000)
    box_height int4 NOT NULL CHECK(box_height > 0), -- mkm (mm * 1000)
    box_weight int4 NOT NULL CHECK(box_weight > 0), -- mg (g * 1000)
    box_status int2 NOT NULL, -- 1=DRAFT, 2=ENABLED, 3=DISABLED
    created_at timestamp with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp with time zone NULL,
    deleted_at timestamp with time zone NULL
);

CREATE UNIQUE INDEX uk_boxes_box_article ON printshop_catalog.boxes (box_article) WHERE deleted_at IS NULL;

-- --------------------------------------------------------------------------------------------------

INSERT INTO printshop_catalog.boxes (box_id, tag_version, box_article, box_caption, box_length, box_width, box_height, box_weight, box_status, created_at, updated_at, deleted_at)
VALUES
    (1, 1, 'T-21-310x260x380', 'СДЭК', 310000, 260000, 380000, 1, 2/*ENABLED*/, '2023-07-28 19:47:00.917593', NULL, NULL),
    (2, 1, 'T-23-300x300x300', 'СДЭК', 300000, 300000, 300000, 15160, 2/*ENABLED*/, '2023-07-28 19:49:04.261215', NULL, NULL),
    (3, 1, 'T-23-310x230x195', 'СДЭК', 310000, 230000, 195000, 7400, 2/*ENABLED*/, '2023-07-30 12:28:57.095098', NULL, NULL);

ALTER SEQUENCE printshop_catalog.boxes_box_id_seq RESTART WITH 4;