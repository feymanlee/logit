package logit

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
)

const (
	// 默认 logger name
	defaultGinLoggerName = "access"
	// 默认慢请求时间 3s
	defaultGinSlowThreshold = time.Second * 3
	// prometheus namespace
	promNamespace = "logit"
)

var (
	// gin prometheus labels
	promGinLabels = []string{
		"status_code",
		"path",
		"method",
	}
	promGinReqCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: promNamespace,
			Name:      "req_count",
			Help:      "gin server request count",
		}, promGinLabels,
	)
	promGinReqLatency = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: promNamespace,
			Name:      "req_latency",
			Help:      "gin server request latency in seconds",
		}, promGinLabels,
	)
)

//
// GetGinTraceIDFromHeader
//  @Description: 从 gin 的 request header 中获取 key 为 TraceIDKeyName 的值作为 traceid
//  @param c
//  @return string
//
func GetGinTraceIDFromHeader(c *gin.Context) string {
	return c.Request.Header.Get(string(TraceIDKeyName))
}

//
// GetGinTraceIDFromQueryString
//  @Description: 从 gin 的 querystring 中获取 key 为 TraceIDKeyName 的值作为 traceid
//  @param c
//  @return string
//
func GetGinTraceIDFromQueryString(c *gin.Context) string {
	return c.Query(string(TraceIDKeyName))
}

//
// GetGinTraceIDFromPostForm
//  @Description: 从 gin 的 post form 中获取 key 为 TraceIDKeyName 的值作为 traceid
//  @param c
//  @return string
//
func GetGinTraceIDFromPostForm(c *gin.Context) string {
	return c.PostForm(string(TraceIDKeyName))
}

// GinLogExtends gin 日志中间件记录的扩展
type GinLogExtends struct {
	// 请求处理耗时 (秒)
	Latency    float64 `json:"latency_seconds"`
	HandleName string  `json:"handle_name"`
}

// GinLoggerConfig GinLogger 支持的配置项字段定义
type GinLoggerConfig struct {
	Name string
	// Optional. Default value is logit.defaultGinLogFormatter
	Formatter func(*gin.Context, GinLogExtends) string
	// SkipPaths is a url path array which logs are not written.
	// Optional.
	SkipPaths []string
	// SkipPathRegexps skip path by regexp
	SkipPathRegexps []string
	// TraceIDFunc 获取或生成 trace id 的函数
	// Optional.
	TraceIDFunc func(*gin.Context) string
	// InitFieldsFunc 获取 baseLogger 初始字段方法 key 为字段名 value 为字段值
	InitFieldsFunc func(*gin.Context) map[string]interface{}
	// 是否使用详细模式打印日志，记录更多字段信息
	// Optional.
	EnableDetails bool
	// 以下选项开启后对性能有影响，适用于接口调试，慎用。
	// 是否打印 context keys
	// Optional.
	EnableContextKeys bool
	// 是否打印请求头信息
	// Optional.
	EnableRequestHeader bool
	// 是否打印请求form信息
	// Optional.
	EnableRequestForm bool
	// 是否打印请求体信息
	// Optional.
	EnableRequestBody bool
	// 是否打印响应体信息
	// Optional.
	EnableResponseBody bool

	// 慢请求时间阈值 请求处理时间超过该值则使用 Error 级别打印日志
	SlowThreshold time.Duration
}

//
// GinLogger
//  @Description: 以默认配置生成 gin 的 Logger 中间件
//  @return gin.HandlerFunc
//
func GinLogger() gin.HandlerFunc {
	return GinLoggerWithConfig(GinLoggerConfig{})
}

//
// defaultGinLogFormatter
//  @Description: 默认访问日志中 msg 字段的输出格式
//  @param c
//  @param m
//  @return string
//
func defaultGinLogFormatter(c *gin.Context, ext GinLogExtends) string {
	msg := fmt.Sprintf("%s|%s|%s%s|%d|%f",
		c.ClientIP(),
		c.Request.Method,
		c.Request.Host,
		c.Request.RequestURI,
		c.Writer.Status(),
		ext.Latency,
	)
	return msg
}

