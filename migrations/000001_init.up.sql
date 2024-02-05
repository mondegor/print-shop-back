-- --------------------------------------------------------------------------------------------------
-- --------------------------------------------------------------------------------------------------
-- --------------------------------------------------------------------------------------------------

CREATE SCHEMA printdata_provider_accounts AUTHORIZATION user_pg;

-- --------------------------------------------------------------------------------------------------

CREATE TABLE printdata_provider_accounts.companies_pages (
    account_id uuid NOT NULL CONSTRAINT PK_companies_pages PRIMARY KEY,
    updated_at timestamp NOT NULL DEFAULT now(),
    rewrite_name character varying(64) NULL CONSTRAINT UK_companies_pages_rewrite_name UNIQUE,
    page_head character varying(128) NOT NULL,
    logo_meta jsonb DEFAULT NULL,
    site_url character varying(256) NOT NULL,
    page_status int2 NOT NULL, -- 1=DRAFT, 2=HIDDEN', 3=PUBLISHED, 4=PUBLISHED_SHARED
    datetime_status timestamp NOT NULL DEFAULT now()
);

-- --------------------------------------------------------------------------------------------------
-- --------------------------------------------------------------------------------------------------
-- --------------------------------------------------------------------------------------------------

CREATE SCHEMA printdata_controls AUTHORIZATION user_pg;

-- --------------------------------------------------------------------------------------------------

CREATE TABLE printdata_controls.element_templates (
    template_id int4 NOT NULL GENERATED BY DEFAULT AS IDENTITY CONSTRAINT PK_element_templates PRIMARY KEY,
    tag_version int4 NOT NULL DEFAULT 1 CHECK(tag_version > 0),
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NULL,
    param_name character varying(32) NULL,
    template_caption character varying(64) NOT NULL,
    element_type int2 NOT NULL, -- 1=GROUP, 2=ELEMENT_LIST
    element_detailing int2 NOT NULL, -- 1=NORMAL, 2=EXTENDED
    element_body jsonb NOT NULL,
    template_status int2 NOT NULL -- 1=DRAFT, 2=ENABLED, 3=DISABLED, 4=REMOVED
);

INSERT INTO printdata_controls.element_templates (template_id, tag_version, created_at, updated_at, param_name, template_caption, element_type, element_detailing, element_body, template_status)
VALUES (1, 1, '2023-07-03 16:22:50.911157', NULL, 'Product', 'Поля листовой продукции', 2/*ELEMENT_LIST*/, 1/*NORMAL*/, '[
  {
    "id": "%parentId%Quantity",
    "name": "Тираж",
    "type": "number",
    "required": true,
    "view": "text",
    "values": [
      {
        "id": "%parentId%_value",
        "defaultValue": 1000,
        "minValue": 1,
        "maxValue": 1000000
      }
    ],
    "unit": "шт"
  },

  {
    "id": "%parentId%SimilarTypes",
    "name": "Видов",
    "type": "number",
    "required": true,
    "view": "text",
    "values": [
      {
        "id": "%parentId%_value",
        "defaultValue": 1,
        "minValue": 1,
        "maxValue": 100
      }
    ]
  },

  {
    "id": "%parentId%FormatX",
    "name": "Длина",
    "type": "number",
    "required": true,
    "view": "text",
    "values": [
      {
        "id": "%parentId%_value",
        "defaultValue": 297,
        "minValue": 1,
        "maxValue": 1020
      }
    ],
    "unit": "мм"
  },

  {
    "id": "%parentId%FormatY",
    "name": "Ширина",
    "type": "number",
    "required": true,
    "view": "text",
    "values": [
      {
        "id": "%parentId%_value",
        "defaultValue": 210,
        "minValue": 1,
        "maxValue": 1020
      }
    ],
    "unit": "мм"
  },

  {
    "id": "%parentId%PrintType",
    "type": "number",
    "name": "Вид печати",
    "required": true,
    "view": "radio",
    "values": [
      {
        "id": "%parentId%_Any",
        "name": "Любая печать"
      },
      {
        "id": "%parentId%_Offset",
        "name": "Офсетная печать"
      },
      {
        "id": "%parentId%_Digital",
        "name": "Цифровая печать"
      }
    ]
  }]', 2/*ENABLED*/),
