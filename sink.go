// 使用 zap.RegisterSink 函数和 Config.OutputPaths 字段添加自定义日志目标。
// RegisterSink 将 URL 方案映射到 Sink 构造函数， OutputPaths 配置日志目的地（编码为 URL ）。
// *lumberjack.Logger 已经实现了几乎所有的 zap.Sink 接口。只缺少 Sync 方法。

package logit

import (
	"net/url"

	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
)

// LumberjackSink 将日志输出到 lumberjack 进行 rotate
type LumberjackSink struct {
	*lumberjack.Logger
}

// Sync lumberjack Logger 默认已实现 Sink 的其他方法，这里实现 Sync 后就成为一个 Sink 对象
func (LumberjackSink) Sync() error {
	return nil
}

// RegisterSink 注册 lumberjack sink
// 在 OutputPaths 中指定输出为 sink.Scheme://log_filename 即可使用
// path url 中不指定日志文件名则使用默认的名称
// 一个 scheme 只能对应一个文件名，相同的 scheme 注册无效，会全部写入同一个文件
func RegisterSink(scheme string, sink zap.Sink) error {
	return zap.RegisterSink(scheme, func(*url.URL) (zap.Sink, error) {
		return sink, nil
	})
}

// NewLumberjackSink
//
//	@Description: 创建 LumberjackSink 对象
//	@param scheme sink scheme
//	@param filename 文件名称
//	@param maxAge 最大生命周期
//	@param maxBackups 最多保留文件个数
//	@param maxSize 单个文件最大 size
//	@param compress 是否压缩
//	@param localtime 是否采用本地时间
//	@return *LumberjackSink
func NewLumberjackSink(filename string, maxAge, maxBackups, maxSize int, compress, localtime bool) *LumberjackSink {
	return &LumberjackSink{
		Logger: &lumberjack.Logger{
			Filename:   filename,
			MaxAge:     maxAge,
			MaxBackups: maxBackups,
			MaxSize:    maxSize,
			Compress:   compress,
			LocalTime:  localtime,
		},
	}
}
