-- --------------------------------------------------------------------------------------------------

DROP TABLE IF EXISTS printshop_provider_accounts.companies_pages;

DROP SCHEMA printshop_provider_accounts;

-- --------------------------------------------------------------------------------------------------

DROP TABLE IF EXISTS printshop_controls.submit_form_elements;
DROP TABLE IF EXISTS printshop_controls.submit_forms;
DROP TABLE IF EXISTS printshop_controls.element_templates;

DROP SCHEMA printshop_controls;

-- --------------------------------------------------------------------------------------------------

DROP TABLE IF EXISTS printshop_dictionaries.laminate_types;
DROP TABLE IF EXISTS printshop_dictionaries.paper_colors;
DROP TABLE IF EXISTS printshop_dictionaries.paper_factures;
DROP TABLE IF EXISTS printshop_dictionaries.print_formats;

DROP SCHEMA printshop_dictionaries;

-- --------------------------------------------------------------------------------------------------

DROP TABLE IF EXISTS printshop_catalog.boxes;
DROP TABLE IF EXISTS printshop_catalog.laminates;
DROP TABLE IF EXISTS printshop_catalog.papers;

DROP SCHEMA printshop_catalog;