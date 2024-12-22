package server

import (
	"bufio"
	"bytes"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"moefile/dist"
	"moefile/internal/cfg"
	"moefile/internal/log"
	"moefile/pkg/dto"

	"github.com/gin-gonic/gin"
)

const (
	QueryPrefixVFS       = "_/"
	QueryPrefixVFSPlayer = "_/player/"
)

type serverConfig struct {
	app         cfg.AppConfig
	absRootPath string
	rootFS      fs.FS
	createdAt   time.Time
}

type handler struct {
	*serverConfig
	*urlInfo
	*gin.Context
}

func Setup(app cfg.AppConfig, e *gin.Engine) {
	absRootPath, err := filepath.Abs(app.RootPath)
	if err != nil {
		log.T("server").Errf("Unable to parse path <%s>: %s", app.RootPath, err)
		os.Exit(1)
	}

	cfg := serverConfig{
		app:         app,
		absRootPath: absRootPath,
		rootFS:      os.DirFS(absRootPath),
		createdAt:   time.Now(),
	}

	e.Use(cfg.serverInfoMiddleware)
	e.Use(cfg.crosMiddleware)
	e.Use(cfg.methodNotAllowedMiddleware)
	e.NoRoute(cfg.handle)
}

func (s *serverConfig) handle(c *gin.Context) {
	url := resolve(c.Request.URL.Path)
	if !url.ok {
		c.XML(http.StatusNotFound, dto.ErrorResponse{Message: "invalid url"})
		return
	}

	handler := handler{
		serverConfig: s,
		urlInfo:      &url,
		Context:      c,
	}

	if ok := handler.handlePlayer(); ok {
		log.T("server").Dbgf("Request handled by: player")
		return
	}

	if ok := handler.handleVFS(); ok {
		log.T("server").Dbgf("Request handled by: vfs")
		return
	}

	if ok := handler.handleXML(); ok {
		log.T("server").Dbgf("Request handled by: xml")
		return
	}

	if ok := handler.handleFile(); ok {
		log.T("server").Dbgf("Request handled by: file")
		return
	}

	c.XML(http.StatusNotFound, dto.ErrorResponse{Message: "not found"})
}

func (c *handler) handlePlayer() bool {
	query := c.Request.URL.RawQuery
	if c.requestURL != "/" || !strings.HasPrefix(query, QueryPrefixVFSPlayer) {
		return false
	}

	if c.Request.Method != "GET" {
		c.Header("Allow", "GET")
		c.XML(http.StatusMethodNotAllowed, dto.ErrorResponse{Message: "player: method not allowed"})
		c.Abort()
		return true
	}

	playerReqURL := resolve("/" + strings.Trim(strings.TrimPrefix(query, QueryPrefixVFSPlayer), "/"))
	if !playerReqURL.ok {
		c.XML(http.StatusNotFound, dto.ErrorResponse{Message: "player: invalid url"})
		c.Abort()
		return true
	}
	log.T("server/player").Dbgf("Player request URL:  %s", playerReqURL.requestURL)

	data, err := c.searchPlayerData(playerReqURL.requestURL)
	if err != nil {
		log.T("server/player").Errf("Unable to search danmaku and subtitles (PlayerData): %v", err)
		data = dto.PlayerData{
			Subs: make([]dto.PlayerSub, 0),
		}
	}

	buf, err := renderPlayerData(data)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return true
	}

	c.Status(http.StatusOK)
	c.Header("Content-Type", "text/html; charset=utf-8")
	_, err = bufio.NewReader(buf).WriteTo(c.Writer)
	c.Abort()
	if err != nil {
		log.T("server/player").Errf("Unable to write player data to response: %v", err)
	}
	return true
}

func (c *handler) handleVFS() bool {
	query := c.Request.URL.RawQuery
	if c.requestURL != "/" || !strings.HasPrefix(query, QueryPrefixVFS) {
		return false
	}

	url := strings.TrimPrefix(query, QueryPrefixVFS)
	buf, err := dist.Embed.ReadFile(url)
	if err != nil {
		log.T("server/vfs").Dbgf("Unable to open file <(vfs)/%s>: %s", url, err)
		c.XML(http.StatusNotFound, dto.ErrorResponse{Message: "vfs: file not found"})
		c.Abort()
		return true
	}

	http.ServeContent(c.Context.Writer, c.Request, url,
		cfg.AppDefaultBuildTime, bytes.NewReader(buf))
	return true
}

func (c *handler) handleXML() bool {
	file, err := c.rootFS.Open(c.relPath)
	if err != nil {
		log.T("server/xml").Dbgf("Unable to open file <(wwwroot)/%s>: %s", c.relPath, err)
		return false
	}

	stat, err := file.Stat()
	if err != nil {
		log.T("server/xml").Dbgf("Unable to stat file <(wwwroot)/%s>: %s", c.relPath, err)
		return false
	}

	if !stat.IsDir() {
		return false
	}

	if !strings.HasSuffix(c.Request.URL.Path, "/") {
		stdURL := c.requestURL + "/"
		log.T("server/xml").Dbgf("Redirecting to tailing slash URL: %s -> %s", c.Request.URL.Path, stdURL)
		c.Redirect(http.StatusTemporaryRedirect, stdURL)
		c.Abort()
		return true
	}

	buf, err := c.createS3XMLFromFSDir(c.requestURL, c.relPath)
	if err != nil {
		c.XML(http.StatusInternalServerError, dto.ErrorResponse{Message: "xml: server error"})
		c.Abort()
		return true
	}

	c.Status(http.StatusOK)
	c.Header("Content-Type", "application/xml; charset=utf-8")
	_, err = c.Writer.Write(buf)
	if err != nil {
		log.T("server/xml").Errf("Unable to write XML response: %v", err)
	}
	return true
}

func (c *handler) handleFile() bool {
	_, err := c.rootFS.Open(c.relPath)
	if err != nil {
		log.T("server/file").Dbgf("Unable to open file <(wwwroot)/%s>: %s", c.relPath, err)
		c.XML(http.StatusNotFound, dto.ErrorResponse{Message: "file not found"})
		c.Abort()
		return true
	}

	http.ServeFileFS(c.Context.Writer, c.Request, c.rootFS, c.relPath)
	return true
}
