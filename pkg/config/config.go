package config

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/hassiimykyta/life-rpg/pkg/helpers"
	"github.com/joho/godotenv"
)

type loadCaps struct {
	useDB      bool
	useCORS    bool
	useJWT     bool
	useStorage bool
	useCache   bool
	useSMTP    bool
}

type Option func(*loadCaps)

func WithDB() Option      { return func(c *loadCaps) { c.useDB = true } }
func WithCORS() Option    { return func(c *loadCaps) { c.useCORS = true } }
func WithJWT() Option     { return func(c *loadCaps) { c.useJWT = true } }
func WithStorage() Option { return func(c *loadCaps) { c.useStorage = true } }
func WithCache() Option   { return func(c *loadCaps) { c.useCache = true } }
func WithSMTP() Option    { return func(c *loadCaps) { c.useSMTP = true } }

type AppConfig struct {
	Env             string
	Host            string
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
}

type DBConfig struct {
	Driver      string
	DSN         string
	MaxOpen     int
	MaxIdle     int
	MaxIdleTime time.Duration
}

type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	ExposedHeaders   []string
	AllowCredentials bool
	MaxAge           int
}

type JWTConfig struct {
	Secret     string
	Issuer     string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

type StorageConfig struct {
	Driver        string
	LocalDir      string
	PublicBaseURL string
	PresignTTL    time.Duration
	S3Bucket      string
	S3Region      string
	S3Endpoint    string
	S3AccessKey   string
	S3SecretKey   string
	S3UsePath     bool
}

type CacheConfig struct {
	Addr         string
	Password     string
	DB           int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	PoolSize     int
	MinIdleConns int
	TLSEnabled   bool
}

type SMTPConfig struct {
	Username string
	Password string
	Host     string
	Port     int
}

type Config struct {
	App     AppConfig
	DB      *DBConfig
	CORS    *CORSConfig
	JWT     *JWTConfig
	Storage *StorageConfig
	Cache   *CacheConfig
	SMTP    *SMTPConfig
}

func (c *Config) Validate(cap loadCaps) error {
	if cap.useCORS && c.CORS != nil && c.App.Env == "prod" {
		if len(c.CORS.AllowedOrigins) == 1 && c.CORS.AllowedOrigins[0] == "*" {
			if c.CORS.AllowCredentials {
				return fmt.Errorf("in prod, CORS_ALLOW_CREDENTIALS=true cannot be used with CORS_ALLOWED_ORIGINS=*")
			}
			log.Println("⚠️  WARNING: prod + CORS_ALLOWED_ORIGINS=* (no credentials). Consider whitelisting domains.")
		}
	}

	if cap.useDB {
		if c.DB == nil {
			return errors.New("DB config required but missing (enable WithDB and provide envs)")
		}
		if c.DB.DSN == "" {
			return errors.New("DB_DSN is required when DB is enabled")
		}
		switch c.DB.Driver {
		case "postgres", "pgx", "mysql":
		default:
			return fmt.Errorf("unsupported DB_DRIVER %q", c.DB.Driver)
		}
	}

	if cap.useJWT {
		if c.JWT == nil {
			return errors.New("JWT config required but missing (enable WithJWT and provide envs)")
		}
		if c.JWT.Secret == "" {
			return errors.New("JWT_SECRET is required when JWT is enabled")
		}
	}

	return nil
}
func (c *Config) RedactedString() string {
	return fmt.Sprintf(
		"env=%s host=%s port=%s db.driver=%s db.maxOpen=%d db.maxIdle=%d",
		c.App.Env, c.App.Host, c.App.Port, c.DB.Driver, c.DB.MaxOpen, c.DB.MaxIdle,
	)
}

func Load(opts ...Option) (*Config, error) {
	_ = godotenv.Load()

	var caps loadCaps
	for _, o := range opts {
		o(&caps)
	}

	cfg := &Config{
		App: AppConfig{
			Env:             helpers.GetEnv("APP_ENV", "dev"),
			Host:            helpers.GetEnv("APP_HOST", "0.0.0.0"),
			Port:            helpers.GetEnv("APP_PORT", "8080"),
			ReadTimeout:     helpers.MustDur(helpers.GetEnv("READ_TIMEOUT", "10s"), 10*time.Second),
			WriteTimeout:    helpers.MustDur(helpers.GetEnv("WRITE_TIMEOUT", "10s"), 10*time.Second),
			IdleTimeout:     helpers.MustDur(helpers.GetEnv("IDLE_TIMEOUT", "60s"), 60*time.Second),
			ShutdownTimeout: helpers.MustDur(helpers.GetEnv("SHUTDOWN_TIMEOUT", "10s"), 10*time.Second),
		},
	}

	if caps.useDB {
		cfg.DB = &DBConfig{
			Driver:      helpers.GetEnv("DB_DRIVER", "postgres"),
			DSN:         helpers.GetEnv("DB_DSN", ""),
			MaxOpen:     helpers.MustInt(helpers.GetEnv("DB_MAX_OPEN", "20"), 20),
			MaxIdle:     helpers.MustInt(helpers.GetEnv("DB_MAX_IDLE", "10"), 10),
			MaxIdleTime: helpers.MustDur(helpers.GetEnv("DB_MAX_IDLE_TIME", "5m"), 5*time.Minute),
		}
	}

	if caps.useCORS {
		cfg.CORS = &CORSConfig{
			AllowedOrigins:   helpers.Csv(helpers.GetEnv("CORS_ALLOWED_ORIGINS", "*")),
			AllowedMethods:   helpers.Csv(helpers.GetEnv("CORS_ALLOWED_METHODS", "GET,POST,PUT,PATCH,DELETE,OPTIONS")),
			AllowedHeaders:   helpers.Csv(helpers.GetEnv("CORS_ALLOWED_HEADERS", "Accept,Authorization,Content-Type,X-CSRF-Token")),
			ExposedHeaders:   helpers.Csv(helpers.GetEnv("CORS_EXPOSE_HEADERS", "")),
			AllowCredentials: helpers.MustBool(helpers.GetEnv("CORS_ALLOW_CREDENTIALS", "true"), true),
			MaxAge:           helpers.MustInt(helpers.GetEnv("CORS_MAX_AGE", "300"), 300),
		}
	}

	if caps.useJWT {
		cfg.JWT = &JWTConfig{
			Secret:     helpers.GetEnv("JWT_SECRET", "dev_secret_change_me"),
			Issuer:     helpers.GetEnv("JWT_ISSUER", "app"),
			AccessTTL:  helpers.MustDur(helpers.GetEnv("JWT_ACCESS_TTL", "15m"), 15*time.Minute),
			RefreshTTL: helpers.MustDur(helpers.GetEnv("JWT_REFRESH_TTL", "720h"), 30*24*time.Hour),
		}
	}

	if caps.useStorage {
		cfg.Storage = &StorageConfig{
			Driver:        helpers.GetEnv("STORAGE_DRIVER", "local"),
			LocalDir:      helpers.GetEnv("STORAGE_LOCAL_DIR", "./var/media"),
			PublicBaseURL: helpers.GetEnv("STORAGE_PUBLIC_BASE_URL", "http://localhost:8080"),
			PresignTTL:    helpers.MustDur(helpers.GetEnv("STORAGE_PRESIGN_TTL", "10m"), 10*time.Minute),
			S3Bucket:      helpers.GetEnv("S3_BUCKET", ""),
			S3Region:      helpers.GetEnv("S3_REGION", ""),
			S3Endpoint:    helpers.GetEnv("S3_ENDPOINT", ""),
			S3AccessKey:   helpers.GetEnv("S3_ACCESS_KEY", ""),
			S3SecretKey:   helpers.GetEnv("S3_SECRET_KEY", ""),
			S3UsePath:     helpers.MustBool(helpers.GetEnv("S3_USE_PATH_STYLE", "true"), true),
		}
	}

	if caps.useCache {
		cfg.Cache = &CacheConfig{
			Addr:         helpers.GetEnv("REDIS_ADDR", "redis:6379"),
			Password:     helpers.GetEnv("REDIS_PASSWORD", ""),
			DB:           helpers.MustInt(helpers.GetEnv("REDIS_DB", "0"), 0),
			PoolSize:     helpers.MustInt(helpers.GetEnv("REDIS_POOL", "50"), 50),
			DialTimeout:  5 * time.Second,
			ReadTimeout:  2 * time.Second,
			WriteTimeout: 2 * time.Second,
		}
	}

	if caps.useSMTP {
		cfg.SMTP = &SMTPConfig{
			Host:     helpers.GetEnv("SMTP_HOST", "mailpit"),
			Port:     helpers.MustInt(helpers.GetEnv("SMTP_PORT", "1025"), 1025),
			Username: helpers.GetEnv("SMTP_USERNAME", "root"),
			Password: helpers.GetEnv("SMTP_PASSWORD", "root"),
		}
	}

	if err := cfg.Validate(caps); err != nil {
		return nil, err
	}
	return cfg, nil
}