(2, 1, '2023-07-03 16:34:02.369491', NULL, 'ProcessMedia', 'Бумага', 1/*GROUP*/, 1/*NORMAL*/, '[
  {
    "id": "%parentId%_Type",
    "name": "Тип бумаги",
    "type": "number",
    "required": true,
    "view": "radio",
    "dictionary": "media-type",
    "values": [
      {
        "id": "%parentId%_None",
        "name": "не указано"
      }
    ]
  },
  {
    "id": "%parentId%_Density",
    "name": "Плотность",
    "type": "number",
    "required": true,
    "view": "radio",
    "dictionary": "media-density",
    "values": [
      {
        "id": "%parentId%_None",
        "name": "не указано"
      }
    ],
    "unit": "г/м2"
  },
  {
    "id": "%parentId%_Texture",
    "name": "Фактура",
    "type": "number",
    "required": true,
    "view": "radio",
    "dictionary": "media-texture",
    "values": [
      {
        "id": "%parentId%_None",
        "name": "не указано"
      }
    ]
  },
  {
    "id": "%parentId%_Color",
    "name": "Цвет",
    "type": "number",
    "required": true,
    "view": "radio",
    "dictionary": "media-color",
    "values": [
      {
        "id": "%parentId%_None",
        "name": "не указано"
      }
    ]
  }
]', 2/*ENABLED*/),
(3, 1, '2023-07-03 16:38:59.254920', NULL, 'ProcessPackaging', 'Упаковка', 1/*GROUP*/, 1/*NORMAL*/, '[
  {
    "id": "%parentId%_Type",
    "name": "Тип упаковки",
    "type": "number",
    "required": true,
    "view": "radio",
    "values": [
      {
        "id": "%parentId%_ShrinkFilm",
        "name": "Термоусадочная пленка"
      },
      {
        "id": "%parentId%_CorrugatedBox",
        "name": "Гофрированная коробка"
      }
    ]
  }
]', 2/*ENABLED*/),
(4, 1, '2023-07-03 16:35:14.078093', NULL, 'ProcessPrinting', 'Печать', 1/*GROUP*/, 1/*NORMAL*/, '[
  {
    "id": "%parentId%_SideFace",
    "name": "Лицевая сторона",
    "type": "group",
    "required": true,
    "view": "block",
    "values": [
      {
        "id": "%parentId%_ColorMode",
        "name": "Количество цветов",
        "type": "number",
        "required": true,
        "view": "combo",
        "values": [
          {
            "id": "%parentId%_1",
            "name": "[ 1 ]"
          },
          {
            "id": "%parentId%_4",
            "name": "[ 4 ]"
          }
        ]
      },

      {
        "id": "%parentId%_Varnish",
        "name": "Лакировка",
        "type": "number",
        "required": false,
        "view": "combo",
        "dictionary": "varnish",
        "values": [
          {
            "id": "%parentId%_None",
            "name": "без лакировки"
          }
        ]
      }
    ]
  },

  {
    "id": "%parentId%_SideBack",
    "name": "Обратная сторона",
    "type": "group",
    "required": false,
    "view": "block",
    "values": [
      {
        "id": "%parentId%_ColorMode",
        "name": "Количество цветов",
        "type": "number",
        "required": true,
        "view": "combo",
        "values": [
          {
            "id": "%parentId%_0",
            "name": "[ 0 ]"
          },
          {
            "id": "%parentId%_1",
            "name": "[ 1 ]"
          },
          {
            "id": "%parentId%_4",
            "name": "[ 4 ]"
          }
        ]
      },

      {
        "id": "%parentId%_Varnish",
        "type": "number",
        "name": "Лакировка",
        "required": false,
        "view": "combo",
        "dictionary": "varnish",
        "values": [
          {
            "id": "%parentId%_None",
            "name": "без лакировки"
          }
        ]
      }
    ]
  }
]', 2/*ENABLED*/),
(5, 1, '2023-07-03 16:36:48.626009', NULL, 'ProcessLaminating', 'Ламинация', 1/*GROUP*/, 1/*NORMAL*/, '[
  {
    "id": "%parentId%_NumberOfSides",
    "name": "Количество сторон",
    "type": "number",
    "required": true,
    "view": "radio",
    "values": [
      {
        "id": "%parentId%_OneSide",
        "name": "Одна сторона"
      },
      {
        "id": "%parentId%_TwoSides",
        "name": "Две стороны"
      }
    ]
  },

  {
    "id": "%parentId%_LaminatingTexture",
    "name": "Тип ламината",
    "type": "number",
    "required": true,
    "view": "combo",
    "dictionary": "laminating-texture",
    "values": [
      {
        "id": "%parentId%_None",
        "name": "не указано"
      }
    ]
  },

  {
    "id": "%parentId%_LaminatingThikness",
    "name": "Толщина ламината",
    "type": "number",
    "required": true,
    "view": "combo",
    "dictionary": "laminating-thikness",
    "values": [
      {
        "id": "%parentId%_None",
        "name": "не указано"
      }
    ],
    "unit": "мм"
  }
]', 2/*ENABLED*/);

