-- CREATE SCHEMA public AUTHORIZATION user_pg;

CREATE TYPE item_status AS ENUM (
    'DRAFT',
    'ENABLED',
    'DISABLED',
    'REMOVED');

CREATE TYPE form_detailing AS ENUM (
	'NORMAL',
	'EXTENDED');

CREATE TYPE form_field_type AS ENUM (
	'GROUP',
	'FIELDS');

CREATE TYPE catalog_paper_sides AS ENUM (
    'SAME',
    'DIFFERENT');

CREATE TABLE form_field_templates (
	template_id int4 NOT NULL GENERATED ALWAYS AS IDENTITY,
    tag_version int4 NOT NULL DEFAULT 1 CHECK(tag_version > 0),
    datetime_created timestamp NOT NULL DEFAULT now(),
	param_name varchar(32) NOT NULL,
	template_caption varchar(128) NOT NULL,
	field_type form_field_type NOT NULL,
    field_detailing form_detailing NOT NULL,
	field_body json NOT NULL,
    template_status item_status NOT NULL,
	CONSTRAINT form_field_templates_pkey PRIMARY KEY (template_id)
);

CREATE TABLE form_data (
	form_id int4 NOT NULL GENERATED ALWAYS AS IDENTITY,
    tag_version int4 NOT NULL DEFAULT 1 CHECK(tag_version > 0),
    datetime_created timestamp NOT NULL DEFAULT now(),
    form_caption varchar(128) NOT NULL,
	form_detailing form_detailing NOT NULL,
	form_body_compiled json NOT NULL,
    form_status item_status NOT NULL,
	CONSTRAINT form_data_pkey PRIMARY KEY (form_id)
);

CREATE TABLE form_fields (
    field_id int4 NOT NULL GENERATED ALWAYS AS IDENTITY,
    form_id int4 NOT NULL,
    tag_version int4 NOT NULL DEFAULT 1 CHECK(tag_version > 0),
    template_id int4 NOT NULL,
    datetime_created timestamp NOT NULL DEFAULT now(),
    param_name varchar(32) NULL,
    field_caption varchar(128) NOT NULL,
    field_required bool NOT NULL,
    field_status item_status NOT NULL,
    prev_field_id int4 NULL,
    next_field_id int4 NULL,
    order_field int8 NULL,
    CONSTRAINT form_fields_pkey PRIMARY KEY (field_id),
    CONSTRAINT form_data_fkey FOREIGN KEY (form_id)
        REFERENCES form_data(form_id) ON DELETE CASCADE,
    CONSTRAINT form_field_templates_fkey FOREIGN KEY (template_id)
        REFERENCES form_field_templates(template_id),
    CONSTRAINT form_fields_param_name UNIQUE (form_id, param_name)
);

CREATE INDEX form_fields_order_field ON form_fields (form_id, order_field);

INSERT INTO form_field_templates (tag_version, datetime_created, param_name, template_caption, field_type, field_detailing, field_body, template_status) VALUES (1, '2023-07-03 16:22:50.911157', 'ProductionFlyers', 'Поля листовой продукции', 'FIELDS', 'NORMAL', '[
  {
    "id": "ProductQuantity",
    "name": "Тираж",
    "type": "number",
    "required": true,
    "view": "text",
    "values": [
      {
        "id": "ProductQuantity_value",
        "defaultValue": 1000,
        "minValue": 1,
        "maxValue": 1000000
      }
    ],
    "unit": "шт"
  },

  {
    "id": "ProductSimilarTypes",
    "name": "Видов",
    "type": "number",
    "required": true,
    "view": "text",
    "values": [
      {
        "id": "ProductSimilarTypes_value",
        "defaultValue": 1,
        "minValue": 1,
        "maxValue": 100
      }
    ]
  },

  {
    "id": "ProductFormatX",
    "name": "Длина",
    "type": "number",
    "required": true,
    "view": "text",
    "values": [
      {
        "id": "ProductFormatX_value",
        "defaultValue": 297,
        "minValue": 1,
        "maxValue": 1020
      }
    ],
    "unit": "мм"
  },

  {
    "id": "ProductFormatY",
    "name": "Ширина",
    "type": "number",
    "required": true,
    "view": "text",
    "values": [
      {
        "id": "ProductFormatY_value",
        "defaultValue": 210,
        "minValue": 1,
        "maxValue": 1020
      }
    ],
    "unit": "мм"
  },

  {
    "id": "PrintType",
    "type": "number",
    "name": "Вид печати",
    "required": true,
    "view": "radio",
    "values": [
      {
        "id": "PrintType_PrintAny",
        "name": "Любая печать"
      },
      {
        "id": "PrintType_PrintOffset",
        "name": "Офсетная печать"
      },
      {
        "id": "PrintType_PrintDigital",
        "name": "Цифровая печать"
      }
    ]
  }]', 'ENABLED');

