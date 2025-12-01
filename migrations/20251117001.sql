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
    created        INTEGER CHECK (created > 0),
    updated        INTEGER CHECK (updated > 0)
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
    created      INTEGER CHECK (created > 0),
    updated      INTEGER CHECK (updated > 0)
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
    created      INTEGER CHECK (created > 0),
    updated      INTEGER CHECK (updated > 0)
);
CREATE INDEX idx_role_role_name ON role (role_name);


create table if not exists role_auth (
    auth_id INT NOT NULL,
    role_id BIGINT NOT NULL,
    PRIMARY KEY (auth_id, role_id)
);
CREATE INDEX idx_role_auth_role_id ON role_auth (role_id);


create table if not exists admin (
    guid         TEXT PRIMARY KEY DEFAULT replace(uuid_generate_v4()::text, '-', ''),
    login_name    VARCHAR(32)  NOT NULL UNIQUE,   -- 登录名 --
    real_name     VARCHAR(32)  UNIQUE,            -- 真实姓名 --
    password      VARCHAR(100) NOT NULL,          -- 密码(加密后) --
    role_ids      VARCHAR(255) DEFAULT '',        -- 角色 ID 列表（字符串存 JSON/CSV） --
    phone         VARCHAR(11) UNIQUE,             -- 手机号 --
    email         VARCHAR(32),                    -- 邮箱 --
    salt          VARCHAR(255) DEFAULT '',        -- 密码盐 --
    last_login    BIGINT DEFAULT 0,               -- 最后登录时间戳 --
    last_ip       VARCHAR(255) DEFAULT '',        -- 最后登录 IP --
    status        INT DEFAULT 1,                  -- 状态(1启用;0禁用) --
    create_id     INT DEFAULT 0,                  -- 创建人 --
    update_id     INT DEFAULT 0,                  -- 修改人 --
    created      INTEGER CHECK (created > 0),
    updated      INTEGER CHECK (updated > 0)
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
    chain_model_type   VARCHAR(10) NOT NULL,
    created            INTEGER CHECK (created > 0),
    updated            INTEGER CHECK (updated > 0)
);

CREATE TABLE if not exists asset (
    guid                 TEXT PRIMARY KEY DEFAULT replace(uuid_generate_v4()::text, '-', ''),
    asset_name           VARCHAR(70) NOT NULL,
    asset_mark           VARCHAR(70) NOT NULL,
    asset_logo           VARCHAR(100),
    asset_active_logo    VARCHAR(100),
    asset_unit           VARCHAR(10) NOT NULL,
    asset_symbol         VARCHAR(70) NOT NULL,
    asset_contract_addr  VARCHAR(70) NOT NULL,
    asset_chain_uuid     VARCHAR(255) DEFAULT '',
    is_hot               VARCHAR(32) NOT NULL,
    created              INTEGER CHECK (created > 0),
    updated              INTEGER CHECK (updated > 0)
);

CREATE TABLE if not exists wallet (
    guid        TEXT PRIMARY KEY DEFAULT replace(uuid_generate_v4()::text, '-', ''),
    device_id   VARCHAR(70) NOT NULL,
    wallet_uuid VARCHAR(70) NOT NULL,
    wallet_name VARCHAR(70) NOT NULL,
    asset_usd   INTEGER CHECK (asset_usd > 0),
    asset_cny   INTEGER CHECK (asset_usd > 0),
    chain_uuid  VARCHAR(255) DEFAULT '',
    created     INTEGER CHECK (created > 0),
    updated     INTEGER CHECK (updated > 0)
);

CREATE TABLE if not exists wallet_address (
    guid             TEXT PRIMARY KEY DEFAULT replace(uuid_generate_v4()::text, '-', ''),
    address_index    VARCHAR(10) NOT NULL,
    address          VARCHAR(70) NOT NULL,
    wallet_id        VARCHAR(255) DEFAULT ''
    created          INTEGER CHECK (created > 0),
    updated          INTEGER CHECK (updated > 0)
);

CREATE TABLE if not exists wallet_asset (
    guid          TEXT PRIMARY KEY DEFAULT replace(uuid_generate_v4()::text, '-', ''),
    contract_addr VARCHAR(70) NOT NULL,
    asset_usd     INTEGER CHECK (asset_usd > 0),
    asset_cny     INTEGER CHECK (asset_cny > 0),
    balance       NUMERIC(65,30) NOT NULL,
    asset_uuid    VARCHAR(255) DEFAULT '',
    chain_uuid    VARCHAR(255) DEFAULT '',
    created       INTEGER CHECK (created > 0),
    updated       INTEGER CHECK (updated > 0)
);

CREATE TABLE if not exists address_amount_stat (
    guid          TEXT PRIMARY KEY DEFAULT replace(uuid_generate_v4()::text, '-', ''),
    amount        NUMERIC(65,30) NOT NULL,
    timedate      VARCHAR(70) NOT NULL,
    asset_uuid    VARCHAR(255) DEFAULT '',
    created       INTEGER CHECK (created > 0),
    updated       INTEGER CHECK (updated > 0)
);

CREATE TABLE if not exists address_asset (
    guid          TEXT PRIMARY KEY DEFAULT replace(uuid_generate_v4()::text, '-', ''),
    asset_usd     INTEGER CHECK (asset_usd > 0),
    asset_cny     INTEGER CHECK (asset_cny > 0),
    balance       INTEGER CHECK (balance > 0),
    address_uuid  VARCHAR(255) DEFAULT '',
    asset_uuid    VARCHAR(255) DEFAULT '',
    wallet_uuid   VARCHAR(255) DEFAULT '',
    created       INTEGER CHECK (created > 0),
    updated       INTEGER CHECK (updated > 0)
);

CREATE TABLE if not exists wallet_tx_record (
    guid          TEXT PRIMARY KEY DEFAULT replace(uuid_generate_v4()::text, '-', ''),
    from_addr     VARCHAR(70) NOT NULL,
    to_addr       VARCHAR(70) NOT NULL,
    amount        NUMERIC(65,30) NOT NULL,
    memo          VARCHAR(70) NOT NULL,
    hash          VARCHAR(70) NOT NULL,
    block_height  VARCHAR(70) NOT NULL,
    tx_time       VARCHAR(70) NOT NULL,
    asset_uuid    VARCHAR(255) DEFAULT '',
    chain_uuid    VARCHAR(255) DEFAULT '',
    created       INTEGER CHECK (created > 0),
    updated       INTEGER CHECK (updated > 0)
);

CREATE TABLE if not exists wallet_address_note (
    guid         TEXT PRIMARY KEY DEFAULT replace(uuid_generate_v4()::text, '-', ''),
    device_id    VARCHAR(255) NOT NULL,
    memo         VARCHAR(255) NOT NULL,
    address      VARCHAR(255) NOT NULL,
    asset_uuid   VARCHAR(255) DEFAULT '',
    chain_uuid   VARCHAR(255) DEFAULT '',
    created      INTEGER CHECK (created > 0),
    updated      INTEGER CHECK (updated > 0)
);