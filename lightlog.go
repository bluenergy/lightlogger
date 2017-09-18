package main

import "time"

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
	 		b. 定位行
			c. 输出日志
			d. stop
			e. multi client
			f. tail
			g. multi log file
			h. daemon
	*/

	var context = InitFlags()
	filePath := *context["path"]

	reader := InitFile(filePath)
	defer reader.Close()

	host := *context["host"]
	serviceName := *context["service"]

	if host != "" {
		daemonMod(reader, host, serviceName)
	} else {
		begin, _ := time.Parse(TIME_FORMAT, *context["begin"])
		end, _ := time.Parse(TIME_FORMAT, *context["end"])
		StandaloneMod(reader, begin, end)
	}

}