INSERT INTO form_field_templates (tag_version, datetime_created, param_name, template_caption, field_type, field_detailing, field_body, template_status) VALUES (1, '2023-07-03 16:34:02.369491', 'ProcessMedia', 'Бумага', 'GROUP', 'NORMAL', '[
  {
    "id": "ProcessMedia_MediaType",
    "name": "Тип бумаги",
    "type": "number",
    "required": true,
    "view": "radio",
    "dictionary": "media-type",
    "values": [
      {
        "id": "ProcessMedia_MediaType_None",
        "name": "не указано"
      }
    ]
  },
  {
    "id": "ProcessMedia_MediaDensity",
    "name": "Плотность",
    "type": "number",
    "required": true,
    "view": "radio",
    "dictionary": "media-density",
    "values": [
      {
        "id": "ProcessMedia_MediaDensity_None",
        "name": "не указано"
      }
    ],
    "unit": "г/м2"
  },
  {
    "id": "ProcessMedia_MediaTexture",
    "name": "Фактура",
    "type": "number",
    "required": true,
    "view": "radio",
    "dictionary": "media-texture",
    "values": [
      {
        "id": "ProcessMedia_MediaTexture_None",
        "name": "не указано"
      }
    ]
  },
  {
    "id": "ProcessMedia_MediaColor",
    "name": "Цвет",
    "type": "number",
    "required": true,
    "view": "radio",
    "dictionary": "media-color",
    "values": [
      {
        "id": "ProcessMedia_MediaColor_None",
        "name": "не указано"
      }
    ]
  }
]', 'ENABLED');

INSERT INTO form_field_templates (tag_version, datetime_created, param_name, template_caption, field_type, field_detailing, field_body, template_status) VALUES (1, '2023-07-03 16:38:59.254920', 'ProcessPackaging', 'Упаковка', 'GROUP', 'NORMAL', '[
  {
    "id": "ProcessPackaging_Type",
    "name": "Тип упаковки",
    "type": "number",
    "required": true,
    "view": "radio",
    "values": [
      {
        "id": "ProcessPackaging_Type_ShrinkFilm",
        "name": "Термоусадочная пленка"
      },
      {
        "id": "ProcessPackaging_Type_CorrugatedBox",
        "name": "Гофрированная коробка"
      }
    ]
  }
]', 'ENABLED');

