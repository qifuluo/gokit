//日志
package gokit

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strconv"
	"sync/atomic"
	"time"
)

//日志结构体
type STLog struct {
	quit        int32
	logChan     chan string
	logFile     *os.File
	logPriority int
	filePath    string
	fileName    string
	curDate     string
}

//新建日志对象
func NewLogger(logPriority int, strName string) *STLog {
	newLogger := &STLog{}

	newLogger.fileName = strName
	newLogger.filePath = GetProPath() + "log/"
	newLogger.logPriority = logPriority
	newLogger.logChan = make(chan string, CHANSIZE_LOG)
	newLogger.quit = 0
	newLogger.curDate = newLogger.getCurDate()

	newLogger.open()
	fmt.Println("log file path:", newLogger.filePath)
	fmt.Println("log file name:", newLogger.fileName)
	fmt.Println("log priority:", newLogger.logPriority)

	go newLogger.loop()

	return newLogger
}

func (self *STLog) open() {
	var err error
	self.logFile, err = os.OpenFile(self.filePath+self.fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if nil != err {
		panic(err)
	}
}

func (self *STLog) renameLog() {
	curDate := self.getCurDate()
	if curDate == self.curDate {
		return
	}

	self.logFile.Close()

	iCount := 0
	newName := self.filePath + self.curDate + ".txt"
	for {
		if !FileExist(newName) {
			break
		}

		iCount++
		newName = self.filePath + self.curDate + "(" + strconv.Itoa(iCount) + ").txt"
	}

	os.Rename(self.filePath+self.fileName, newName)
	self.open()
	self.curDate = curDate
}

//写
func (self *STLog) loop() {
	defer self.logFile.Close()

	for 0 == atomic.LoadInt32(&self.quit) {
		select {
		case strLog := <-self.logChan:
			self.renameLog()

			strLog = fmt.Sprintf("[%v]%v", time.Now().Format("2006-01-02 15:04:05"), strLog)
			if LV_DEBUG == self.logPriority {
				fmt.Print(strLog)
			}
			self.logFile.Write([]byte(strLog))
		}
	}
}

func (self *STLog) getCurDate() string {
	nowTime := time.Now()
	return fmt.Sprintf("%d-%d-%d", nowTime.Year(), nowTime.Month(), nowTime.Day())
}

func (self *STLog) getPriority() int {
	return self.logPriority
}

//获取对应的字符串
func (self *STLog) getLogLVStr(logLV int) string {
	var strLV string

	switch logLV {
	case LV_SYS:
		strLV = "SYSTEM"
	case LV_DEBUG:
		strLV = "DEBUG"
	case LV_INFO:
		strLV = "INFO"
	case LV_WARN:
		strLV = "WARNING"
	case LV_ERROR:
		strLV = "ERROR"
	default:
		strLV = "UNKNOWN"
	}

	return strLV
}

//日志信息格式化
func (self *STLog) Log(logLV int, v ...interface{}) {
	if 0 != atomic.LoadInt32(&self.quit) {
		return
	}

	_, file, line, ok := runtime.Caller(2)
	if !ok {
		fmt.Println("get log position error.")
	}

	strMsg := fmt.Sprint(v...)
	strLogInfo := fmt.Sprintf("[%v][%v %v]%v\n",
		self.getLogLVStr(logLV), path.Base(file), line, strMsg)

	self.logChan <- strLogInfo
}

//关闭退出
func (self *STLog) Close() {
	if 0 != atomic.LoadInt32(&self.quit) {
		return
	}

	atomic.AddInt32(&self.quit, 1)
}