ALTER SEQUENCE printdata_controls.element_templates_template_id_seq RESTART WITH 6;

-- --------------------------------------------------------------------------------------------------

CREATE TABLE printdata_controls.forms (
    form_id int4 NOT NULL GENERATED BY DEFAULT AS IDENTITY CONSTRAINT PK_forms PRIMARY KEY,
    tag_version int4 NOT NULL DEFAULT 1 CHECK(tag_version > 0),
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NULL,
    param_name character varying(32) NULL CONSTRAINT UK_forms_param_name UNIQUE,
    form_caption character varying(128) NOT NULL,
    form_detailing int2 NOT NULL, -- 1=NORMAL, 2=EXTENDED
    form_body_compiled jsonb NOT NULL,
    form_status int2 NOT NULL -- 1=DRAFT, 2=ENABLED, 3=DISABLED, 4=REMOVED
);

INSERT INTO printdata_controls.forms (form_id, tag_version, created_at, updated_at, param_name, form_caption, form_detailing, form_body_compiled, form_status)
VALUES (1, 1, '2023-07-03 19:33:28.945816', NULL, 'Flyers', 'Флаеры', 1/*NORMAL*/, '[]', 2/*ENABLED*/);

ALTER SEQUENCE printdata_controls.forms_form_id_seq RESTART WITH 2;

-- --------------------------------------------------------------------------------------------------

CREATE TABLE printdata_controls.form_elements (
    element_id int4 NOT NULL GENERATED BY DEFAULT AS IDENTITY CONSTRAINT PK_form_elements PRIMARY KEY,
    tag_version int4 NOT NULL DEFAULT 1 CHECK(tag_version > 0),
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NULL,
    form_id int4 NOT NULL CONSTRAINT FK_form_elements_form_id REFERENCES printdata_controls.forms(form_id) ON DELETE CASCADE,
    param_name character varying(32) NULL,
    element_caption character varying(64) NOT NULL,
    template_id int4 NOT NULL CONSTRAINT FK_form_elements_template_id REFERENCES printdata_controls.element_templates (template_id),
    element_required bool NOT NULL,
    prev_field_id int4 NULL CHECK(prev_field_id IS NULL OR prev_field_id > 0),
    next_field_id int4 NULL CHECK(next_field_id IS NULL OR next_field_id > 0),
    order_field int8 NULL CHECK(order_field IS NULL OR order_field > 0),
    CONSTRAINT UK_form_elements_form_id_param_name UNIQUE (form_id, param_name)
);

CREATE INDEX IX_form_elements_form_id_order_field ON printdata_controls.form_elements (form_id, order_field);

INSERT INTO printdata_controls.form_elements (element_id, tag_version, created_at, updated_at, form_id, param_name, element_caption, template_id, element_required, prev_field_id, next_field_id, order_field)
VALUES (1, 1, '2023-07-15 13:49:58.567032', NULL, 1, 'Product', 'Поля листовой продукции', 1, true, null, null, null);

ALTER SEQUENCE printdata_controls.form_elements_element_id_seq RESTART WITH 2;

-- --------------------------------------------------------------------------------------------------
-- --------------------------------------------------------------------------------------------------
-- --------------------------------------------------------------------------------------------------

