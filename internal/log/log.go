package log

import (
	"fmt"
	"moefile/pkg/logger"

	"github.com/gin-gonic/gin"
)

var AppLogger = logger.NewStdout()

func T(name string) *logger.Tag {
	return AppLogger.Tag(name)
}

func Setup(minLevel logger.LogLevel) {
	AppLogger.MinLevel = minLevel
	AppLogger.TagColor = map[string]logger.LogColor{
		"gin":    logger.CMagenta,
		"main":   logger.CYellow,
		"http":   logger.CGreen,
		"server": logger.CBlue,
	}
}

func SetupGin1() {
	gin.DebugPrintFunc = T("gin").Wrnf
	gin.DefaultWriter = T("http").LogWriter(logger.LInf)
}

func SetupGin2(e *gin.Engine) {
	e.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		statusCode := fmt.Sprintf("%d", param.StatusCode)
		statusColor := logger.CReset
		if param.StatusCode >= 400 {
			statusColor = logger.CRed
		}
		if param.StatusCode >= 300 && param.StatusCode < 400 {
			statusColor = logger.CYellow
		}
		if param.StatusCode >= 200 && param.StatusCode < 300 {
			statusColor = logger.CGreen
		}
		coloredStatausCode := fmt.Sprint(statusColor, statusCode, logger.CReset)
		return fmt.Sprint(
			fmt.Sprintf("%-15s", param.ClientIP), " ",
			coloredStatausCode, " ",
			fmt.Sprintf("%-6s", param.Method), " ",
			param.Path, " ",
			"t=", param.Latency, " ",
			"ua=", param.Request.UserAgent(), " ",
			"msg=", param.ErrorMessage, "\n",
		)
	}))
}
