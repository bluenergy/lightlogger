package main

import (
	"time"
	"./utils"
	"./processors"
)

func main() {

	/*TODO:
		1. 初始化
			a. 环境变量， 参数
			b. 检查环境变量， 参数， 决定启动方式
		2. 运行方式：
			a. standalone模式
			b. daemon模式
		3. 标准流程：
			a. 文件句柄， 错误处理
			d. stop - done
			e. multi client
			f. tail
			g. multi log file - done
			h. daemon
	*/

	var context = utils.InitFlags()
	dirPath := *context["path"]

	//dir := utils.InitFile(dirPath)
	//defer dir.Close()

	host := *context["host"]
	serviceName := *context["service"]

	if host != "" {
		processors.DaemonMod(dirPath, host, serviceName)
	} else {
		begin, _ := time.Parse(processors.TIME_FORMAT, *context["begin"])
		end, _ := time.Parse(processors.TIME_FORMAT, *context["end"])
		processors.StandaloneMod(dirPath, begin, end)
	}

}