CREATE SCHEMA printdata_dictionaries AUTHORIZATION user_pg;

-- --------------------------------------------------------------------------------------------------

CREATE TABLE printdata_dictionaries.laminate_types (
    type_id int4 NOT NULL GENERATED BY DEFAULT AS IDENTITY CONSTRAINT PK_laminate_types PRIMARY KEY,
    tag_version int4 NOT NULL DEFAULT 1 CHECK(tag_version > 0),
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NULL,
    type_caption character varying(64) NOT NULL,
    type_status int4 NOT NULL -- 1=DRAFT, 2=ENABLED, 3=DISABLED, 4=REMOVED
);

INSERT INTO printdata_dictionaries.laminate_types (type_id, tag_version, created_at, updated_at, type_caption, type_status)
VALUES
    (1, 1, '2023-07-30 12:30:52.651613', NULL, 'глянцевый', 2/*ENABLED*/);

ALTER SEQUENCE printdata_dictionaries.laminate_types_type_id_seq RESTART WITH 2;

-- --------------------------------------------------------------------------------------------------

CREATE TABLE printdata_dictionaries.paper_colors (
    color_id int4 NOT NULL GENERATED BY DEFAULT AS IDENTITY CONSTRAINT PK_paper_colors PRIMARY KEY,
    tag_version int4 NOT NULL DEFAULT 1 CHECK(tag_version > 0),
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NULL,
    color_caption character varying(64) NOT NULL,
    color_status int2 NOT NULL -- 1=DRAFT, 2=ENABLED, 3=DISABLED, 4=REMOVED
);

INSERT INTO printdata_dictionaries.paper_colors (color_id, tag_version, created_at, updated_at, color_caption, color_status)
VALUES
    (1, 1, '2023-07-26 20:56:43.894664', NULL, 'белый', 2/*ENABLED*/);

ALTER SEQUENCE printdata_dictionaries.paper_colors_color_id_seq RESTART WITH 2;

-- --------------------------------------------------------------------------------------------------

CREATE TABLE printdata_dictionaries.paper_factures (
    facture_id int4 NOT NULL GENERATED BY DEFAULT AS IDENTITY CONSTRAINT PK_paper_factures PRIMARY KEY,
    tag_version int4 NOT NULL DEFAULT 1 CHECK(tag_version > 0),
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NULL,
    facture_caption character varying(64) NOT NULL,
    facture_status int2 NOT NULL -- 1=DRAFT, 2=ENABLED, 3=DISABLED, 4=REMOVED
);

INSERT INTO printdata_dictionaries.paper_factures (facture_id, tag_version, created_at, updated_at, facture_caption, facture_status)
VALUES
    (1, 2, '2023-07-26 20:53:18.942332', NULL, 'глянцевая', 2/*ENABLED*/),
    (2, 2, '2023-07-26 20:53:39.478106', NULL, 'матовая', 2/*ENABLED*/),
    (3, 1, '2023-07-29 13:30:58.387279', NULL, 'гладкая', 2/*ENABLED*/);

ALTER SEQUENCE printdata_dictionaries.paper_factures_facture_id_seq RESTART WITH 4;

-- --------------------------------------------------------------------------------------------------

CREATE TABLE printdata_dictionaries.print_formats (
    format_id int4 NOT NULL GENERATED BY DEFAULT AS IDENTITY CONSTRAINT PK_print_formats PRIMARY KEY,
    tag_version int4 NOT NULL DEFAULT 1 CHECK(tag_version > 0),
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NULL,
    format_caption character varying(64) NOT NULL,
    format_length int4 NOT NULL CHECK(format_length > 0 AND format_length < 1000000000), -- mm * 1000
    format_width int4 NOT NULL CHECK(format_width > 0 AND format_width < 1000000000), -- mm * 1000
    format_status int2 NOT NULL -- 1=DRAFT, 2=ENABLED, 3=DISABLED, 4=REMOVED
);

