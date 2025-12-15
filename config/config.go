package config

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Migrations                string           `yaml:"migrations"`
	MasterDB                  DBConfig         `yaml:"master_db"`
	SlaveDB                   DBConfig         `yaml:"slave_db"`
	SlaveDbEnable             bool             `yaml:"slave_db_enable"`
	ApiCacheEnable            bool             `yaml:"api_cache_enable"`
	CacheConfig               CacheConfig      `yaml:"cache_config"`
	RpcServer                 ServerConfig     `yaml:"rpc_server"`
	MetricsServer             ServerConfig     `yaml:"metrics_server"`
	HttpServer                ServerConfig     `yaml:"http_server"`
	WebsocketServer           ServerConfig     `yaml:"websocket_server"`
	EmailConfig               EmailConfig      `yaml:"email_config"`
	SMSConfig                 SMSConfig        `yaml:"sms_config"`
	MinioConfig               MinioConfig      `yaml:"minio_config"`
	KodoConfig                KodoConfig       `yaml:"kodo_config"`
	S3Config                  S3Config         `yaml:"s3_config"`
	CORSAllowedOrigins        string           `yaml:"cors_allowed_origins"`
	JWTSecret                 string           `yaml:"jwt_secret"`
	Domain                    string           `yaml:"domain"`
	PrivateKey                string           `yaml:"private_key"`
	NumConfirmations          uint64           `yaml:"num_confirmations"`
	SafeAbortNonceTooLowCount uint64           `yaml:"safe_abort_nonce_too_low_count"`
	CallerAddress             string           `yaml:"caller_address"`
	RedisConfig               RedisConfig      `yaml:"redis_config"`
	AggregatorConfig          AggregatorConfig `yaml:"aggregator_config"`

	RpcConfig RpcConfig `yaml:"rpc_config"`
	Chains    []string  `yaml:"chains"`
}

