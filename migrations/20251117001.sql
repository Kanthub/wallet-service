DO
$$
   BEGIN
        IF
            NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'uint256') THEN
            CREATE DOMAIN UINT256 AS NUMERIC CHECK (VALUE >= 0 AND VALUE < POWER(CAST(2 AS NUMERIC), CAST(256 AS NUMERIC)) AND SCALE(VALUE) = 0);
        ELSE
            ALTER DOMAIN UINT256 DROP CONSTRAINT uint256_check;
            ALTER DOMAIN UINT256 ADD CHECK (VALUE >= 0 AND VALUE < POWER(CAST(2 AS NUMERIC), CAST(256 AS NUMERIC)) AND SCALE(VALUE) = 0);
        END IF;
   END
$$;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp" cascade;

create table if not exists sys_log (
    guid          TEXT PRIMARY KEY DEFAULT replace(uuid_generate_v4()::text, '-', ''),
    action        VARCHAR(100) DEFAULT '', -- 路径 --
    desc          VARCHAR(100) DEFAULT '', -- 描述 --
    admin         VARCHAR(30)  DEFAULT '', -- 操作管理员 --
    ip            VARCHAR(30)  DEFAULT '', -- 操作管理员 IP --
    cate          SMALLINT DEFAULT 0,      -- 类型(0表示其他;1=>表示登陆;2=>表示财务操作) --
    status        SMALLINT DEFAULT -1,     -- 登陆状态(0=>成功;1=>失败) --
    asset         VARCHAR(255) DEFAULT '', -- 币种 --
    before        VARCHAR(255) DEFAULT '', -- 修改前 --
    after         VARCHAR(255) DEFAULT '', -- 修改后 --
    user_id       BIGINT DEFAULT 0,
    order_number  VARCHAR(64) DEFAULT '',
    op            SMALLINT DEFAULT -1,     -- 操作类型(0添加;1编辑)
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_syslog_cate ON sys_log (cate);
CREATE INDEX idx_syslog_status ON sys_log (status);
CREATE INDEX idx_syslog_order_number ON sys_log (order_number);


create table if not exists auth (
    guid         TEXT PRIMARY KEY DEFAULT replace(uuid_generate_v4()::text, '-', ''),
    auth_name    VARCHAR(255) DEFAULT '' COMMENT '权限名称',
    auth_url     VARCHAR(255) DEFAULT '' COMMENT '权限路径/接口地址',
    user_id      INT DEFAULT 0 COMMENT '所属用户/管理员ID',
    pid          INT DEFAULT 0 COMMENT '父级权限ID',
    sort         INT DEFAULT 0 COMMENT '排序',
    icon         VARCHAR(255) DEFAULT '' COMMENT '图标',
    is_show      INT DEFAULT 1 COMMENT '是否显示(1显示;0隐藏)',
    status       INT DEFAULT 1 COMMENT '状态(1启用;0禁用)',
    create_id    INT DEFAULT 0 COMMENT '创建人ID',
    update_id    INT DEFAULT 0 COMMENT '修改人ID',
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_auth_user_id   ON auth (user_id);
CREATE INDEX idx_auth_pid       ON auth (pid);
CREATE INDEX idx_auth_create_id ON auth (create_id);
CREATE INDEX idx_auth_update_id ON auth (update_id);


create table if not exists role (
    guid         TEXT PRIMARY KEY DEFAULT replace(uuid_generate_v4()::text, '-', ''),
    role_name    VARCHAR(100) DEFAULT '' COMMENT '角色名称',
    detail       VARCHAR(255) DEFAULT '' COMMENT '角色描述/说明',
    status       INT DEFAULT 1 COMMENT '状态(1启用;0禁用)',
    create_id    INT DEFAULT 0 COMMENT '创建人ID',
    update_id    INT DEFAULT 0 COMMENT '修改人ID',
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_role_role_name ON role (role_name);


create table if not exists role_auth (
    auth_id  INT NOT NULL,
    role_id  BIGINT NOT NULL,
    PRIMARY KEY (auth_id, role_id)
);
CREATE INDEX idx_role_auth_role_id ON role_auth (role_id);


create table if not exists admin (
    guid          TEXT PRIMARY KEY DEFAULT replace(uuid_generate_v4()::text, '-', ''),
    login_name    VARCHAR(32)  NOT NULL UNIQUE,   -- 登录名 --
    real_name     VARCHAR(32)  UNIQUE,            -- 真实姓名 --
    password      VARCHAR(100) NOT NULL,          -- 密码(加密后) --
    role_ids      VARCHAR(255) DEFAULT '',        -- 角色 ID 列表（字符串存 JSON/CSV） --
    phone         VARCHAR(11) UNIQUE,             -- 手机号 --
    email         VARCHAR(32),                    -- 邮箱 --
    salt          VARCHAR(255) DEFAULT '',        -- 密码盐 --
    last_login    INTEGER DEFAULT 0,              -- 最后登录时间戳 --
    last_ip       VARCHAR(255) DEFAULT '',        -- 最后登录 IP --
    status        INT DEFAULT 1,                  -- 状态(1启用;0禁用) --
    create_id     INT DEFAULT 0,                  -- 创建人 --
    update_id     INT DEFAULT 0,                  -- 修改人 --
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_admin_status ON admin (status);
CREATE INDEX idx_admin_create_id ON admin (create_id);
CREATE INDEX idx_admin_update_id ON admin (update_id);
CREATE INDEX idx_admin_last_login ON admin (last_login);

CREATE TABLE if not exists chain (
    guid               TEXT PRIMARY KEY DEFAULT replace(uuid_generate_v4()::text, '-', ''),
    chain_name         VARCHAR(70) NOT NULL,
    chain_mark         VARCHAR(70) NOT NULL,
    chain_logo         VARCHAR(200) NOT NULL,
    chain_active_logo  VARCHAR(200) NOT NULL,
    chain_model_type   VARCHAR(10) NOT NULL,  --utxo/account--
    created_at         TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at         TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_chain_guid ON chain (guid);
CREATE INDEX idx_chain_name_guid ON chain (chain_name);


CREATE TABLE if not exists token (
    guid                    TEXT PRIMARY KEY DEFAULT replace(uuid_generate_v4()::text, '-', ''),
    token_name              VARCHAR(70) DEFAULT '',
    token_mark              VARCHAR(70) DEFAULT '',
    token_logo              VARCHAR(100) DEFAULT '',
    token_active_logo       VARCHAR(100) DEFAULT '',
    token_decimal           VARCHAR(10) DEFAULT '18',
    token_symbol            VARCHAR(70) DEFAULT '',
    token_contract_address  VARCHAR(70) NOT NULL,
    token_chain_uuid        VARCHAR(255) DEFAULT '',
    is_hot                  VARCHAR(32) NOT NULL DEFAULT 'hot',
    created_at              TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at              TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_token_guid ON chain (guid);
CREATE INDEX idx_token_token_name ON chain (token_name);

CREATE TABLE if not exists chain_token (
    guid          TEXT PRIMARY KEY DEFAULT replace(uuid_generate_v4()::text, '-', ''),
    chain_uuid    VARCHAR(255) DEFAULT '',
    token_uuid    VARCHAR(255) NOT NULL,
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)
CREATE INDEX idx_chain_token_guid ON chain_token (guid);
CREATE INDEX idx_chain_token_chain_uuid ON chain_token (chain_uuid);

CREATE TABLE if not exists wallet (
    guid          TEXT PRIMARY KEY DEFAULT replace(uuid_generate_v4()::text, '-', ''),
    device_uuid   VARCHAR(255) NOT NULL,
    wallet_uuid   VARCHAR(255) NOT NULL,
    chain_uuid    VARCHAR(255) DEFAULT '',
    wallet_name   VARCHAR(70) DEFAULT 'roothash',
    asset_usdt    NUMERIC(20, 8) NOT NULL,
    asset_usd     NUMERIC(20, 8) NOT NULL,
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_wallet_guid ON wallet (guid);
CREATE INDEX idx_wallet_wallet_uuid ON wallet (wallet_uuid);

CREATE TABLE if not exists wallet_address (
    guid             TEXT PRIMARY KEY DEFAULT replace(uuid_generate_v4()::text, '-', ''),
    address_index    INTEGER CHECK (address_index > 0),
    address          VARCHAR(70) NOT NULL,
    wallet_uuid      VARCHAR(255) DEFAULT ''
    created_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_wallet_address_guid ON wallet_address (guid);
CREATE INDEX idx_wallet_address_wallet_uuid ON wallet_address (wallet_uuid);

CREATE TABLE if not exists wallet_asset (
    guid          TEXT PRIMARY KEY DEFAULT replace(uuid_generate_v4()::text, '-', ''),
    token_uuid    VARCHAR(255) DEFAULT '',
    chain_uuid    VARCHAR(255) DEFAULT '',
    balance       INTEGER CHECK (balance > 0),
    asset_usdt    NUMERIC(20, 8) NOT NULL,
    asset_usd     NUMERIC(20, 8) NOT NULL,
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_wallet_asset_guid ON wallet_asset (guid);
CREATE INDEX idx_wallet_asset_chain_uuid ON wallet_address (chain_uuid);

CREATE TABLE if not exists asset_amount_stat (
    guid          TEXT PRIMARY KEY DEFAULT replace(uuid_generate_v4()::text, '-', ''),
    asset_uuid    VARCHAR(255) DEFAULT '',
    time_date     VARCHAR(255) NOT NULL,
    amount        INTEGER CHECK (amount > 0),
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_asset_amount_stat_guid ON asset_amount_stat (guid);
CREATE INDEX idx_asset_amount_stat_asset_uuid ON asset_amount_stat (asset_uuid);


CREATE TABLE if not exists address_asset (
    guid          TEXT PRIMARY KEY DEFAULT replace(uuid_generate_v4()::text, '-', ''),
    token_uuid    VARCHAR(255) DEFAULT '',
    wallet_uuid   VARCHAR(255) DEFAULT '',
    address_uuid  VARCHAR(255) DEFAULT '',
    asset_usdt    NUMERIC(20, 8) NOT NULL,
    asset_usd     NUMERIC(20, 8) NOT NULL,
    balance       INTEGER CHECK (balance > 0),
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_address_asset_guid ON address_asset (guid);
CREATE INDEX idx_address_asset_address_uuid ON address_asset (address_uuid);


CREATE TABLE if not exists wallet_tx_record (
    guid          TEXT PRIMARY KEY DEFAULT replace(uuid_generate_v4()::text, '-', ''),
    tx_time       VARCHAR(500) NOT NULL,
    chain_uuid    VARCHAR(255) DEFAULT '',
    token_uuid    VARCHAR(255) DEFAULT '',
    from_address  VARCHAR(70) NOT NULL,
    to_address    VARCHAR(70) NOT NULL,
    amount        INTEGER CHECK (amount > 0),
    memo          VARCHAR(500) NOT NULL,
    hash          VARCHAR(500) NOT NULL,
    block_height  VARCHAR(500) NOT NULL,
    explorer_url  VARCHAR(500) NOT NULL,
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_wallet_tx_record_guid ON wallet_tx_record (guid);
CREATE INDEX idx_wallet_tx_record_from_address ON wallet_tx_record (from_address);
CREATE INDEX idx_wallet_tx_record_to_address ON wallet_tx_record (to_address);


CREATE TABLE if not exists wallet_address_note (
    guid         TEXT PRIMARY KEY DEFAULT replace(uuid_generate_v4()::text, '-', ''),
    device_uuid  VARCHAR(255) NOT NULL,
    chain_uuid   VARCHAR(255) DEFAULT '',
    memo         VARCHAR(255) NOT NULL,
    address      VARCHAR(255) NOT NULL,
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_wallet_address_note_guid ON wallet_address_note (guid);
CREATE INDEX idx_wallet_address_note_device_uuid ON wallet_address_note (device_uuid);


CREATE TABLE if not exists fiat_currency_rate (
    guid         TEXT PRIMARY KEY DEFAULT replace(uuid_generate_v4()::text, '-', ''),
    key_name     VARCHAR(255) NOT NULL,
    value_data   VARCHAR(255) NOT NULL,
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)
CREATE INDEX idx_fiat_currency_rate_guid ON fiat_currency_rate (guid);
CREATE INDEX idx_fiat_currency_rate_key_name ON fiat_currency_rate (key_name);


CREATE TABLE if not exists market_price (
    guid         TEXT PRIMARY KEY DEFAULT replace(uuid_generate_v4()::text, '-', ''),
    chain_uuid   VARCHAR(255) DEFAULT '',
    token_uuid   VARCHAR(255) DEFAULT '',
    usdt_price   NUMERIC(20, 8) NOT NULL,
    usd_price    NUMERIC(20, 8) NOT NULL,
    market_cap   INTEGER CHECK (market_cap > 0),
    liquidity    INTEGER CHECK (liquidity > 0),
    24h_volume   INTEGER CHECK (liquidity > 0),
    price_change VARCHAR(255) NOT NULL,
    ranking      VARCHAR(255) NOT NULL,
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)
CREATE INDEX idx_market_price_guid ON market_price (guid);
CREATE INDEX idx_market_price_token_uuid ON market_price (key_name);


CREATE TABLE IF NOT EXISTS kline (
    guid          TEXT PRIMARY KEY DEFAULT replace(uuid_generate_v4()::text, '-', ''),
    token_id      VARCHAR NOT NULL,
    time_interval VARCHAR NOT NULL,          -- K线周期，如 1m, 5m, 1h, 1d --
    open_time     TIMESTAMP NOT NULL,        -- 开始时间 --
    open_price    NUMERIC(20, 8) NOT NULL,   -- 开盘价 --
    high_price    NUMERIC(20, 8) NOT NULL,   -- 最高价 --
    low_price     NUMERIC(20, 8) NOT NULL,   -- 最低价 --
    close_price   NUMERIC(20, 8) NOT NULL,   -- 收盘价 --
    volume        UINT256 NOT NULL,          -- 成交量（币数量） --
    quote_volume  UINT256 DEFAULT 0,         -- 成交额（计价货币） --
    trade_count   UINT256 DEFAULT 0,         -- 成交笔数 --
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_kline_interval_time ON kline(token_id, time_interval, open_time);
CREATE INDEX IF NOT EXISTS idx_kline_time ON kline(symbol_id, open_time DESC);

CREATE TABLE IF NOT EXISTS newsletter_cat (
    guid          VARCHAR PRIMARY KEY,
    cat_name      VARCHAR NOT NULL,
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_newsletter_cat_guid ON newsletter_cat(guid);

CREATE TABLE IF NOT EXISTS newsletter (
    guid          VARCHAR PRIMARY KEY,
    cat_uuid      VARCHAR(255) NOT NULL,
    title         VARCHAR(255) NOT NULL,
    image         VARCHAR(700) NOT NULL,
    detail        TEXT DEFAULT '',
    link_url      VARCHAR(255) NOT NULL,
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_newsletter_guid ON newsletter(guid);
CREATE UNIQUE INDEX IF NOT EXISTS idx_newsletter_cat_uuid ON newsletter(cat_uuid);
