package handler

import (
	"github.com/gin-gonic/gin"
	"logscan/pkg"
)

// 搜索页面日志信息
func LogMessageInfoHandler(ctx *gin.Context) {

	// 预检请求的处理
	startTime := ctx.Query("start")
	endTime := ctx.Query("end")
	lucene := ctx.Query("lucene")
	response := pkg.SearchLogFromEs(startTime, endTime, lucene)
	ctx.JSON(200, response.Hits.Hits)

}
