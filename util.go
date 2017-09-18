package main

import (
	"flag"
	"time"
	"os"
	"fmt"
)

const DEFAULT_TIME_FORMAT = "2006-01-02 15:04:05"

var LOG_PATH = os.Getenv("LL_LOG_PATH")
var HOST = os.Getenv("LL_HOST")

type Flag struct {
	value string
	desc  string
}

var FlagMap = map[string]Flag{
	"begin": {
		time.Unix(0, 0).Format(DEFAULT_TIME_FORMAT),
		"start time",
	},
	"end": {
		time.Now().Format(DEFAULT_TIME_FORMAT),
		"end time",
	},
	"host": {
		"",
		"host",
	},
	"timeMatcher": {
		"",
		"matched word",
	},
	"path": {
		"",
		"notice path",
	},
	"service": {
		"",
		"service name",
	},
}

func InitFlags() map[string]*string {
	var context map[string] *string
	context = make(map[string] *string)
	for name, v := range FlagMap {
		context[name] = flag.String(name, v.value, v.desc)
	}
	flag.Parse()
	return context
}

func InitFile(filePath string) *os.File {
	reader, err := os.Open(filePath)

	if err != nil {
		notice("ERROR: No such file:", filePath);
		os.Exit(0)
	}

	return reader
}

func notice(arg ...string) {
	fmt.Println(arg)
}