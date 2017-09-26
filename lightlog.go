package main

import (
	"time"
	"./utils"
	"./processors"
	"github.com/sevlyar/go-daemon"
	"log"
)

func main() {

	/*TODO:
		1. 初始化
			a. 环境变量， 参数 - done
			b. 检查环境变量， 参数， 决定启动方式 -done
		2. 运行方式：
			a. standalone模式
			b. daemon模式 - done
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
	host := *context["host"]
	serviceName := *context["service"]
	daemonMod := *context["daemon"]
	begin := *context["begin"]
	end := *context["end"]

	if daemonMod == "true" {
		runDaemon(host, dirPath, serviceName, begin, end)
	} else {
		run(host, dirPath, serviceName, begin, end)
	}
}

func run(host string, dirPath string, serviceName string, begin string, end string) {
	if host != "" {
		processors.DaemonMod(dirPath, host, serviceName)
	} else {
		begin, _ := time.Parse(processors.TIME_FORMAT, begin)
		end, _ := time.Parse(processors.TIME_FORMAT, end)
		processors.StandaloneMod(dirPath, begin, end)
	}

}

func runDaemon(host string, dirPath string, serviceName string, begin string, end string) {
	cntxt := &daemon.Context{
		PidFileName: "pid",
		PidFilePerm: 0644,
		LogFileName: "log",
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
		Args:        []string{"[lightlogger]"},
	}

	d, err := cntxt.Reborn()
	if err != nil {
		log.Fatal("Unable to run: ", err)
	}
	if d != nil {
		return
	}
	defer cntxt.Release()

	utils.Notice("daemon started")

	run(host, dirPath, serviceName, begin, end)
}
