package main

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/Si40Code/kit/web"
	"github.com/gin-gonic/gin"
)

func main() {
	server := web.New(
		web.WithMode(web.DebugMode),
		web.WithServiceName("file-upload-example"),
		web.WithMaxMultipartMemory(64 << 20), // 64MB
	)

	engine := server.Engine()

	// 单文件上传
	engine.POST("/upload/single", func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			web.Error(c, http.StatusBadRequest, "No file provided")
			return
		}

		// 保存文件
		filename := filepath.Base(file.Filename)
		if err := c.SaveUploadedFile(file, "./uploads/"+filename); err != nil {
			web.Error(c, http.StatusInternalServerError, err.Error())
			return
		}

		web.Success(c, gin.H{
			"filename": filename,
			"size":     file.Size,
		})
	})

	// 多文件上传
	engine.POST("/upload/multiple", func(c *gin.Context) {
		form, err := c.MultipartForm()
		if err != nil {
			web.Error(c, http.StatusBadRequest, err.Error())
			return
		}

		files := form.File["files"]
		var savedFiles []map[string]interface{}

		for _, file := range files {
			filename := filepath.Base(file.Filename)
			if err := c.SaveUploadedFile(file, "./uploads/"+filename); err != nil {
				web.Error(c, http.StatusInternalServerError, err.Error())
				return
			}
			savedFiles = append(savedFiles, map[string]interface{}{
				"filename": filename,
				"size":     file.Size,
			})
		}

		web.Success(c, gin.H{
			"count": len(savedFiles),
			"files": savedFiles,
		})
	})

	// 启动服务器
	fmt.Println("File upload server starting on :8080")
	server.RunWithGracefulShutdown(":8080")
}