INSERT INTO printdata_dictionaries.print_formats (format_id, tag_version, created_at, updated_at, format_caption, format_length, format_width, format_status)
VALUES
    (1,1, '2023-01-01 00:00:00', NULL, 'A4 (210 x 297 mm)',210000,297000,2/*ENABLED*/),
    (2,1, '2023-01-01 00:00:00', NULL, 'A3 (297 x 420 mm)',297000,420000,2/*ENABLED*/),
    (3,1, '2023-01-01 00:00:00', NULL, '90 x 50 mm',90000,50000,2/*ENABLED*/),
    (4,1, '2023-01-01 00:00:00', NULL, 'A5 (148 x 210 mm)',148000,210000,2/*ENABLED*/),
    (5,1, '2023-01-01 00:00:00', NULL, 'A6 (105 x 148 mm)',105000,148000,2/*ENABLED*/),
    (6,1, '2023-01-01 00:00:00', NULL, '99 x 210 mm',99000,210000,2/*ENABLED*/),
    (7,1, '2023-01-01 00:00:00', NULL, '55 x 85 mm',55000,85000,2/*ENABLED*/),
    (8,1, '2023-01-01 00:00:00', NULL, '110 x 210 mm',110000,210000,2/*ENABLED*/),
    (9,1, '2023-01-01 00:00:00', NULL, '150 x 297 mm',150000,297000,2/*ENABLED*/),
    (10,1, '2023-01-01 00:00:00', NULL, '70 x 100 mm',70000,100000,2/*ENABLED*/),
    (11,1, '2023-01-01 00:00:00', NULL, 'Euro (110 x 220 mm)',110000,220000,2/*ENABLED*/),
    (12,1, '2023-01-01 00:00:00', NULL, 'C5 (162 x 229 mm)',162000,229000,2/*ENABLED*/),
    (13,1, '2023-01-01 00:00:00', NULL, 'C4 (229 x 324 mm)',229000,324000,2/*ENABLED*/),
    (14,1, '2023-01-01 00:00:00', NULL, '90 x 90 mm',90000,90000,2/*ENABLED*/),
    (15,1, '2023-01-01 00:00:00', NULL, '100 x 150 mm',100000,150000,2/*ENABLED*/),
    (16,1, '2023-01-01 00:00:00', NULL, '100 x 200 mm',100000,200000,2/*ENABLED*/),
    (17,1, '2023-01-01 00:00:00', NULL, '230 x 410 mm (47х65 см/3)',410000,230000,2/*ENABLED*/),
    (18,1, '2023-01-01 00:00:00', NULL, '520 x 360 mm',520000,360000,2/*ENABLED*/),
    (19,1, '2023-01-01 00:00:00', NULL, '250 x 330 mm',250000,330000,2/*ENABLED*/),
    (20,1, '2023-01-01 00:00:00', NULL, '520 x 360 mm (72x104/4)',520000,360000,2/*ENABLED*/),
    (23,1, '2023-01-01 00:00:00', NULL, 'SRA3 (320 x 450 mm)',320000,450000,2/*ENABLED*/),
    (24,1, '2023-01-01 00:00:00', NULL, 'max (330 x 485 mm)',330000,485000,2/*ENABLED*/),
    (25,1, '2023-01-01 00:00:00', NULL, '350 x 500 mm',350000,500000,2/*ENABLED*/),
    (26,1, '2023-01-01 00:00:00', NULL, '250 х 450 mm (70х100 см/6)',250000,450000,2/*ENABLED*/),
    (30,1, '2023-01-01 00:00:00', NULL, '320 х 450 mm (64х90 см/4)',320000,450000,2/*ENABLED*/),
    (31,1, '2023-01-01 00:00:00', NULL, '325 х 470 mm (65х47 см/2)',325000,470000,2/*ENABLED*/),
    (32,1, '2023-01-01 00:00:00', NULL, '650 x 470 mm',650000,470000,2/*ENABLED*/),
    (33,1, '2023-01-01 00:00:00', NULL, '640 x 450 mm',650000,450000,2/*ENABLED*/),
    (34,1, '2023-01-01 00:00:00', NULL, '500 x 700 mm',500000,700000,2/*ENABLED*/),
    (35,1, '2023-01-01 00:00:00', NULL, '620 x 470 mm',620000,470000,2/*ENABLED*/),
    (36,1, '2023-01-01 00:00:00', NULL, '520 x 720 mm (72x104/2)',520000,720000,2/*ENABLED*/),
    (37,1, '2023-01-01 00:00:00', NULL, '330 x 700 mm (70x100 см/3)',330000,700000,2/*ENABLED*/),
    (38,1, '2023-01-01 00:00:00', NULL, '700 x 1000 mm',700000,1000000,2/*ENABLED*/),
    (39,1, '2023-01-01 00:00:00', NULL, '640 x 900 mm',640000,900000,2/*ENABLED*/),
    (40,1, '2023-01-01 00:00:00', NULL, '620 x 940 mm',620000,940000,2/*ENABLED*/),
    (41,1, '2023-01-01 00:00:00', NULL, '720 х 1040 мм',720000,1040000,2/*ENABLED*/),
    (42,1, '2023-01-01 00:00:00', NULL, 'А1 (594 х 841 мм)',841000,594000,2/*ENABLED*/),
    (43,1, '2023-01-01 00:00:00', NULL, 'А2 (420 х 594 мм)',594000,420000,2/*ENABLED*/),
    (44,1, '2023-01-01 00:00:00', NULL, '330 x 350 mm (70x100 см/6)',350000,330000,2/*ENABLED*/),
    (45,1, '2023-01-01 00:00:00', NULL, '297 х 210 мм',297000,210000,2/*ENABLED*/),
    (46,1, '2023-01-01 00:00:00', NULL, '420 х 297 мм',420000,297000,2/*ENABLED*/);

