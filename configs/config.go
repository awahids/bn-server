package configs

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv      string
	Server      ServerConfig
	DB          DBConfig
	JWT         JWTConfig
	Google      GoogleConfig
	CORS        CORSConfig
	AuthCookie  AuthCookieConfig
	RateLimiter RateLimiterConfig
	AI          AIConfig
	Push        PushConfig
}

type ServerConfig struct {
	Host string
	Port string
}

type DBConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	Name            string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type JWTConfig struct {
	Issuer          string
	Secret          string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

type GoogleConfig struct {
	ClientID     string
	ClientSecret string
}

type CORSConfig struct {
	AllowOrigins []string
}

type AuthCookieConfig struct {
	RefreshTokenName     string
	RefreshTokenPath     string
	RefreshTokenDomain   string
	RefreshTokenSecure   bool
	RefreshTokenMaxAge   int
	RefreshTokenHTTPOnly bool
	RefreshTokenSameSite http.SameSite
}

type RateLimiterConfig struct {
	AuthWindowSeconds int
	AuthGoogleLimit   int
	AuthRefreshLimit  int
	AuthLogoutLimit   int
}

type AIConfig struct {
	APIKey  string
	BaseURL string
	Model   string
}

type PushConfig struct {
	Enabled          bool
	VAPIDPublicKey   string
	VAPIDPrivateKey  string
	VAPIDSubject     string
	DispatchInterval time.Duration
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load("server/.env")
	_ = godotenv.Load(".env")

	cfg := &Config{
		AppEnv: getEnv("APP_ENV", "development"),
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Port: getEnv("SERVER_PORT", "8080"),
		},
		DB: DBConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnv("DB_PORT", "5432"),
			User:            getEnv("DB_USER", "postgres"),
			Password:        getEnv("DB_PASSWORD", "postgres"),
			Name:            getEnv("DB_NAME", "bn_mobile"),
			SSLMode:         getEnv("DB_SSL_MODE", "disable"),
			MaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 20),
			MaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getEnvDuration("DB_CONN_MAX_LIFETIME", 30*time.Minute),
		},
		JWT: JWTConfig{
			Issuer:          getEnv("JWT_ISSUER", "bn-mobile-server"),
			Secret:          getEnv("JWT_SECRET", ""),
			AccessTokenTTL:  getEnvDuration("JWT_ACCESS_TOKEN_TTL", 15*time.Minute),
			RefreshTokenTTL: getEnvDuration("JWT_REFRESH_TOKEN_TTL", 7*24*time.Hour),
		},
		Google: GoogleConfig{
			ClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
			ClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		},
		CORS: CORSConfig{
			AllowOrigins: parseCommaSeparated(getEnv("CORS_ALLOW_ORIGINS", "http://localhost:3000")),
		},
		AuthCookie: AuthCookieConfig{
			RefreshTokenName:     getEnv("AUTH_REFRESH_COOKIE_NAME", "bn_refresh_token"),
			RefreshTokenPath:     getEnv("AUTH_REFRESH_COOKIE_PATH", "/api/v1/auth"),
			RefreshTokenDomain:   getEnv("AUTH_REFRESH_COOKIE_DOMAIN", ""),
			RefreshTokenSecure:   getEnvBool("AUTH_REFRESH_COOKIE_SECURE", false),
			RefreshTokenMaxAge:   getEnvInt("AUTH_REFRESH_COOKIE_MAX_AGE", 7*24*60*60),
			RefreshTokenHTTPOnly: true,
			RefreshTokenSameSite: parseSameSite(getEnv("AUTH_REFRESH_COOKIE_SAMESITE", "lax")),
		},
		RateLimiter: RateLimiterConfig{
			AuthWindowSeconds: getEnvInt("AUTH_RATE_LIMIT_WINDOW_SECONDS", 60),
			AuthGoogleLimit:   getEnvInt("AUTH_RATE_LIMIT_GOOGLE", 10),
			AuthRefreshLimit:  getEnvInt("AUTH_RATE_LIMIT_REFRESH", 30),
			AuthLogoutLimit:   getEnvInt("AUTH_RATE_LIMIT_LOGOUT", 30),
		},
		AI: AIConfig{
			APIKey:  getEnv("AI_API_KEY", ""),
			BaseURL: getEnv("AI_BASE_URL", "https://ai.sumopod.com/v1"),
			Model:   getEnv("AI_MODEL", "glm-ocr"),
		},
		Push: PushConfig{
			Enabled:          getEnvBool("PUSH_ENABLED", false),
			VAPIDPublicKey:   getEnv("PUSH_VAPID_PUBLIC_KEY", ""),
			VAPIDPrivateKey:  getEnv("PUSH_VAPID_PRIVATE_KEY", ""),
			VAPIDSubject:     getEnv("PUSH_VAPID_SUBJECT", "mailto:admin@example.com"),
			DispatchInterval: getEnvDuration("PUSH_DISPATCH_INTERVAL", time.Minute),
		},
	}

	if cfg.JWT.Secret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	if cfg.AppEnv == "production" && cfg.Google.ClientID == "" {
		return nil, fmt.Errorf("GOOGLE_CLIENT_ID is required in production")
	}
	if cfg.AppEnv == "production" && cfg.Google.ClientSecret == "" {
		return nil, fmt.Errorf("GOOGLE_CLIENT_SECRET is required in production")
	}

	if cfg.Push.Enabled {
		if strings.TrimSpace(cfg.Push.VAPIDPublicKey) == "" || strings.TrimSpace(cfg.Push.VAPIDPrivateKey) == "" {
			return nil, fmt.Errorf("PUSH_VAPID_PUBLIC_KEY and PUSH_VAPID_PRIVATE_KEY are required when PUSH_ENABLED=true")
		}
		if cfg.Push.DispatchInterval < 10*time.Second {
			return nil, fmt.Errorf("PUSH_DISPATCH_INTERVAL must be at least 10s")
		}
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return parsed
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	parsed, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}

	return parsed
}

func getEnvBool(key string, fallback bool) bool {
	value := strings.TrimSpace(strings.ToLower(os.Getenv(key)))
	if value == "" {
		return fallback
	}

	switch value {
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	default:
		return fallback
	}
}

func parseCommaSeparated(value string) []string {
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			out = append(out, trimmed)
		}
	}
	if len(out) == 0 {
		return []string{"http://localhost:3000"}
	}
	return out
}

func parseSameSite(value string) http.SameSite {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "strict":
		return http.SameSiteStrictMode
	case "none":
		return http.SameSiteNoneMode
	case "lax":
		fallthrough
	default:
		return http.SameSiteLaxMode
	}
}
