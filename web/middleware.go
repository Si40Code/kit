package web

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// loggingMiddleware 日志中间件
func (s *Server) loggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查是否跳过
		if s.shouldSkipPath(c.Request.URL.Path) {
			c.Next()
			return
		}

		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// 读取请求 body（处理文件上传）
		var reqBody string
		if !s.isMultipartForm(c) && c.Request.ContentLength > 0 {
			reqBody = s.readRequestBody(c)
		} else if s.isMultipartForm(c) {
			reqBody = fmt.Sprintf("[multipart/form-data, size: %d]", c.Request.ContentLength)
		}

		// 使用自定义 ResponseWriter 捕获响应
		blw := &bodyLogWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
			maxSize:        int(s.options.maxBodyLogSize),
		}
		c.Writer = blw

		// 处理请求
		c.Next()

		// 计算耗时
		latency := time.Since(start)
		statusCode := c.Writer.Status()

		// 读取响应 body
		respBody := blw.body.String()
		if len(respBody) > int(s.options.maxBodyLogSize) {
			respBody = fmt.Sprintf("[response too large, size: %d, showing first %d bytes] %s...",
				len(respBody), s.options.maxBodyLogSize, respBody[:s.options.maxBodyLogSize])
		}

		// 构建日志字段
		fields := map[string]interface{}{
			"client_ip":  c.ClientIP(),
			"method":     c.Request.Method,
			"path":       path,
			"query":      query,
			"status":     statusCode,
			"latency_ms": latency.Milliseconds(),
			"user_agent": c.Request.UserAgent(),
		}

		// 添加请求体（如果不是文件）
		if reqBody != "" {
			fields["req_body"] = reqBody
		}

		// 添加响应体
		if respBody != "" {
			fields["resp_body"] = respBody
		}

		// 添加错误信息
		if len(c.Errors) > 0 {
			fields["errors"] = c.Errors.String()
		}

		// 记录日志
		if s.options.logger != nil {
			logLevel := s.determineLogLevel(statusCode, latency)
			switch logLevel {
			case "error":
				s.options.logger.Error(c.Request.Context(), "HTTP request", fields)
			case "warn":
				s.options.logger.Warn(c.Request.Context(), "HTTP request", fields)
			default:
				s.options.logger.Info(c.Request.Context(), "HTTP request", fields)
			}
		}

		// 记录指标
		if s.options.enableMetric && s.options.metricRecorder != nil {
			s.options.metricRecorder.RecordRequest(MetricData{
				Method:   c.Request.Method,
				Path:     path,
				Status:   statusCode,
				Duration: latency,
			})
		}
	}
}

// recoveryMiddleware Panic 恢复中间件
func (s *Server) recoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 记录错误日志
				if s.options.logger != nil {
					s.options.logger.Error(c.Request.Context(), "Panic recovered", map[string]interface{}{
						"error":      fmt.Sprintf("%v", err),
						"path":       c.Request.URL.Path,
						"method":     c.Request.Method,
						"client_ip":  c.ClientIP(),
						"user_agent": c.Request.UserAgent(),
					})
				}

				// 返回 500 错误
				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    http.StatusInternalServerError,
					"message": "Internal server error",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}

// corsMiddleware CORS 中间件
func (s *Server) corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// Helper functions

func (s *Server) shouldSkipPath(path string) bool {
	for _, skip := range s.options.skipPaths {
		if path == skip {
			return true
		}
	}
	return false
}

func (s *Server) isMultipartForm(c *gin.Context) bool {
	contentType := c.GetHeader("Content-Type")
	return strings.HasPrefix(contentType, "multipart/form-data")
}

func (s *Server) readRequestBody(c *gin.Context) string {
	if c.Request.Body == nil {
		return ""
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return "[failed to read body]"
	}

	// 恢复 body 供后续处理使用
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	// 限制日志大小
	if int64(len(body)) > s.options.maxBodyLogSize {
		return fmt.Sprintf("[body too large, size: %d, showing first %d bytes] %s...",
			len(body), s.options.maxBodyLogSize, string(body[:s.options.maxBodyLogSize]))
	}

	return string(body)
}

func (s *Server) determineLogLevel(statusCode int, latency time.Duration) string {
	if statusCode >= 500 {
		return "error"
	}
	if statusCode >= 400 {
		return "warn"
	}
	if latency > s.options.slowRequestThresh {
		return "warn"
	}
	return "info"
}

// bodyLogWriter 用于捕获响应 body
type bodyLogWriter struct {
	gin.ResponseWriter
	body    *bytes.Buffer
	maxSize int
}

func (w *bodyLogWriter) Write(b []byte) (int, error) {
	// 限制缓存大小
	if w.body.Len() < w.maxSize {
		if w.body.Len()+len(b) > w.maxSize {
			w.body.Write(b[:w.maxSize-w.body.Len()])
		} else {
			w.body.Write(b)
		}
	}
	return w.ResponseWriter.Write(b)
}
