package httpclient

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/go-resty/resty/v2"
)

// Client HTTP 客户端封装
type Client struct {
	*resty.Client
	options *options
}

// New 创建一个新的 HTTP 客户端
func New(opts ...Option) *Client {
	options := newOptions(opts...)

	// 创建自定义 transport
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * options.timeout,
			KeepAlive: options.keepAlive,
		}).DialContext,
		MaxIdleConns:          options.maxIdleConns,
		IdleConnTimeout:       options.idleConnTimeout,
		TLSHandshakeTimeout:   options.tlsHandshakeTimeout,
		ExpectContinueTimeout: 1 * options.timeout,
	}

	// TLS 配置
	if options.insecureSkipVerify {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	// 创建 resty 客户端
	restyClient := resty.NewWithClient(&http.Client{
		Transport: transport,
		Timeout:   options.timeout,
	})

	// 启用 trace 信息收集
	restyClient.EnableTrace()

	// 设置 resty 的调试输出（如果有 logger）
	if options.enableLog && options.logger != nil {
		restyClient.SetDebug(false) // 禁用默认的 stderr 输出
		// resty 会通过钩子函数记录日志，而不是直接输出
	}

	// 设置重试
	if options.retryCount > 0 {
		restyClient.
			SetRetryCount(options.retryCount).
			SetRetryWaitTime(options.retryWaitTime).
			SetRetryMaxWaitTime(options.retryMaxWaitTime)
	}

	client := &Client{
		Client:  restyClient,
		options: options,
	}

	// 注册钩子
	client.setupHooks()

	return client
}

// setupHooks 设置请求/响应钩子
func (c *Client) setupHooks() {
	// OnBeforeRequest: 请求前处理
	c.OnBeforeRequest(func(client *resty.Client, req *resty.Request) error {
		ctx := req.Context()

		// Trace: 创建 span
		if c.options.enableTrace {
			newCtx, span := createSpan(ctx, c.options.serviceName, req.Method, req.URL)
			req.SetContext(newCtx)

			// 设置初始属性
			if span.IsRecording() {
				attrs := map[string]interface{}{
					"http.request.body_size": len(fmt.Sprintf("%v", req.Body)),
				}
				if len(req.Header) > 0 {
					attrs["http.request.headers"] = fmt.Sprintf("%v", req.Header)
				}
				setSpanAttributes(span, attrs)
			}
		}

		// Log: 记录请求开始
		if c.options.enableLog && c.options.logger != nil {
			c.logRequest(req)
		}

		return nil
	})

	// OnAfterResponse: 响应后处理
	c.OnAfterResponse(func(client *resty.Client, resp *resty.Response) error {
		ctx := resp.Request.Context()
		ti := resp.Request.TraceInfo()

		// Trace: 更新 span
		if c.options.enableTrace {
			span := getSpanFromContext(ctx)
			if span.IsRecording() {
				// 设置响应属性
				attrs := map[string]interface{}{
					"http.status_code":        resp.StatusCode(),
					"http.response.body_size": len(resp.Body()),
					"http.response.time_ms":   ti.TotalTime.Milliseconds(),
					"http.dns_lookup_ms":      ti.DNSLookup.Milliseconds(),
					"http.tcp_conn_ms":        ti.TCPConnTime.Milliseconds(),
					"http.tls_handshake_ms":   ti.TLSHandshake.Milliseconds(),
					"http.server_time_ms":     ti.ServerTime.Milliseconds(),
					"http.conn_reused":        ti.IsConnReused,
				}
				setSpanAttributes(span, attrs)

				// 根据状态码判断成功/失败
				if resp.StatusCode() >= 400 {
					markSpanError(span, nil, fmt.Sprintf("HTTP %d", resp.StatusCode()))
				} else {
					markSpanSuccess(span)
				}
			}
			span.End()
		}

		// Log: 记录响应
		if c.options.enableLog && c.options.logger != nil {
			c.logResponse(resp)
		}

		// Metric: 记录指标
		if c.options.enableMetric && c.options.metricRecorder != nil {
			c.recordMetric(resp)
		}

		return nil
	})

	// OnError: 错误处理
	c.OnError(func(req *resty.Request, err error) {
		ctx := req.Context()
		ti := req.TraceInfo()

		// Trace: 标记错误
		if c.options.enableTrace {
			span := getSpanFromContext(ctx)
			if span.IsRecording() {
				attrs := map[string]interface{}{
					"http.error":         err.Error(),
					"http.total_time_ms": ti.TotalTime.Milliseconds(),
				}
				setSpanAttributes(span, attrs)
				markSpanError(span, err, "HTTP request failed")
			}
			span.End()
		}

		// Log: 记录错误
		if c.options.enableLog && c.options.logger != nil {
			c.logError(req, err)
		}

		// Metric: 记录失败指标
		if c.options.enableMetric && c.options.metricRecorder != nil {
			c.recordErrorMetric(req, err)
		}
	})
}

