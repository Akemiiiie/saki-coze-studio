IF NOT EXISTS (SELECT * FROM sys.schemas WHERE name = 'dbo')
BEGIN
    EXEC('CREATE SCHEMA dbo')
END

IF NOT EXISTS (SELECT * FROM sysobjects WHERE name='plugin_draft' AND xtype='U')
CREATE TABLE dbo.plugin_draft (
    id bigint NOT NULL DEFAULT 0,
    space_id bigint NOT NULL DEFAULT 0,
    developer_id bigint NOT NULL DEFAULT 0,
    app_id bigint NOT NULL DEFAULT 0,
    icon_uri nvarchar(512) NOT NULL DEFAULT '',
    server_url nvarchar(512) NOT NULL DEFAULT '',
    plugin_type tinyint NOT NULL DEFAULT 0,
    created_at bigint NOT NULL DEFAULT 0,
    updated_at bigint NOT NULL DEFAULT 0,
    deleted_at datetime NULL,
    manifest nvarchar(max) NULL,
    openapi_doc nvarchar(max) NULL,
    CONSTRAINT PK_plugin_draft PRIMARY KEY (id)
);

CREATE NONCLUSTERED INDEX IX_plugin_draft_app_id ON dbo.plugin_draft (app_id, id);
CREATE NONCLUSTERED INDEX IX_plugin_draft_space_app_created_at ON dbo.plugin_draft (space_id, app_id, created_at);
CREATE NONCLUSTERED INDEX IX_plugin_draft_space_app_updated_at ON dbo.plugin_draft (space_id, app_id, updated_at);

IF NOT EXISTS (SELECT * FROM sysobjects WHERE name='plugin_oauth_auth' AND xtype='U')
CREATE TABLE dbo.plugin_oauth_auth (
    id bigint NOT NULL DEFAULT 0,
    user_id nvarchar(255) NOT NULL DEFAULT '',
    plugin_id bigint NOT NULL DEFAULT 0,
    is_draft bit NOT NULL DEFAULT 0,
    oauth_config nvarchar(max) NULL,
    access_token nvarchar(max) NULL,
    refresh_token nvarchar(max) NULL,
    token_expired_at bigint NULL,
    next_token_refresh_at bigint NULL,
    last_active_at bigint NULL,
    created_at bigint NOT NULL DEFAULT 0,
    updated_at bigint NOT NULL DEFAULT 0,
    CONSTRAINT PK_plugin_oauth_auth PRIMARY KEY (id)
);

CREATE NONCLUSTERED INDEX IX_plugin_oauth_auth_last_active_at ON dbo.plugin_oauth_auth (last_active_at);
CREATE NONCLUSTERED INDEX IX_plugin_oauth_auth_token_expired_at ON dbo.plugin_oauth_auth (token_expired_at);
CREATE NONCLUSTERED INDEX IX_plugin_oauth_auth_next_token_refresh_at ON dbo.plugin_oauth_auth (next_token_refresh_at);
CREATE UNIQUE NONCLUSTERED INDEX IX_plugin_oauth_auth_user_plugin_is_draft ON dbo.plugin_oauth_auth (user_id, plugin_id, is_draft);