INSERT INTO form_field_templates (tag_version, datetime_created, param_name, template_caption, field_type, field_detailing, field_body, template_status) VALUES (1, '2023-07-03 16:35:14.078093', 'ProcessPrinting', 'Печать', 'GROUP', 'NORMAL', '[
  {
    "id": "ProcessPrinting_SideFace",
    "name": "Лицевая сторона",
    "type": "group",
    "required": true,
    "view": "block",
    "values": [
      {
        "id": "ProcessPrinting_SideFace_ColorMode",
        "name": "Количество цветов",
        "type": "number",
        "required": true,
        "view": "combo",
        "values": [
          {
            "id": "ProcessPrinting_SideFace_ColorMode_1",
            "name": "[ 1 ]"
          },
          {
            "id": "ProcessPrinting_SideFace_ColorMode_4",
            "name": "[ 4 ]"
          }
        ]
      },

      {
        "id": "ProcessPrinting_SideFace_Varnish",
        "name": "Лакировка",
        "type": "number",
        "required": false,
        "view": "combo",
        "dictionary": "varnish",
        "values": [
          {
            "id": "ProcessPrinting_SideFace_Varnish_None",
            "name": "без лакировки"
          }
        ]
      }
    ]
  },

  {
    "id": "ProcessPrinting_SideBack",
    "name": "Обратная сторона",
    "type": "group",
    "required": false,
    "view": "block",
    "enabledValues": [
      {
        "id": "ProcessPrinting_SideBack_Disabled",
        "checked": false
      },
      {
        "id": "ProcessPrinting_SideBack_Enabled",
        "checked": true
      }
    ],
    "values": [
      {
        "id": "ProcessPrinting_SideBack_ColorMode",
        "name": "Количество цветов",
        "type": "number",
        "required": true,
        "view": "combo",
        "values": [
          {
            "id": "ProcessPrinting_SideBack_ColorMode_0",
            "name": "[ 0 ]"
          },
          {
            "id": "ProcessPrinting_SideBack_ColorMode_1",
            "name": "[ 1 ]"
          },
          {
            "id": "ProcessPrinting_SideBack_ColorMode_4",
            "name": "[ 4 ]"
          }
        ]
      },

      {
        "id": "ProcessPrinting_SideBack_Varnish",
        "type": "number",
        "name": "Лакировка",
        "required": false,
        "view": "combo",
        "dictionary": "varnish",
        "values": [
          {
            "id": "ProcessPrinting_SideBack_Varnish_None",
            "name": "без лакировки"
          }
        ]
      }
    ]
  }
]', 'ENABLED');

INSERT INTO form_field_templates (tag_version, datetime_created, param_name, template_caption, field_type, field_detailing, field_body, template_status) VALUES (1, '2023-07-03 16:36:48.626009', 'ProcessLaminating', 'Ламинация', 'GROUP', 'NORMAL', '[
  {
    "id": "ProcessLaminating_NumberOfSides",
    "name": "Количество сторон",
    "type": "number",
    "required": true,
    "view": "radio",
    "values": [
      {
        "id": "ProcessLaminating_NumberOfSides_OneSide",
        "name": "Одна сторона"
      },
      {
        "id": "ProcessLaminating_NumberOfSides_TwoSides",
        "name": "Две стороны"
      }
    ]
  },

  {
    "id": "ProcessLaminating_LaminatingTexture",
    "name": "Тип ламината",
    "type": "number",
    "required": true,
    "view": "combo",
    "dictionary": "laminating-texture",
    "values": [
      {
        "id": "ProcessLaminating_LaminatingTexture_None",
        "name": "не указано"
      }
    ]
  },

  {
    "id": "ProcessLaminating_LaminatingThikness",
    "name": "Толщина ламината",
    "type": "number",
    "required": true,
    "view": "combo",
    "dictionary": "laminating-thikness",
    "values": [
      {
        "id": "ProcessLaminating_LaminatingThikness_None",
        "name": "не указано"
      }
    ],
    "unit": "мм"
  }
]', 'ENABLED');

INSERT INTO form_data (tag_version, datetime_created, form_caption, form_detailing, form_body_compiled, form_status) VALUES (1, '2023-07-03 19:33:28.945816', 'Флаеры', 'NORMAL', '[]', 'ENABLED');


CREATE TABLE catalog_print_formats (
    format_id int4 NOT NULL GENERATED ALWAYS AS IDENTITY,
    tag_version int4 NOT NULL DEFAULT 1 CHECK(tag_version > 0),
    datetime_created timestamp NOT NULL DEFAULT now(),
    format_caption varchar(64) NOT NULL,
    format_length int4 NOT NULL CHECK(format_length > 0 AND format_length < 1000000000), -- mm * 1000
    format_width int4 NOT NULL CHECK(format_width > 0 AND format_width < 1000000000), -- mm * 1000
    format_status item_status NOT NULL,
    CONSTRAINT catalog_print_formats_pkey PRIMARY KEY (format_id)
);

INSERT INTO catalog_print_formats (format_id, tag_version, datetime_created, format_caption, format_length, format_width, format_status)
    OVERRIDING SYSTEM VALUE