// logRequest 记录请求日志
func (c *Client) logRequest(req *resty.Request) {
	fields := map[string]interface{}{
		"http.method": req.Method,
		"http.url":    req.URL,
	}

	// 记录请求头（过滤敏感信息）
	if len(req.Header) > 0 {
		safeHeaders := c.filterSensitiveHeaders(req.Header)
		fields["http.request.headers"] = safeHeaders
	}

	// 记录查询参数
	if len(req.QueryParam) > 0 {
		fields["http.request.query_params"] = req.QueryParam
	}

	// 记录表单数据
	if len(req.FormData) > 0 {
		fields["http.request.form_data"] = req.FormData
	}

	// 检查是否是文件上传
	contentType := req.Header.Get("Content-Type")
	isMultipart := c.isMultipartRequest(contentType)

	if isMultipart {
		// 文件上传：只记录元信息
		fields["http.request.content_type"] = contentType
		if req.RawRequest != nil && req.RawRequest.ContentLength > 0 {
			fields["http.request.content_length"] = req.RawRequest.ContentLength
			fields["http.request.body"] = fmt.Sprintf("[multipart/form-data, size: %d bytes]", req.RawRequest.ContentLength)
		} else {
			fields["http.request.body"] = "[multipart/form-data]"
		}
	} else if req.Body != nil {
		// 普通请求：记录 body 内容（限制大小）
		bodyStr := fmt.Sprintf("%v", req.Body)
		if int64(len(bodyStr)) > c.options.maxBodyLogSize {
			fields["http.request.body"] = bodyStr[:c.options.maxBodyLogSize] + "..."
			fields["http.request.body_truncated"] = true
			fields["http.request.body_original_size"] = len(bodyStr)
		} else {
			fields["http.request.body"] = bodyStr
		}
	}

	c.options.logger.InfoMap(req.Context(), "HTTP request started", fields)
}

