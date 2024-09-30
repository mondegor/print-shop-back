-- --------------------------------------------------------------------------------------------------

CREATE TABLE printshop_global.settings (
    setting_id int4 NOT NULL CONSTRAINT pk_settings PRIMARY KEY,
    setting_name character varying(64) NOT NULL,
    setting_type int2 NOT NULL, -- 1=STRING, 2=STRING_LIST, 3=INTEGER, 4=INTEGER_LIST, 5=BOOLEAN
    setting_value character varying(65536) NOT NULL,
    setting_description character varying(1024) NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp with time zone NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX uk_settings_setting_name ON printshop_global.settings (setting_name);
CREATE INDEX ix_settings_updated_at ON printshop_global.settings (updated_at);

-- --------------------------------------------------------------------------------------------------

INSERT INTO printshop_global.settings (setting_id, setting_name, setting_type, setting_value, setting_description, created_at, updated_at)
VALUES
    (1, 'providers.registration.enabled', 5/*BOOLEAN*/, 'true', 'Возможность регистрации провайдеров', '2023-01-01 12:15:59.981966', '2023-01-01 12:15:59.981966');