//
// defaultGinTraceIDFunc
//  @Description: 默认从 context 中获取 traceID 的方法
//  @param c
//  @return traceID
//
func defaultGinTraceIDFunc(c *gin.Context) (traceID string) {
	traceID = GetGinTraceIDFromHeader(c)
	if traceID != "" {
		return
	}
	traceID = GetGinTraceIDFromPostForm(c)
	if traceID != "" {
		return
	}
	traceID = GetGinTraceIDFromQueryString(c)
	if traceID != "" {
		return
	}
	traceID = CtxTraceID(c)
	return
}

//
// GinLoggerWithConfig
//  @Description:  根据配置信息生成 gin 的 Logger 中间件
// 中间件会记录访问信息，根据状态码确定日志级别， 500 以上为 Error ， 400-500 默认为 Warn ， 400 以下默认为 Info
// api 请求进来的 context 的函数无需在其中打印 err ，使用 c.Error(err)会在请求完成时自动打印 error
// context 中有 error 则日志忽略返回码始终使用 error 级别
//  @param conf
//  @return gin.HandlerFunc
//
func GinLoggerWithConfig(conf GinLoggerConfig) gin.HandlerFunc {
	formatter := conf.Formatter
	if formatter == nil {
		formatter = defaultGinLogFormatter
	}
	getTraceID := conf.TraceIDFunc
	if getTraceID == nil {
		getTraceID = defaultGinTraceIDFunc
	}

	var skipRegexps []*regexp.Regexp
	for _, p := range conf.SkipPathRegexps {
		if r, err := regexp.Compile(p); err != nil {
			panic("skip path regexps compile " + p + " error:" + err.Error())
		} else {
			skipRegexps = append(skipRegexps, r)
		}
	}

	if conf.SlowThreshold.Seconds() <= 0 {
		conf.SlowThreshold = defaultGinSlowThreshold
	}
	ginLogger := CloneLogger(conf.Name)

	var ginLogExtends = GinLogExtends{}

	return func(c *gin.Context) {
		if skipLog(c.Request.URL.Path, conf.SkipPaths, skipRegexps) {
			c.Next()
			return
		}
		start := time.Now()

		traceID := getTraceID(c)
		// 设置 trace id 到 request header 中
		c.Request.Header.Set(string(TraceIDKeyName), traceID)
		// 设置 trace id 到 response header 中
		c.Writer.Header().Set(string(TraceIDKeyName), traceID)
		// 设置 trace id 和 ctxLogger 到 context 中
		if conf.InitFieldsFunc != nil {
			for k, v := range conf.InitFieldsFunc(c) {
				ginLogger = ginLogger.With(zap.Any(k, v))
			}
		}
		_, ctxLogger := NewCtxLogger(c, ginLogger, traceID)
		_, shortHandlerName := path.Split(c.HandlerName())
		ginLogExtends.HandleName = shortHandlerName
		// 创建基础 logger，可以记录基础的信息
		accessLogger := ctxLogger.With(
			zap.Time("req_time", start),
			zap.String("client_ip", c.ClientIP()),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("host", c.Request.Host),
			zap.String("handle", shortHandlerName),
		)

		if conf.Name == "" {
			conf.Name = defaultGinLoggerName
		}
		// 判断是否打印请求 header
		if conf.EnableRequestHeader {
			accessLogger = accessLogger.With(zap.Any("request_header", c.Request.Header))
		}
		// 判断是否打印请求 form
		if conf.EnableRequestForm {
			accessLogger = accessLogger.With(zap.Any("request_form", c.Request.Form))
		}
		// 判断是否打印请求 body
		if conf.EnableRequestBody {
			accessLogger = accessLogger.With(zap.Any("request_body", string(GetGinRequestBody(c))))
		}
		rspBodyWriter := &responseBodyWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		if conf.EnableResponseBody {
			// 开启记录响应 body 时，保存 body 到 rspBodyWriter.body 中
			c.Writer = rspBodyWriter
		}

		defer func() {
			ginLogExtends.Latency = time.Since(start).Seconds()
			// 记录 status code 和 latency
			accessLogger = accessLogger.With(
				zap.Int("status_code", c.Writer.Status()),
				zap.Float64("latency_seconds", ginLogExtends.Latency),
			)
			// handler 中使用 c.Error(err) 后，会打印到 context_errors 字段中
			if len(c.Errors) > 0 {
				accessLogger = accessLogger.With(zap.String("context_errors", c.Errors.String()))
			}
			// 判断是否打印 context keys
			if conf.EnableContextKeys {
				accessLogger = accessLogger.With(zap.Any("context_keys", c.Keys))
			}
			// 判断是否打印响应 body
			if conf.EnableResponseBody {
				accessLogger = accessLogger.With(zap.Any("response_body", rspBodyWriter.body.String()))
			}
			detailFields := []zap.Field{
				zap.String("query", c.Request.URL.RawQuery),
				zap.String("proto", c.Request.Proto),
				zap.Int("content_length", int(c.Request.ContentLength)),
				zap.String("remote_addr", c.Request.RemoteAddr),
				zap.String("request_uri", c.Request.RequestURI),
				zap.String("referer", c.Request.Referer()),
				zap.String("user_agent", c.Request.UserAgent()),
				zap.String("content_type", c.ContentType()),
				zap.Int("body_size", c.Writer.Size()),
			}
			if conf.EnableDetails {
				accessLogger = accessLogger.With(detailFields...)
			}
			log := accessLogger.Info
			// 打印访问日志，根据状态码确定日志打印级别
			if c.Writer.Status() >= http.StatusInternalServerError || len(c.Errors) > 0 {
				// 500+ 始终打印带 details 的 error 级别日志
				// 无视配置开关，打印全部能搜集的信息
				accessLogger = accessLogger.With(detailFields...)
				accessLogger = accessLogger.With(zap.Any("context_keys", c.Keys))
				accessLogger = accessLogger.With(zap.Any("request_header", c.Request.Header))
				accessLogger = accessLogger.With(zap.Any("request_form", c.Request.Form))
				accessLogger = accessLogger.With(zap.String("request_body", string(GetGinRequestBody(c))))
				accessLogger = accessLogger.With(zap.String("response_body", rspBodyWriter.body.String()))
				log = accessLogger.Error
			} else if c.Writer.Status() >= http.StatusBadRequest {
				// 400+ 默认使用 warn 级别
				log = accessLogger.Warn
			}

			// 慢请求使用 Warn 记录
			if ginLogExtends.Latency > conf.SlowThreshold.Seconds() {
				accessLogger.Warn(
					formatter(c, ginLogExtends)+" hit slow request.",
					zap.Float64("slow_threshold", conf.SlowThreshold.Seconds()),
				)
			} else {
				log(formatter(c, ginLogExtends))
			}

			// update prometheus info
			labels := []string{fmt.Sprint(c.Writer.Status()), c.Request.URL.Path, c.Request.Method}
			promGinReqCount.WithLabelValues(labels...).Inc()
			promGinReqLatency.WithLabelValues(labels...).Observe(ginLogExtends.Latency)
		}()

		c.Next()
	}
}

//
// skipLog
//  @Description: 判断是否需要跳过日志记录
//  @param path
//  @param SkipPaths
//  @param skipRegexps
//  @return bool
//
func skipLog(path string, SkipPaths []string, skipRegexps []*regexp.Regexp) bool {
	for _, skipPath := range SkipPaths {
		if skipPath == path {
			return true
		}
	}
	for _, p := range skipRegexps {
		if p.MatchString(path) {
			return true
		}
	}

	return false
}

//
// GetGinRequestBody
//  @Description: 获取请求 body
//  @param c
//  @return []byte
//
func GetGinRequestBody(c *gin.Context) []byte {
	// 获取请求 body
	var requestBody []byte
	if c.Request.Body != nil {
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			_ = c.Error(err)
		} else {
			requestBody = body
			// body 被 read 、 bind 之后会被置空，需要重置
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		}
	}
	return requestBody
}

// 用于记录响应 body
type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

//
// Write
//  @Description: 覆盖 ResponseWriter 接口的 Write 方法，将 body 保存到 responseBodyWriter.body 中
//  @receiver w
//  @param b
//  @return int
//  @return error
//
func (w responseBodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
