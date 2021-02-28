package handlers

import (
	"bytes"
	"embed"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func StaticHandler(fs embed.FS) func(c *gin.Context) {
	return func(c *gin.Context) {
		fileParam := c.Param("file")
		if fileParam == "" {
			fileParam = "/index.html"
		} else {
			fileParam = "/static" + fileParam
		}

		fileUrl := "assets/src/build" + fileParam
		fileContent, err := fs.ReadFile(fileUrl)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		contentType := "text/html"
		if strings.HasSuffix(fileParam, "css") {
			contentType = "text/css"
		} else if strings.HasSuffix(fileParam, "js") {
			contentType = "text/javascript"
		}

		c.DataFromReader(http.StatusOK, int64(len(fileContent)), contentType, bytes.NewReader(fileContent), nil)
	}
}
