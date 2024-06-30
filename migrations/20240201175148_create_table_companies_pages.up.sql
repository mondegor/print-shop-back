-- --------------------------------------------------------------------------------------------------

CREATE TABLE printshop_providers.companies_pages (
    account_id uuid NOT NULL CONSTRAINT pk_companies_pages PRIMARY KEY,
    rewrite_name character varying(32) NULL CONSTRAINT uk_companies_pages_rewrite_name UNIQUE,
    page_title character varying(128) NOT NULL,
    logo_meta jsonb DEFAULT NULL,
    site_url character varying(512) NOT NULL,
    page_status int2 NOT NULL, -- 1=DRAFT, 2=HIDDEN', 3=PUBLISHED, 4=PUBLISHED_SHARED
    created_at timestamp with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp with time zone NULL
);