// logResponse 记录响应日志
func (c *Client) logResponse(resp *resty.Response) {
	ti := resp.Request.TraceInfo()

	remoteAddr := ""
	if ti.RemoteAddr != nil {
		remoteAddr = ti.RemoteAddr.String()
	}

	fields := map[string]interface{}{
		"http.method":            resp.Request.Method,
		"http.url":               resp.Request.URL,
		"http.status_code":       resp.StatusCode(),
		"http.status":            resp.Status(),
		"http.proto":             resp.Proto(),
		"http.total_time_ms":     ti.TotalTime.Milliseconds(),
		"http.dns_lookup_ms":     ti.DNSLookup.Milliseconds(),
		"http.tcp_conn_ms":       ti.TCPConnTime.Milliseconds(),
		"http.tls_handshake_ms":  ti.TLSHandshake.Milliseconds(),
		"http.server_time_ms":    ti.ServerTime.Milliseconds(),
		"http.response_time_ms":  ti.ResponseTime.Milliseconds(),
		"http.conn_reused":       ti.IsConnReused,
		"http.conn_was_idle":     ti.IsConnWasIdle,
		"http.conn_idle_time_ms": ti.ConnIdleTime.Milliseconds(),
		"http.request_attempt":   ti.RequestAttempt,
		"http.remote_addr":       remoteAddr,
	}

	// 记录响应头（过滤敏感信息）
	if len(resp.Header()) > 0 {
		safeHeaders := c.filterSensitiveHeaders(resp.Header())
		fields["http.response.headers"] = safeHeaders
	}

	// 检查是否是文件下载
	contentType := resp.Header().Get("Content-Type")
	contentDisposition := resp.Header().Get("Content-Disposition")
	isFileDownload := c.isFileResponse(contentType, contentDisposition)

	if isFileDownload {
		// 文件下载：只记录元信息
		fields["http.response.content_type"] = contentType
		fields["http.response.content_length"] = len(resp.Body())
		fields["http.response.body"] = fmt.Sprintf("[file download, size: %d bytes]", len(resp.Body()))

		if contentDisposition != "" {
			fields["http.response.content_disposition"] = contentDisposition
		}
	} else {
		// 普通响应：记录 body 内容（限制大小）
		respBody := string(resp.Body())
		bodyLen := int64(len(respBody))

		if bodyLen > c.options.maxBodyLogSize {
			fields["http.response.body"] = respBody[:c.options.maxBodyLogSize] + "..."
			fields["http.response.body_truncated"] = true
			fields["http.response.body_original_size"] = bodyLen
		} else if bodyLen > 0 {
			fields["http.response.body"] = respBody
		}

		fields["http.response.content_type"] = contentType
	}

	// 读取请求体（如果需要）
	contentType = resp.Request.Header.Get("Content-Type")
	isMultipart := c.isMultipartRequest(contentType)

	if !isMultipart && resp.Request.RawRequest.GetBody != nil {
		if bodyReader, err := resp.Request.RawRequest.GetBody(); err == nil {
			if bodyBytes, err := io.ReadAll(bodyReader); err == nil {
				if int64(len(bodyBytes)) > c.options.maxBodyLogSize {
					fields["http.request.body"] = string(bodyBytes[:c.options.maxBodyLogSize]) + "..."
					fields["http.request.body_truncated"] = true
				} else if len(bodyBytes) > 0 {
					fields["http.request.body"] = string(bodyBytes)
				}
			}
			bodyReader.Close()
		}
	} else if isMultipart {
		fields["http.request.body"] = "[multipart/form-data]"
	}

	// 根据状态码选择日志级别
	if resp.StatusCode() >= 500 {
		c.options.logger.ErrorMap(resp.Request.Context(), "HTTP request completed with server error", fields)
	} else if resp.StatusCode() >= 400 {
		c.options.logger.WarnMap(resp.Request.Context(), "HTTP request completed with client error", fields)
	} else {
		c.options.logger.InfoMap(resp.Request.Context(), "HTTP request completed successfully", fields)
	}
}

// logError 记录错误日志
func (c *Client) logError(req *resty.Request, err error) {
	ti := req.TraceInfo()

	fields := map[string]interface{}{
		"http.method":           req.Method,
		"http.url":              req.URL,
		"http.error":            err.Error(),
		"http.total_time_ms":    ti.TotalTime.Milliseconds(),
		"http.dns_lookup_ms":    ti.DNSLookup.Milliseconds(),
		"http.tcp_conn_ms":      ti.TCPConnTime.Milliseconds(),
		"http.tls_handshake_ms": ti.TLSHandshake.Milliseconds(),
		"http.request_attempt":  ti.RequestAttempt,
	}

	c.options.logger.ErrorMap(req.Context(), "HTTP request failed", fields)
}