ALTER SEQUENCE printdata_dictionaries.print_formats_format_id_seq RESTART WITH 47;

-- --------------------------------------------------------------------------------------------------
-- --------------------------------------------------------------------------------------------------
-- --------------------------------------------------------------------------------------------------

CREATE SCHEMA printdata_catalog AUTHORIZATION user_pg;

-- --------------------------------------------------------------------------------------------------

CREATE TABLE printdata_catalog.boxes (
    box_id int4 NOT NULL GENERATED BY DEFAULT AS IDENTITY CONSTRAINT PK_boxes PRIMARY KEY,
    tag_version int4 NOT NULL DEFAULT 1 CHECK(tag_version > 0),
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NULL,
    box_article character varying(32) NULL CONSTRAINT UK_boxes_box_article UNIQUE,
    box_caption character varying(64) NOT NULL,
    box_length int4 NOT NULL CHECK(box_length > 0 AND box_length < 10000001),-- mm * 1000
    box_width int4 NOT NULL CHECK(box_width > 0 AND box_width < 10000001), -- mm * 1000
    box_depth int4 NOT NULL CHECK(box_depth > 0 AND box_depth < 10000001), -- mm * 1000
    box_status int2 NOT NULL -- 1=DRAFT, 2=ENABLED, 3=DISABLED, 4=REMOVED
);

INSERT INTO printdata_catalog.boxes (box_id, tag_version, created_at, updated_at, box_article, box_caption, box_length, box_width, box_depth, box_status)
VALUES
    (1, 1, '2023-07-28 19:47:00.917593', NULL, 'T-21-310x26x380', 'СДЭК', 310, 260, 380, 2/*ENABLED*/),
    (2, 1, '2023-07-28 19:49:04.261215', NULL, 'T-23-300x300x300', 'СДЭК', 300, 300, 300, 2/*ENABLED*/),
    (3, 1, '2023-07-30 12:28:57.095098', NULL, 'T-23-310x230x195', 'СДЭК', 310, 230, 195, 2/*ENABLED*/);

ALTER SEQUENCE printdata_catalog.boxes_box_id_seq RESTART WITH 4;

-- --------------------------------------------------------------------------------------------------

