package utils

import (
	"flag"
	"time"
	"os"
	"fmt"
	"strings"
	"path/filepath"
	"errors"
)

const DEFAULT_TIME_FORMAT = "2006-01-02 15:04:05"

type Flag struct {
	env string
	defaultValue string
	desc  string
}

var FlagMap = map[string]Flag{
	"begin": {
		"",
		time.Unix(0, 0).Format(DEFAULT_TIME_FORMAT),
		"start time",
	},
	"end": {
		"",
		time.Now().Format(DEFAULT_TIME_FORMAT),
		"end time",
	},
	"host": {
		"LL_HOST",
		"localhost:9000",
		"host",
	},
	"timeMatcher": {
		"LL_TIME_MATCHER",
		"",
		"matched word",
	},
	"path": {
		"LL_LOG_PATH",
		".",
		"Notice path",
	},
	"service": {
		"LL_SERVICE_NAME",
		"",
		"service name",
	},
	"daemon": {
		"",
		"",
		"daemonize",
	},
}

func getEnvOrDefaultValue(key string) string {
	if FlagMap[key].env != "" {
		return os.Getenv(FlagMap[key].env)
	} else {
		return FlagMap[key].defaultValue
	}
}

func WalkDir(dirPth, suffix string, process func(filename string) bool) error {
	suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写
	err := filepath.Walk(dirPth, func(filename string, fi os.FileInfo, err error) error { //遍历目录
		if err != nil { //忽略错误
		 return err
		}
		if fi.IsDir() { // 忽略目录
			return nil
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) {
			if !process(filename) {
				return errors.New("whoops")
			}
		}
		return nil
	})
	return err
}

func InitFlags() map[string]*string {
	var context map[string] *string
	context = make(map[string] *string)
	for name, v := range FlagMap {
		context[name] = flag.String(name, getEnvOrDefaultValue(name), v.desc)
	}
	flag.Parse()

	for name := range context {
		Notice("Config:", name, "(", *context[name], ")")
	}
	return context
}

func InitFile(filePath string) *os.File {
	reader, err := os.Open(filePath)

	if err != nil {
		Notice("ERROR: No such file:", filePath);
		os.Exit(0)
	}

	return reader
}

func Notice(arg ...string) {
	fmt.Println(arg)
}