VALUES
    (1,1, '2023-01-01 00:00:00','A4 (210 x 297 mm)',210000,297000,'ENABLED'),
    (2,1, '2023-01-01 00:00:00','A3 (297 x 420 mm)',297000,420000,'ENABLED'),
    (3,1, '2023-01-01 00:00:00','90 x 50 mm',90000,50000,'ENABLED'),
    (4,1, '2023-01-01 00:00:00','A5 (148 x 210 mm)',148000,210000,'ENABLED'),
    (5,1, '2023-01-01 00:00:00','A6 (105 x 148 mm)',105000,148000,'ENABLED'),
    (6,1, '2023-01-01 00:00:00','99 x 210 mm',99000,210000,'ENABLED'),
    (7,1, '2023-01-01 00:00:00','55 x 85 mm',55000,85000,'ENABLED'),
    (8,1, '2023-01-01 00:00:00','110 x 210 mm',110000,210000,'ENABLED'),
    (9,1, '2023-01-01 00:00:00','150 x 297 mm',150000,297000,'ENABLED'),
    (10,1, '2023-01-01 00:00:00','70 x 100 mm',70000,100000,'ENABLED'),
    (11,1, '2023-01-01 00:00:00','Euro (110 x 220 mm)',110000,220000,'ENABLED'),
    (12,1, '2023-01-01 00:00:00','C5 (162 x 229 mm)',162000,229000,'ENABLED'),
    (13,1, '2023-01-01 00:00:00','C4 (229 x 324 mm)',229000,324000,'ENABLED'),
    (14,1, '2023-01-01 00:00:00','90 x 90 mm',90000,90000,'ENABLED'),
    (15,1, '2023-01-01 00:00:00','100 x 150 mm',100000,150000,'ENABLED'),
    (16,1, '2023-01-01 00:00:00','100 x 200 mm',100000,200000,'ENABLED'),
    (17,1, '2023-01-01 00:00:00','230 x 410 mm (47х65 см/3)',410000,230000,'ENABLED'),
    (18,1, '2023-01-01 00:00:00','520 x 360 mm',520000,360000,'ENABLED'),
    (19,1, '2023-01-01 00:00:00','250 x 330 mm',250000,330000,'ENABLED'),
    (20,1, '2023-01-01 00:00:00','520 x 360 mm (72x104/4)',520000,360000,'ENABLED'),
    (23,1, '2023-01-01 00:00:00','SRA3 (320 x 450 mm)',320000,450000,'ENABLED'),
    (24,1, '2023-01-01 00:00:00','max (330 x 485 mm)',330000,485000,'ENABLED'),
    (25,1, '2023-01-01 00:00:00','350 x 500 mm',350000,500000,'ENABLED'),
    (26,1, '2023-01-01 00:00:00','250 х 450 mm (70х100 см/6)',250000,450000,'ENABLED'),
    (30,1, '2023-01-01 00:00:00','320 х 450 mm (64х90 см/4)',320000,450000,'ENABLED'),
    (31,1, '2023-01-01 00:00:00','325 х 470 mm (65х47 см/2)',325000,470000,'ENABLED'),
    (32,1, '2023-01-01 00:00:00','650 x 470 mm',650000,470000,'ENABLED'),
    (33,1, '2023-01-01 00:00:00','640 x 450 mm',650000,450000,'ENABLED'),
    (34,1, '2023-01-01 00:00:00','500 x 700 mm',500000,700000,'ENABLED'),
    (35,1, '2023-01-01 00:00:00','620 x 470 mm',620000,470000,'ENABLED'),
    (36,1, '2023-01-01 00:00:00','520 x 720 mm (72x104/2)',520000,720000,'ENABLED'),
    (37,1, '2023-01-01 00:00:00','330 x 700 mm (70x100 см/3)',330000,700000,'ENABLED'),
    (38,1, '2023-01-01 00:00:00','700 x 1000 mm',700000,1000000,'ENABLED'),
    (39,1, '2023-01-01 00:00:00','640 x 900 mm',640000,900000,'ENABLED'),
    (40,1, '2023-01-01 00:00:00','620 x 940 mm',620000,940000,'ENABLED'),
    (41,1, '2023-01-01 00:00:00','720 х 1040 мм',720000,1040000,'ENABLED'),
    (42,1, '2023-01-01 00:00:00','А1 (594 х 841 мм)',841000,594000,'ENABLED'),
    (43,1, '2023-01-01 00:00:00','А2 (420 х 594 мм)',594000,420000,'ENABLED'),
    (44,1, '2023-01-01 00:00:00','330 x 350 mm (70x100 см/6)',350000,330000,'ENABLED'),
    (45,1, '2023-01-01 00:00:00','297 х 210 мм',297000,210000,'ENABLED'),
    (46,1, '2023-01-01 00:00:00','420 х 297 мм',420000,297000,'ENABLED');

