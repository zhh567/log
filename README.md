# log

This a simple log library.

It can be classified, output JSON format, and log file segmentation.

## Async mode
```go
package main

import (
	"context"
	"github.com/zhh567/log"
	"time"
)

func main() {
	ctx, cancelLog := context.WithCancel(context.Background())
	logger := log.NewLoggerSync(ctx, log.DEBUG)

	logger.SetJsonFormat()
	logger.SetPrefix("日志库中的测试代码")

	logger.Debug("日志")
	logger.Info("日志")
	logger.Warn("日志")
	logger.Error("日志")
	logger.FATAL("日志")

	cancelLog()
	time.Sleep(time.Second * 3)
}
```

## Sync mode
```go
package main

import (
	"github.com/zhh567/log"
)

func main() {
	logger := log.NewLogger(log.DEBUG)

	logger.SetFlag(log.Time|log.Line|log.FuncName|log.FileName|log.DirName)
	logger.SetJsonFormat()
	logger.SetPrefix("日志库中的测试代码")

	logger.Debug("日志")
	logger.Info("日志")
	logger.Warn("日志")
	logger.Error("日志")
	logger.FATAL("日志")
}
```
