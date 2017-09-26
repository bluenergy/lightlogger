package processors

import (
	"os"
	"net"
	"time"
	"github.com/golang/protobuf/proto"
	"../utils"
	"../protos"
	"strings"
	"io"
)

const STOP = "STOP"
const UPDATE = "UPDATE"
const RESULT = "RESULT"
const SEARCH = "SEARCH"
const SUFFIX = ".log"

func DaemonMod(dir string, host string, serviceName string) {
	if serviceName == "" {
		utils.Notice("ERROR: service name is required")
		os.Exit(0)
	}

	getConn(serviceName, host, dir)
}

func getConn(serviceName string, host string, dir string) {
	for {
		conn, err := net.Dial(TCP, host)
		if err != nil {
			utils.Notice("ERROR: can not connect to host", err.Error())
			timer := time.NewTimer(1 * time.Second)
			<-timer.C
			continue
		} else {
			data, _ := commandBuilder(UPDATE, "", int64(0), int64(0), []byte(serviceName))
			conn.Write(data)
		}

		utils.Notice("Successfully connected to the host:", host)
		defer conn.Close()

		receiver(conn, dir)
	}
}

func receiver(conn net.Conn, dir string) {
	recBuf := make([]byte, BUF_SIZE)

	for {
		conn.SetDeadline(time.Now().Add(RETRY_DURATION))
		size, err := conn.Read(recBuf)
		if err != nil {
			readBufError := err.Error()
			utils.Notice(readBufError)
			if err == io.EOF || strings.Contains(readBufError, "connection reset by peer") {
				utils.Notice("ERROR: read buf readBufError, trying to reconnect", err.Error())
				break
			} else if strings.Contains(readBufError, "i/o timeout") {
				utils.Notice("INFO: read buf timeout")
				continue
			}
		}

		cmd, err := commandParser(size, recBuf)

		if err != nil {
			utils.Notice("ERROR: read cmd error:", err.Error())
			continue
		}
		commandHandler(cmd, conn, dir)
	}
}

func commandParser(size int, buf []byte) (lightlog.Cmd, error) {
	cmd := &lightlog.Cmd{}
	err := proto.Unmarshal(buf[:size], cmd)

	return *cmd, err
}

func commandBuilder(cmd string, id string, startTime int64, endTime int64, buf []byte) ([]byte, error) {
	command := &lightlog.Cmd{
		Cmd:       cmd,
		Id:       id,
		StartTime: startTime,
		EndTime:   endTime,
		Data:      buf,
	}
	var err error
	data, err := proto.Marshal(command)
	if err != nil {
		utils.Notice("marshaling error: ", err.Error())
		//os.Exit(0)
	}

	return data, err
}

func commandHandler(cmd lightlog.Cmd, conn net.Conn, dir string) {
	stopSigCH := make(chan string, 2)
	utils.Notice(cmd.Cmd)
	switch cmd.Cmd {
	case SEARCH:
		go sender(conn, dir, cmd, stopSigCH)
	case STOP:
		stopSigCH <- STOP
	}
}

func sender(conn net.Conn, dir string, cmd lightlog.Cmd, stopSigCH chan string) {
	utils.WalkDir(dir, SUFFIX, func(filename string) bool {
		reader := utils.InitFile(filename)
		defer reader.Close()

		start := time.Unix(cmd.GetStartTime(), 0)
		end := time.Unix(cmd.GetEndTime(), 0)

		_, err := fastLocateStartPos(reader, start)

		if err != nil {
			return false
		}

		process(reader, 0, end, func(buf []byte) bool {
			select {
			case <-stopSigCH:
				return false
			default:
				data, err := commandBuilder(RESULT, cmd.GetId(), int64(0), int64(0), buf)
				if err == nil {
					conn.Write(data)
				}
				return true
			}
		})

		return true
	})
}
