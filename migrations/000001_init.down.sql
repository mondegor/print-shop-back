-- --------------------------------------------------------------------------------------------------

DROP TABLE ps_provider_accounts.companies_pages;

DROP TYPE ps_provider_accounts.resource_status;

DROP SCHEMA ps_provider_accounts;

-- --------------------------------------------------------------------------------------------------

DROP TABLE ps_controls.form_elements;
DROP TABLE ps_controls.forms;
DROP TABLE ps_controls.element_templates;

DROP TYPE ps_controls.item_status;
DROP TYPE ps_controls.element_detailing;
DROP TYPE ps_controls.element_type;

DROP SCHEMA ps_controls;

-- --------------------------------------------------------------------------------------------------

DROP TABLE ps_dictionaries.laminate_types;
DROP TABLE ps_dictionaries.paper_colors;
DROP TABLE ps_dictionaries.paper_factures;
DROP TABLE ps_dictionaries.print_formats;

DROP TYPE ps_dictionaries.item_status;

DROP SCHEMA ps_dictionaries;

-- --------------------------------------------------------------------------------------------------

DROP TABLE ps_catalog.boxes;
DROP TABLE ps_catalog.laminates;
DROP TABLE ps_catalog.papers;

DROP TYPE ps_catalog.item_status;
DROP TYPE ps_catalog.paper_sides;

DROP SCHEMA ps_catalog;