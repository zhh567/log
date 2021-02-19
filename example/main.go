package main

import (
	"context"
	"fmt"
	"github.com/zhh567/log"
	"os"
	"time"
)

func main() {
	ctx, cancelLog := context.WithCancel(context.Background())
	logger := log.NewLoggerSync(ctx, log.DEBUG)

	//logger := log.NewLogger(log.INFO)

	//logger.SetFlag(log.Time|log.Line|log.FuncName|log.FileName|log.DirName)
	logger.SetJsonFormat()
	logger.SetPrefix("日志库中的测试代码")
	logger.SetMaxSize(1024)
	f, err := os.OpenFile("F:/code/golang/log/example/test.log", os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("open file failed, err: ", err)
		return
	}
	defer f.Close()
	logger.SetOutput(f)

	for {
		logger.Debug("日志")
		logger.Info("日志")
		logger.Warn("日志")
		logger.Error("日志")
		logger.FATAL("日志")
		time.Sleep(time.Second * 2)
	}


	cancelLog()
	time.Sleep(time.Second * 3)
}

