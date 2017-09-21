package processors

import (
	"os"
	"time"
	"io"
	"regexp"
	"../utils"
)

const REGEXP_MATCHER = "\\d{4,}\\-\\d{1,}\\-\\d{1,}\\ \\d{1,}\\:\\d{1,}\\:\\d{2,}?"
const TIME_FORMAT = "2006-01-02 15:04:05"
const TCP = "tcp"
const BUF_SIZE = 4096
const WINDOW_SIZE = 512
const RETRY_DURATION = time.Second * 2

func process(reader *os.File, startPos int64, end time.Time, handler func(buf []byte) bool) {
	reader.Seek(startPos, 0)

	timeMatcher(reader, handler, func(curTime time.Time) bool {
		if curTime.After(end) {
			return true
		}
		return false
	}, nil)
}

func timeMatcher(reader io.Reader, handler func(buf []byte) bool, timeMatchedCallback func(curTime time.Time) bool, timeNotMatchedCallback func()) (timeMatched time.Time, err error) {
	buf := make([]byte, WINDOW_SIZE)
	for {
		_, err = reader.Read(buf)
		if err == io.EOF {
			break
		}

		if handler != nil && !handler(buf) {
			break
		}

		r, _ := regexp.Compile(REGEXP_MATCHER)
		curTimeStr := r.FindString(string(buf))
		if curTimeStr != "" {
			curTime, _ := time.Parse(TIME_FORMAT, curTimeStr)
			if timeMatchedCallback != nil && timeMatchedCallback(curTime) {
				return curTime, err
			}
		} else {
			if timeNotMatchedCallback != nil {
				timeNotMatchedCallback()
			}
		}
	}

	return time.Now(), err
}

func seekTime(reader io.Reader, pos int64) (time.Time, int64, error) {
	timeMatched, err := timeMatcher(reader, nil, func(curTime time.Time) bool {
		return true
	}, func() {
		pos += WINDOW_SIZE
	})
	return timeMatched, pos, err
}

func fastLocateStartPos(reader *os.File, beginTime time.Time) (pos int64, err error) {
	fileStat, _ := reader.Stat()
	fileSize := fileStat.Size()
	high := fileSize
	low := int64(0)
	for low <= high {
		pos = (low + high) / 2
		reader.Seek(pos, 0)
		curTime, pos, err := seekTime(reader, pos)

		if err != nil {
			utils.Notice("Error: locate pos error:", err.Error())
			break
		}

		if (curTime.Before(beginTime)) {
			low = pos + 1
		} else {
			high = pos - 1
		}
		break
	}
	return pos, err
}
