package main

import (
	"runtime"
	"time"
	"xmn/core/algorithm/ssc"
	_ "xmn/core/model"
	"log"
	"os"
)

func main(){
	log.Println("服务启动中．．．　进程ID:", os.Getpid())
	runtime.GOMAXPROCS(runtime.NumCPU())
	for {
		select {
		case <-time.After(1 * time.Minute):
			//todo
			// 时时彩　包含数据包　算法　邮件报警
			go ssc.Contain()
		}
	}
}