ALTER SEQUENCE catalog_print_formats_format_id_seq RESTART WITH 47;


CREATE TABLE catalog_boxes (
    box_id int4 NOT NULL GENERATED ALWAYS AS IDENTITY,
    tag_version int4 NOT NULL DEFAULT 1 CHECK(tag_version > 0),
    datetime_created timestamp NOT NULL DEFAULT now(),
    box_article varchar(64) NOT NULL,
    box_caption varchar(64) NOT NULL,
    box_length int4 NOT NULL CHECK(box_length > 0 AND box_length < 1000000000),-- mm * 1000
    box_width int4 NOT NULL CHECK(box_width > 0 AND box_width < 1000000000), -- mm * 1000
    box_depth int4 NOT NULL CHECK(box_depth > 0 AND box_depth < 1000000000), -- mm * 1000
    box_status item_status NOT NULL,
    CONSTRAINT catalog_boxes_pkey PRIMARY KEY (box_id),
    CONSTRAINT catalog_boxes_box_article UNIQUE (box_article)
);

-- ALTER SEQUENCE catalog_boxes_box_id_seq RESTART WITH 1;


CREATE TABLE catalog_paper_colors (
    color_id int4 NOT NULL GENERATED ALWAYS AS IDENTITY,
    tag_version int4 NOT NULL DEFAULT 1 CHECK(tag_version > 0),
    datetime_created timestamp NOT NULL DEFAULT now(),
    color_caption varchar(64) NOT NULL,
    color_status item_status NOT NULL,
    CONSTRAINT catalog_paper_colors_pkey PRIMARY KEY (color_id)
);

-- ALTER SEQUENCE catalog_paper_colors_color_id_seq RESTART WITH 1;


CREATE TABLE catalog_paper_factures (
    facture_id int4 NOT NULL GENERATED ALWAYS AS IDENTITY,
    tag_version int4 NOT NULL DEFAULT 1 CHECK(tag_version > 0),
    datetime_created timestamp NOT NULL DEFAULT now(),
    facture_caption varchar(64) NOT NULL,
    facture_status item_status NOT NULL,
    CONSTRAINT catalog_paper_factures_pkey PRIMARY KEY (facture_id)
);

-- ALTER SEQUENCE catalog_paper_factures_facture_id_seq RESTART WITH 1;


CREATE TABLE catalog_papers (
    paper_id int4 NOT NULL GENERATED ALWAYS AS IDENTITY,
    tag_version int4 NOT NULL DEFAULT 1 CHECK(tag_version > 0),
    datetime_created timestamp NOT NULL DEFAULT now(),
    paper_article varchar(64) NOT NULL,
    paper_caption varchar(64) NOT NULL,
    paper_length int4 NOT NULL CHECK(paper_length > 0 AND paper_length < 1000000000), -- mm * 1000
    paper_width int4 NOT NULL CHECK(paper_width > 0 AND paper_width < 1000000000), -- mm * 1000
    paper_density int4 NOT NULL CHECK(paper_density > 0 AND paper_density < 10000), -- g/m2
    color_id int4 NOT NULL CHECK(color_id > 0),
    facture_id int4 NOT NULL CHECK(facture_id > 0),
    paper_thickness int4 NOT NULL CHECK(paper_thickness > 0 AND paper_thickness < 1000000000), -- mm * 1000
    paper_sides catalog_paper_sides NOT NULL,
    paper_status item_status NOT NULL,
    CONSTRAINT catalog_papers_pkey PRIMARY KEY (paper_id),
    CONSTRAINT catalog_papers_paper_article UNIQUE (paper_article),
    CONSTRAINT catalog_paper_colors_fkey FOREIGN KEY (color_id)
        REFERENCES catalog_paper_colors(color_id),
    CONSTRAINT catalog_paper_factures_fkey FOREIGN KEY (facture_id)
        REFERENCES catalog_paper_factures(facture_id)
);

