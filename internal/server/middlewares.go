package server

import (
	"fmt"
	"net/http"
	"strings"

	"moefile/internal/meta"

	"github.com/gin-gonic/gin"
)

var (
	HTTPAllowedMethods = "GET, HEAD, OPTIONS"
	HTTPHeadersVary    = fmt.Sprintf("Origin, %s", CROSAllowedHeaders)
	CROSAllowedMethods = HTTPAllowedMethods
	CROSAllowedHeaders = "Range, If-Modified-Since"
	CROSExposeHeaders  = "*"
	CROSMaxAge         = map[bool]string{true: "3600", false: "0"}[meta.BuildMode == "production"]
)

func (s *serverConfig) serverInfoMiddleware(c *gin.Context) {
	c.Header("Server", fmt.Sprintf("%s/%s (%s)", meta.AppName, meta.AppVersion, s.app.ServerName))
	c.Header("Vary", HTTPHeadersVary)
}

func (s *serverConfig) crosMiddleware(c *gin.Context) {
	origin := c.GetHeader("Origin")
	if s.app.IsAllowedOrigin(origin) {
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Methods", CROSAllowedMethods)
		c.Header("Access-Control-Allow-Headers", CROSAllowedHeaders)
		c.Header("Access-Control-Expose-Headers", CROSExposeHeaders)
		c.Header("Access-Control-Max-Age", CROSMaxAge)
	}
	if c.Request.Method == http.MethodOptions {
		c.AbortWithStatus(http.StatusNoContent)
	}
}

func (s *serverConfig) methodNotAllowedMiddleware(c *gin.Context) {
	for _, allowed := range strings.Split(HTTPAllowedMethods, ",") {
		if c.Request.Method == strings.TrimSpace(allowed) {
			return
		}
	}
	c.Header("Allow", HTTPAllowedMethods)
	c.AbortWithStatus(http.StatusMethodNotAllowed)
}