CREATE TABLE printdata_catalog.laminates (
    laminate_id int4 NOT NULL GENERATED BY DEFAULT AS IDENTITY CONSTRAINT PK_laminates PRIMARY KEY,
    tag_version int4 NOT NULL DEFAULT 1 CHECK(tag_version > 0),
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NULL,
    laminate_article character varying(32) NULL CONSTRAINT UK_laminates_laminate_article UNIQUE,
    laminate_caption character varying(64) NOT NULL,
    type_id int4 NOT NULL CHECK(type_id > 0),
    laminate_length int4 NOT NULL CHECK(laminate_length > 0 AND laminate_length < 1000000001), -- mm * 1000
    laminate_weight int4 NOT NULL CHECK(laminate_weight > 0 AND laminate_weight < 10001), -- g/m2
    laminate_thickness int4 NOT NULL CHECK(laminate_thickness > 0 AND laminate_thickness < 1000001), -- mkm
    laminate_status int2 NOT NULL -- 1=DRAFT, 2=ENABLED, 3=DISABLED, 4=REMOVED
);

CREATE INDEX IX_laminates_type_id ON printdata_catalog.laminates (type_id);

INSERT INTO printdata_catalog.laminates (laminate_id, tag_version, created_at, updated_at, laminate_article, laminate_caption, type_id, laminate_length, laminate_weight, laminate_thickness, laminate_status)
VALUES
    (1, 1, '2023-07-26 21:08:40.908057', NULL, 'lam-1', 'Глянцевый 450', 1, 10000, 450, 30, 2/*ENABLED*/);

ALTER SEQUENCE printdata_catalog.laminates_laminate_id_seq RESTART WITH 2;

-- --------------------------------------------------------------------------------------------------

CREATE TABLE printdata_catalog.papers (
    paper_id int4 NOT NULL GENERATED BY DEFAULT AS IDENTITY CONSTRAINT PK_papers PRIMARY KEY,
    tag_version int4 NOT NULL DEFAULT 1 CHECK(tag_version > 0),
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NULL,
    paper_article character varying(32) NULL CONSTRAINT UK_papers_paper_article UNIQUE,
    paper_caption character varying(64) NOT NULL,
    color_id int4 NOT NULL CHECK(color_id > 0),
    facture_id int4 NOT NULL CHECK(facture_id > 0),
    paper_length int4 NOT NULL CHECK(paper_length > 0 AND paper_length < 10000001), -- mm * 1000
    paper_width int4 NOT NULL CHECK(paper_width > 0 AND paper_width < 10000001), -- mm * 1000
    paper_density int4 NOT NULL CHECK(paper_density > 0 AND paper_density < 10001), -- g/m2
    paper_thickness int4 NOT NULL CHECK(paper_thickness > 0 AND paper_thickness < 1000001), -- mm * 1000
    paper_sides int2 NOT NULL, -- 1=SAME, 2=DIFFERENT
    paper_status int2 NOT NULL -- 1=DRAFT, 2=ENABLED, 3=DISABLED, 4=REMOVED
);

CREATE INDEX IX_papers_color_id ON printdata_catalog.papers (color_id);
CREATE INDEX IX_papers_facture_id ON printdata_catalog.papers (facture_id);

INSERT INTO printdata_catalog.papers (paper_id, tag_version, created_at, updated_at, paper_article, paper_caption, color_id, facture_id, paper_length, paper_width, paper_density, paper_thickness, paper_sides, paper_status)
VALUES
    (1, 3, '2023-07-29 13:08:45.348283', NULL, 'mel130-64х90', 'мелованная', 1, 2, 900, 640, 130, 100, 2/*DIFFERENT*/, 2/*ENABLED*/),
    (2, 1, '2023-07-29 13:13:23.991912', NULL, 'v130-70х100', 'мелованная', 1, 2, 1000, 700, 130, 100, 2/*DIFFERENT*/, 2/*ENABLED*/),
    (3, 1, '2023-07-29 13:38:24.199768', NULL, 'offset-80-70х100', 'офсетная', 1, 3, 1000, 700, 80, 100, 1/*SAME*/, 2/*ENABLED*/),
    (4, 1, '2023-07-30 12:41:25.813125', NULL, 'offset-170-70х100', 'мелованная', 1, 2, 1000, 700, 170, 130, 1/*SAME*/, 2/*ENABLED*/),
    (5, 1, '2023-07-30 12:43:54.456152', NULL, 'coatGlossy-170-70х100', 'мелованная', 1, 1, 1000, 700, 170, 120, 1/*SAME*/, 2/*ENABLED*/);

ALTER SEQUENCE printdata_catalog.papers_paper_id_seq RESTART WITH 6;