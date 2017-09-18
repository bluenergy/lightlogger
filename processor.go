package main

import (
	"os"
	"time"
	"net"
	"io"
	"regexp"
	"./protos"
	"github.com/golang/protobuf/proto"
)

const REGEXP_MATCHER = "\\d{4,}\\-\\d{1,}\\-\\d{1,}\\ \\d{1,}\\:\\d{1,}\\:\\d{2,}?"
const TIME_FORMAT = "2006-01-02 15:04:05"
const TCP = "tcp"
const BUF_SIZE = 4096
const WINDOW_SIZE = 512
const RETRY_DURATION = time.Second * 2

func StandaloneMod(reader *os.File, begin time.Time, end time.Time) {
	startPos := fastLocateStartPos(reader, begin)
	process(reader, startPos, end, nil)
}

func process(reader *os.File, startPos int64, end time.Time, handler func(buf []byte) (n int, err error)) {
	reader.Seek(startPos, 0)

	timeMatcher(reader, handler, func(curTime time.Time) bool {
		if curTime.After(end) {
			return true
		}
		return false
	}, nil)
}

func getConn(serviceName string, host string, reader *os.File) {
	for {
		conn, err := net.Dial(TCP, host)
		if err != nil {
			notice("ERROR: can not connect to host", err.Error())
			timer := time.NewTimer(1 * time.Second)
			<-timer.C
			continue
		} else {
			conn.Write([]byte(serviceName))
		}

		notice("Successfully connected to the host:", host)
		defer conn.Close()

		receiver(conn, reader)
	}
}

func daemonMod(reader *os.File, host string, serviceName string) {
	if serviceName == "" {
		notice("ERROR: service name is required")
		os.Exit(0)
	}

	getConn(serviceName, host, reader)
}

func sender(conn net.Conn, reader *os.File, cmd lightlog.Cmd) {
	start := time.Unix(cmd.GetStartTime(), 0)
	end := time.Unix(cmd.GetEndTime(), 0)
	startPos := fastLocateStartPos(reader, start)
	process(reader, startPos, end, conn.Write)
}

func receiver(conn net.Conn, reader *os.File) {
	recBuf := make([]byte, BUF_SIZE)

	for {
		conn.SetDeadline(time.Now().Add(RETRY_DURATION))
		recSize, err := conn.Read(recBuf)
		if err != nil {
			notice("ERROR: read buf error, trying to reconnect", err.Error())
			break
		}
		cmd := &lightlog.Cmd{}
		err = proto.Unmarshal(recBuf[:recSize], cmd)

		if err != nil {
			notice("ERROR: read cmd error:", err.Error())
			continue
		}

		go sender(conn, reader, *cmd)
	}
}

func timeMatcher(reader io.Reader, handler func(buf []byte) (int, error), timeMatchedCallback func(curTime time.Time) bool, timeNotMatchedCallback func()) time.Time {
	buf := make([]byte, WINDOW_SIZE)
	for {
		count, err := reader.Read(buf)
		if err == io.EOF && count == 0 {
			break
		}

		if handler != nil {
			handler(buf)
		}

		r, _ := regexp.Compile(REGEXP_MATCHER)
		curTimeStr := r.FindString(string(buf))
		if curTimeStr != "" {
			curTime, _ := time.Parse(TIME_FORMAT, curTimeStr)
			if timeMatchedCallback != nil && timeMatchedCallback(curTime) {
				return curTime
			}
		} else {
			if timeNotMatchedCallback != nil {
				timeNotMatchedCallback()
			}
		}
	}

	return time.Now()
}

func seekTime(reader io.Reader, pos int64) (time.Time, int64) {
	timeMatched := time.Now()
	timeMatched = timeMatcher(reader, nil, func(curTime time.Time) bool {
		return true
	}, func() {
		pos += WINDOW_SIZE
	})
	return timeMatched, pos
}

func fastLocateStartPos(reader *os.File, beginTime time.Time) int64 {
	fileStat, _ := reader.Stat()
	fileSize := fileStat.Size()
	high := fileSize
	low := int64(0)
	var pos int64
	for low <= high {
		pos = (low + high) / 2
		reader.Seek(pos, 0)
		curTime, pos := seekTime(reader, pos)
		if (curTime.Before(beginTime)) {
			low = pos + 1
		} else {
			high = pos - 1
		}
	}
	return pos
}