type RpcConfig struct {
	EthRpc      string `yaml:"eth_rpc"`
	ArbitrumRpc string `yaml:"arbitrum_rpc"`
	PolygonRpc  string `yaml:"polygon_rpc"`
	BaseRpc     string `yaml:"base_rpc"`
	RootHashRpc string `yaml:"roothash_rpc"`
	OpRpc       string `yaml:"op_rpc"`
	DataApiUrl  string `yaml:"data_api_url"`
	DataApiKey  string `yaml:"data_api_key"`
}

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Name     string `yaml:"name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type EmailConfig struct {
	SMTPHost     string `yaml:"smtp_host"`     // SMTP服务器地址
	SMTPPort     int    `yaml:"smtp_port"`     // SMTP端口
	SMTPUser     string `yaml:"smtp_user"`     // SMTP用户名
	SMTPPassword string `yaml:"smtp_password"` // SMTP密码
	FromName     string `yaml:"from_name"`     // 发件人名称
	FromEmail    string `yaml:"from_email"`    // 发件人邮箱
	UseSSL       bool   `yaml:"use_ssl"`       // 是否使用SSL/TLS
}

type SMSConfig struct {
	AccessKeyId     string `yaml:"access_key_id"`     // 阿里云AccessKeyId
	AccessKeySecret string `yaml:"access_key_secret"` // 阿里云AccessKeySecret
	SignName        string `yaml:"sign_name"`         // 短信签名
	TemplateCode    string `yaml:"template_code"`     // 短信模板代码
	Endpoint        string `yaml:"endpoint"`          // 短信服务端点
}

type MinioConfig struct {
	Endpoint        string `yaml:"endpoint"`
	AccessKeyID     string `yaml:"access_key_id"`
	SecretAccessKey string `yaml:"secret_access_key"`
	UseSSL          bool   `yaml:"use_ssl"`
	BucketName      string `yaml:"bucket_name"`
	BaseURL         string `yaml:"base_url"`
}

type CacheConfig struct {
	ListSize         int           `yaml:"list_size"`
	DetailSize       int           `yaml:"detail_size"`
	ListExpireTime   time.Duration `yaml:"list_expire_time"`
	DetailExpireTime time.Duration `yaml:"detail_expire_time"`
}

type ServerConfig struct {
	Scheme string `yaml:"scheme"` // http / https
	Host   string `yaml:"host"`
	Port   int    `yaml:"port"`
	Path   string `yaml:"path"` // 可选，比如 /v2/xxx
}

type KodoConfig struct {
	AccessKey     string `yaml:"access_key"`
	SecretKey     string `yaml:"secret_key"`
	Bucket        string `yaml:"bucket"`
	Domain        string `yaml:"domain"`
	Zone          string `yaml:"zone"`
	UseHTTPS      bool   `yaml:"use_https"`
	UseCdnDomains bool   `yaml:"use_cdn_domains"`
}

type S3Config struct {
	AccessKey    string `yaml:"access_key"`
	SecretKey    string `yaml:"secret_key"`
	Bucket       string `yaml:"bucket"`
	Region       string `yaml:"region"`
	Endpoint     string `yaml:"endpoint"`
	CDNDomain    string `yaml:"cdn_domain"`
	UsePathStyle bool   `yaml:"use_path_style"`
}

type RedisConfig struct {
	Addr     string `yaml:"addr"`     // Redis server address (host:port)
	Password string `yaml:"password"` // Redis password (optional)
	DB       int    `yaml:"db"`       // Redis database number
}

type AggregatorConfig struct {
	WalletAccountAddr          string            `yaml:"wallet_account_addr"`           // wallet-chain-account gRPC address
	WalletAccountConsumerToken string            `yaml:"wallet_account_consumer_token"` // default consumer token for wallet-chain-account
	ChainConsumerTokens        map[string]string `yaml:"chain_consumer_tokens"`         // optional per-chain tokens keyed by chain_id
	ZeroXAPIURL                string            `yaml:"zerox_api_url"`                 // 0x Protocol API URL
	ZeroXAPIKey                string            `yaml:"zerox_api_key"`                 // 0x Protocol API Key
	OneInchAPIURL              string            `yaml:"oneinch_api_url"`               // 1inch API URL
	OneInchAPIKey              string            `yaml:"oneinch_api_key"`               // 1inch API Key
	JupiterAPIURL              string            `yaml:"jupiter_api_url"`               // Jupiter API URL
	LiFiAPIURL                 string            `yaml:"lifi_api_url"`                  // LiFi API URL
	LiFiAPIKey                 string            `yaml:"lifi_api_key"`                  // LiFi API Key
	EnableProviders            map[string]bool   `yaml:"enable_providers"`              // Enable/disable specific providers
}

func New(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", path, err)
	}

	cfg := &Config{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml: %w", err)
	}

	if reflect.DeepEqual(cfg, &Config{}) {
		return nil, fmt.Errorf("config file %s is empty or invalid", path)
	}

	return cfg, nil
}

func (c ServerConfig) RPCURL() string {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "http"
	}

	host := c.Host
	if c.Port > 0 {
		host = fmt.Sprintf("%s:%d", c.Host, c.Port)
	}

	if c.Path != "" {
		return fmt.Sprintf("%s://%s/%s", scheme, host, strings.TrimPrefix(c.Path, "/"))
	}

	return fmt.Sprintf("%s://%s", scheme, host)
}

func (c RpcConfig) RPC(chain string) (string, error) {
	switch chain {
	case "eth":
		return c.EthRpc, nil
	case "arbitrum":
		return c.ArbitrumRpc, nil
	case "polygon":
		return c.PolygonRpc, nil
	case "base":
		return c.BaseRpc, nil
	case "op":
		return c.OpRpc, nil
	case "roothash":
		return c.RootHashRpc, nil
	default:
		return "", fmt.Errorf("unsupported chain %s", chain)
	}
}
