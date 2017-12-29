package gokit

import (
	"runtime"
	"fmt"
	"path/filepath"
	"os"
	"strings"
	"time"
	"math"
)

//产生panic后调用栈打印
func PanicStack(extras ...interface{}) string {
	var str = ""
	if x := recover(); x != nil {
		i := 0
		funcName, file, line, ok := runtime.Caller(i)
		for ok {
			str += fmt.Sprintf("frame %v:[func:%v,file:%v,line:%v]\n", i, runtime.FuncForPC(funcName).Name(), file, line)
			i++
			funcName, file, line, ok = runtime.Caller(i)
		}
		for k := range extras {
			str += fmt.Sprintf("EXRAS#%v DATA:type:%T, val:%+v\n", k, extras[k], extras[k])
		}
	}

	return str
}

//调用堆栈
func Stack() string {
	stack := make([]byte, 4096)
	runtime.Stack(stack, false)

	return string(stack)
}

//获取当前程序所在路径
func GetProPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if nil != err {
		return ""
	}

	return strings.Replace(dir, "\\", "/", -1) + "/"
}

//文件是否存在
func FileExist(strFile string) bool {
	var bExist = true
	if _, err := os.Stat(strFile); os.IsNotExist(err) {
		bExist = false
	}

	return bExist
}

func AbsInt64(x int64) int64 {
	if x >= 0 {
		return x
	} else {
		return -1 * x
	}
}

//strTime 2017-10-16 00:00:00
func StrToTime(strTime string) int64 {
	loc, _ := time.LoadLocation("Local")
	tm, _ := time.ParseInLocation("2006-01-02 15:04:05", strTime, loc)

	return tm.Unix()
}

//当前毫秒数
func Millisecond() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

//获取一个月多少天
var days = [12]int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

func GetDayInMon(year int, mon int) int {
	var day int
	if 2 == mon {
		if (year%4 == 0 && year%100 != 0) || year%400 == 0 {
			day = 29
		} else {
			day = 28
		}
	} else {
		day = days[mon-1]
	}

	return day
}

//检查是否包含NaN Inf
func IsNaN(params ...float64) bool {
	for _, param := range params {
		if math.IsNaN(param) || math.IsInf(param, 0) {
			return true
		}
	}

	return false
}


