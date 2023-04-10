package main

import (
	"context"

	"github.com/feymanlee/logit"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()
	/* zap Debug */
	logit.Debug(ctx, "Debug message", zap.Int("intType", 123), zap.Bool("boolType", false), zap.Ints("sliceInt", []int{1, 2, 3}), zap.Reflect("map", map[string]interface{}{"i": 1, "s": "s"}))
	// Output:
	// {"level":"DEBUG","time":"2020-04-15 18:12:11.991006","logger":"logit.ctx_logger","msg":"Debug message","pid":45713,"intType":123,"boolType":false,"sliceInt":[1,2,3],"map":{"i":1,"s":"s"}}

	/* zap sugared logger Debug */
	logit.Debugs(ctx, "Debugs message", 123, false, []int{1, 2, 3}, map[string]interface{}{"i": 1, "s": "s"})
	// Output:
	// {"level":"DEBUG","time":"2020-04-15 18:12:11.991239","logger":"logit.ctx_logger","msg":"Debugs message123 false [1 2 3] map[i:1 s:s]","pid":45713}

	/* zap sugared logger Debugf */
	logit.Debugf(ctx, "Debugf message, %s", "ok")
	// Output:
	// {"level":"DEBUG","time":"2020-04-15 18:12:11.991268","logger":"logit.ctx_logger","msg":"Debugf message, ok","pid":45713}

	/* zap sugared logger Debugw */
	logit.Debugw(ctx, "Debugw message", "name", "axiaoxin", "age", 18)
	// Output:
	// {"level":"DEBUG","time":"2020-04-15 18:12:11.991277","logger":"logit.ctx_logger","msg":"Debugw message","pid":45713,"name":"axiaoxin","age":18}

	/* with context */
	c, _ := logit.NewCtxLogger(context.Background(), logit.CloneLogger("myname"), "trace-id-123")
	logit.Debug(c, "Debug with trace id")
	// Output:
	// {"level":"DEBUG","time":"2020-04-15 18:12:11.991314","logger":"logit.myname","msg":"Debug with trace id","pid":45713,"traceID":"trace-id-123"}

	/* extra fields */
	logit.Debug(c, "extra fields demo", logit.ExtraField("k1", "v1", "k2", 2, "k3", true))
	// Output:
	// {"level":"DEBUG","time":"2020-04-15 18:12:11.991348","logger":"logit.myname","msg":"extra fields demo","pid":45713,"traceID":"trace-id-123","extra":{"k1":"v1","k2":2,"k3":true}}
}
