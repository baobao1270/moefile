package main

import (
	"moefile/internal/cfg"
	"moefile/internal/log"
	"moefile/internal/meta"
	"moefile/internal/server"
	"os"

	"github.com/gin-gonic/gin"
)

var (
	DefaultFlagName         = meta.AppName
	DefaultFlagRoot         = ""
	DefaultFlagAddr         = "0.0.0.0:3328" // 0x0d00 - Ciallo～(∠・ω< )⌒★
	DefalutFlagCROS         = "*"
	DefaultFlagIndent       = false
	DefaultFlagTrustProxies = "off"
)

func main() {
	app := cfg.NewAppConfigFromFlag()
	log.Setup(app.ParseLogLevel())
	log.T("main").Inff("%s %s (Build %s)", meta.AppName, meta.AppVersion, meta.BuildTimestamp)
	log.T("main").Inff("Copyrigyt (c) %s %s, distributed under the %s license",
		meta.AppCopyRight, meta.AppAuthor, meta.AppLicense)
	log.T("main").Inff(" - Build mode: %s", meta.BuildMode)
	log.T("main").Inff(" - Server name: %s", app.ServerName)
	log.T("main").Inff(" - Log level: %s (0x%02x)", app.LogLevel, app.ParseLogLevel())

	log.SetupGin1()
	if meta.BuildMode == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	e := gin.New()
	log.SetupGin2(e)
	err := e.SetTrustedProxies(app.TrustedProxiesList())
	if err != nil {
		log.T("main").Errf("Failed to set trusted proxies: %v", err)
		os.Exit(1)
	}

	log.T("main").Inff("Server is listening on http://%s", app.ListenAddr)
	server.Setup(app, e)
	err = e.Run(app.ListenAddr)
	if err != nil {
		log.T("main").Errf("Failed to start server: %v", err)
		os.Exit(1)
	}
}
