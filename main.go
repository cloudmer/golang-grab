package main

import (
	"runtime"
	"time"
	"xmn/core/algorithm/ssc"
	_ "xmn/core/model"
	"log"
	"os"
	"xmn/core/algorithm/shishicai/CustomPackage"
	"xmn/core/algorithm/shishicai/play1"
	"xmn/core/algorithm/shishicai/play22"
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
			// 时时彩 连号 算法 邮件报警
			go ssc.Consecutive()
			// 时时彩 连续AB表报警
			go ssc.ContailMultiple()
			// 时时彩 AB包 自定义A包周期 报警
			go CustomPackage.Calculation()
			// 时时彩 a出现几期的b
			go play1.Calculation()
			// 时时彩 2连1站报警
			go play22.Consecutive()
		}
	}
}
