package webframework

import (
	"log"
	"strings"

	"github.com/dvwright/xss-mw"
	"github.com/gin-gonic/gin"
)

/**
 * @description: 启动xss过滤
 * @param {*gin.Engine} r
 * @return {*}
 */
func useXssMiddleware(r *gin.Engine) {
	log.Println("Use Xss Middleware!")
	xssMdwr := &xss.XssMw{}
	if webConfigCopy.Security.Xss.FieldsToSkip != "" {
		fields := strings.Split(webConfigCopy.Security.Xss.FieldsToSkip, ",")
		xssMdwr.FieldsToSkip = fields
		xssMdwr.BmPolicy = "UGCPolicy"
	}
	r.Use(xssMdwr.RemoveXss())
}
