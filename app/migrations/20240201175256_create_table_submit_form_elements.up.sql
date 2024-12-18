-- --------------------------------------------------------------------------------------------------

CREATE TABLE printshop_controls.submit_form_elements (
    element_id int8 NOT NULL GENERATED BY DEFAULT AS IDENTITY CONSTRAINT pk_submit_form_elements PRIMARY KEY,
    tag_version int4 NOT NULL DEFAULT 1 CHECK(tag_version > 0),
    form_id uuid NOT NULL CONSTRAINT fk_submit_form_elements_form_id REFERENCES printshop_controls.submit_forms(form_id) ON DELETE CASCADE,
    param_name character varying(32) NOT NULL,
    element_caption character varying(64) NOT NULL,
    template_id int8 NOT NULL CONSTRAINT fk_submit_form_elements_template_id REFERENCES printshop_controls.element_templates (template_id),
    template_version int4 NOT NULL CHECK(template_version > 0),
    element_required bool NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp with time zone NOT NULL DEFAULT NOW(),
    deleted_at timestamp with time zone NULL,
    prev_field_id int8 NULL CHECK(prev_field_id IS NULL OR prev_field_id > 0),
    next_field_id int8 NULL CHECK(next_field_id IS NULL OR next_field_id > 0),
    order_index int8 NULL CHECK(order_index IS NULL OR order_index > 0)
);

CREATE UNIQUE INDEX uk_submit_form_elements_form_id_param_name ON printshop_controls.submit_form_elements (form_id, param_name) WHERE deleted_at IS NULL;
CREATE INDEX ix_submit_form_elements_form_id_order_index ON printshop_controls.submit_form_elements (form_id, order_index) WHERE deleted_at IS NULL;

-- --------------------------------------------------------------------------------------------------

INSERT INTO printshop_controls.submit_form_elements (element_id, tag_version, form_id, param_name, element_caption, template_id, template_version, element_required, created_at, updated_at, deleted_at, prev_field_id, next_field_id, order_index)
VALUES (1, 1, 'aa22434f-f09d-4b20-89c0-0785948cdc04', 'Product', 'Поля листовой продукции', 1, 1, true, '2023-07-15 13:49:58.567032', '2023-07-15 13:49:58.567032', NULL, NULL, NULL, 1048576);

ALTER SEQUENCE printshop_controls.submit_form_elements_element_id_seq RESTART WITH 2;