-- ALTER SEQUENCE catalog_papers_paper_id_seq RESTART WITH 1;


CREATE TABLE catalog_laminate_types (
    type_id int4 NOT NULL GENERATED ALWAYS AS IDENTITY,
    tag_version int4 NOT NULL DEFAULT 1 CHECK(tag_version > 0),
    datetime_created timestamp NOT NULL DEFAULT now(),
    type_caption varchar(64) NOT NULL,
    type_status item_status NOT NULL,
    CONSTRAINT catalog_laminate_types_pkey PRIMARY KEY (type_id)
);

-- ALTER SEQUENCE catalog_laminate_types_type_id_seq RESTART WITH 1;


CREATE TABLE catalog_laminates (
    laminate_id int4 NOT NULL GENERATED ALWAYS AS IDENTITY,
    tag_version int4 NOT NULL DEFAULT 1 CHECK(tag_version > 0),
    datetime_created timestamp NOT NULL DEFAULT now(),
    laminate_article varchar(64) NOT NULL,
    laminate_caption varchar(64) NOT NULL,
    type_id int4 NOT NULL CHECK(type_id > 0),
    laminate_length int4 NOT NULL CHECK(laminate_length > 0 AND laminate_length < 1000000000), -- mm * 1000
    laminate_weight int4 NOT NULL CHECK(laminate_weight > 0 AND laminate_weight < 10000), -- g/m2
    laminate_thickness int4 NOT NULL CHECK(laminate_thickness > 0 AND laminate_thickness < 1000000), -- mkm
    laminate_status item_status NOT NULL,
    CONSTRAINT catalog_laminates_pkey PRIMARY KEY (laminate_id),
    CONSTRAINT catalog_laminates_laminate_article UNIQUE (laminate_article),
    CONSTRAINT catalog_laminate_types_fkey FOREIGN KEY (type_id)
        REFERENCES catalog_laminate_types(type_id)
);

-- ALTER SEQUENCE catalog_laminates_laminate_id_seq RESTART WITH 1;


CREATE TABLE form_fields_test (
    field_id int4 NOT NULL GENERATED ALWAYS AS IDENTITY,
    field_value int4 NOT NULL,
    prev_field_id int4 NULL,
    next_field_id int4 NULL,
    order_field int8 NULL,
    CONSTRAINT form_fields_test_pkey PRIMARY KEY (field_id)
);

CREATE INDEX form_fields_test_order_field ON form_fields_test (order_field);

INSERT INTO form_fields_test (field_id, field_value, order_field, prev_field_id, next_field_id)
OVERRIDING SYSTEM VALUE
VALUES  (2, 100000, 17408, 4, 14),
        (37, 200000, 19456, 14, 1),
        (15, 200000, 4160, 17, 3),
        (3, 100000, 7168, 15, 4),
        (14, 200000, 18432, 2, 37),
        (29, 200000, 2306, 28, 12),
        (28, 200000, 2177, 19, 29),
        (19, 200000, 2112, 18, 28),
        (18, 200000, 2080, 36, 19),
        (36, 200000, 2064, 13, 18),
        (13, 200000, 2056, 35, 36),
        (35, 200000, 2052, 21, 13),
        (4, 100000, 8192, 3, 2),
        (12, 200000, 2564, 29, 17),
        (17, 200000, 4112, 12, 15),
        (24, 200000, null, null, null),
        (23, 200000, null, null, null),
        (32, 200000, null, null, null),
        (21, 200000, 2049, 33, 35),
        (30, 200000, null, null, null),
        (1, 100000, 20480, 37, 5),
        (16, 200000, 1024, null, 33),
        (25, 200000, null, null, null),
        (31, 200000, null, null, null),
        (22, 200000, null, null, null),
        (33, 200000, 1536, 16, 21),
        (20, 200000, null, null, null),
        (5, 100000, 21504, 1, null),
        (34, 200000, null, null, null),
        (39, 200000, null, null, null),
        (38, 200000, null, null, null);

ALTER SEQUENCE form_fields_test_field_id_seq RESTART WITH 40;
