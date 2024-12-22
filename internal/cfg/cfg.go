package cfg

import (
	"flag"
	"strings"
	"time"

	"moefile/internal/meta"
	"moefile/pkg/logger"
)

const (
	Wildcard               = "*"
	WildcardCIDRListString = "0.0.0.0/0,::/0"
)

var (
	AppIsDevelopmentMode     = meta.BuildMode == "development"
	AppDefaultServerName     = meta.AppName
	AppDefaultListenAddr     = "0.0.0.0:3328" // 0x0d00 - Ciallo～(∠・ω< )⌒★
	AppDefaultRootPath       = ""
	AppDefaultLogLevel       = map[bool]string{true: "dbg", false: "inf"}[AppIsDevelopmentMode]
	AppDefaultAllowedOrigins = map[bool]string{true: Wildcard, false: ""}[AppIsDevelopmentMode]
	AppDefaultTrustedProxies = map[bool]string{true: WildcardCIDRListString, false: "127.0.0.1"}[AppIsDevelopmentMode]
	AppDefaultXMLIndent      = AppIsDevelopmentMode
	AppDefaultBuildTime      = parseBuildTime()
)

type AppConfig struct {
	ServerName     string
	ListenAddr     string
	RootPath       string
	LogLevel       string
	AllowedOrigins string
	TrustedProxies string
	XMLIndent      bool
}

func (cfg *AppConfig) IsDevelopmentMode() bool {
	return AppIsDevelopmentMode
}

func (cfg *AppConfig) ParseLogLevel() logger.LogLevel {
	switch strings.ToLower(cfg.LogLevel) {
	case "dbg":
		return logger.LDbg
	case "inf":
		return logger.LInf
	case "wrn":
		return logger.LWrn
	case "err":
		return logger.LErr
	default:
		return map[bool]logger.LogLevel{true: logger.LDbg, false: logger.LInf}[cfg.IsDevelopmentMode()]
	}
}

func (cfg *AppConfig) TrustedProxiesList() []string {
	trustedProxies := cfg.TrustedProxies
	if trustedProxies == Wildcard {
		trustedProxies = WildcardCIDRListString
	}
	return strings.Split(trustedProxies, ",")
}

func (cfg *AppConfig) IsAllowedOrigin(origin string) bool {
	allows := strings.Split(strings.TrimSpace(cfg.AllowedOrigins), ",")
	for _, allow := range allows {
		allow = strings.TrimSpace(allow)
		if allow == Wildcard || allow == origin {
			return true
		}
	}
	return false
}

func NewAppConfigFromFlag() AppConfig {
	serverName := flag.String("server", AppDefaultServerName, "app name")
	listenAddr := flag.String("listen", AppDefaultListenAddr, "listen address")
	rootPath := flag.String("root", AppDefaultRootPath, "server web root path")
	logLevel := flag.String("level", AppDefaultLogLevel, "log level, available values: dbg, inf, wrn, err")
	allowedOrigin := flag.String("origins", AppDefaultAllowedOrigins, "allowed CROS origins, split by comma")
	trustedProxies := flag.String("proxies", AppDefaultTrustedProxies, "trusted proxies, split by comma, or '*' for all")
	xmlIndent := flag.Bool("xmltab", AppDefaultXMLIndent, "pretty print JSON/XML in response")

	flag.Parse()
	return AppConfig{
		ServerName:     *serverName,
		ListenAddr:     *listenAddr,
		RootPath:       *rootPath,
		LogLevel:       *logLevel,
		AllowedOrigins: *allowedOrigin,
		TrustedProxies: *trustedProxies,
		XMLIndent:      *xmlIndent,
	}
}

func parseBuildTime() time.Time {
	t, err := time.Parse(time.RFC3339, meta.BuildTimestamp)
	if err == nil {
		return t
	}

	t, err = time.Parse(time.RFC3339Nano, meta.BuildTimestamp)
	if err == nil {
		return t
	}

	t, err = time.Parse("2006-01-02 15:04:05 -0700 MST", meta.BuildTimestamp)
	if err == nil {
		return t
	}

	t, err = time.Parse("2006-01-02 15:04:05 -0700", meta.BuildTimestamp)
	if err == nil {
		return t
	}

	return time.Now()
}
