package logit

import (
	"runtime"
	"strings"
	"time"

	"go.uber.org/zap/zapcore"
)

//
// FuncName
//  @Description: FuncName 返回调用本函数的函数名称
//  @param pc runtime.Caller 返回的第一个值
//  @return string
//
func FuncName(pc uintptr) string {
	funcName := runtime.FuncForPC(pc).Name()
	sFuncName := strings.Split(funcName, ".")
	return sFuncName[len(sFuncName)-1]
}

//
// CallerEncoder
//  @Description: serializes a caller in package/file:funcname:line format
//  @param caller
//  @param enc
//
func CallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	shortCaller := caller.TrimmedPath()
	shortCallerSplited := strings.Split(shortCaller, ":")
	funcName := FuncName(caller.PC)
	result := shortCallerSplited[0] + ":" + funcName + ":" + shortCallerSplited[1]
	enc.AppendString(result)
}

//
// TimeEncoder
//  @Description: 自定义日志时间格式, 不带时区信息， YYYY-mm-dd H:M:S.xxxxxx
//  @param t
//  @param enc
//
func TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000000"))
}