// recordMetric 记录响应指标
func (c *Client) recordMetric(resp *resty.Response) {
	ti := resp.Request.TraceInfo()

	remoteAddr := ""
	if ti.RemoteAddr != nil {
		remoteAddr = ti.RemoteAddr.String()
	}

	data := MetricData{
		Method:         resp.Request.Method,
		Host:           resp.Request.URL,
		Path:           resp.Request.RawRequest.URL.Path,
		StatusCode:     resp.StatusCode(),
		TotalTime:      ti.TotalTime,
		DNSLookup:      ti.DNSLookup,
		TCPConn:        ti.TCPConnTime,
		TLSHandshake:   ti.TLSHandshake,
		ServerTime:     ti.ServerTime,
		ResponseTime:   ti.ResponseTime,
		IsConnReused:   ti.IsConnReused,
		IsConnWasIdle:  ti.IsConnWasIdle,
		ConnIdleTime:   ti.ConnIdleTime,
		RequestAttempt: ti.RequestAttempt,
		RemoteAddr:     remoteAddr,
	}

	c.options.metricRecorder.RecordRequest(data)
}

// recordErrorMetric 记录错误指标
func (c *Client) recordErrorMetric(req *resty.Request, err error) {
	ti := req.TraceInfo()

	data := MetricData{
		Method:         req.Method,
		Host:           req.URL,
		Path:           req.RawRequest.URL.Path,
		StatusCode:     0, // 错误情况下状态码为 0
		TotalTime:      ti.TotalTime,
		DNSLookup:      ti.DNSLookup,
		TCPConn:        ti.TCPConnTime,
		TLSHandshake:   ti.TLSHandshake,
		RequestAttempt: ti.RequestAttempt,
	}

	c.options.metricRecorder.RecordRequest(data)
}

// R 创建一个新的请求，自动传递 context（如果设置了）
func (c *Client) R(ctx context.Context) *resty.Request {
	return c.Client.R().SetContext(ctx)
}

// isMultipartRequest 检查是否是 multipart 请求（文件上传）
func (c *Client) isMultipartRequest(contentType string) bool {
	return contentType != "" && strings.HasPrefix(contentType, "multipart/")
}

// isFileResponse 检查响应是否是文件下载
func (c *Client) isFileResponse(contentType, contentDisposition string) bool {
	// 检查 Content-Disposition 是否包含 attachment 或 inline
	if contentDisposition != "" {
		lower := strings.ToLower(contentDisposition)
		if strings.HasPrefix(lower, "attachment") || strings.HasPrefix(lower, "inline") {
			return true
		}
	}

	// 检查常见的文件 MIME 类型
	if contentType == "" {
		return false
	}

	fileTypes := []string{
		"application/octet-stream",
		"application/pdf",
		"application/zip",
		"application/x-rar-compressed",
		"application/x-7z-compressed",
		"application/x-tar",
		"application/gzip",
		"image/",
		"video/",
		"audio/",
		"application/vnd.ms-excel",
		"application/vnd.openxmlformats-officedocument",
		"application/msword",
	}

	lower := strings.ToLower(contentType)
	for _, fileType := range fileTypes {
		if strings.HasPrefix(lower, fileType) {
			return true
		}
	}

	return false
}

// filterSensitiveHeaders 过滤敏感的请求头信息
func (c *Client) filterSensitiveHeaders(headers http.Header) map[string]string {
	safeHeaders := make(map[string]string)
	sensitiveKeys := []string{
		"authorization",
		"cookie",
		"set-cookie",
		"x-api-key",
		"api-key",
		"apikey",
		"token",
		"access-token",
		"refresh-token",
		"x-auth-token",
		"x-csrf-token",
		"password",
		"secret",
	}

	for key, values := range headers {
		lowerKey := strings.ToLower(key)

		// 检查是否是敏感的 key
		isSensitive := false
		for _, sensitiveKey := range sensitiveKeys {
			if lowerKey == sensitiveKey || strings.HasPrefix(lowerKey, sensitiveKey) {
				isSensitive = true
				break
			}
		}

		if isSensitive {
			safeHeaders[key] = "******"
		} else {
			if len(values) > 0 {
				safeHeaders[key] = values[0]
			}
		}
	}

	return safeHeaders
}
