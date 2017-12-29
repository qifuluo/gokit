package gokit

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

/*
读取格式为：
姓名 年龄 住址
小明 20   龙泉
....
的文本文件
注释为//
tbFile := utile.NewTableFile("E:\\user.txt", "	")
for !tbFile.Eof() {
	fmt.Println(tbFile.ReadInt("id", 0), tbFile.ReadString("name", "默认"), tbFile.ReadFloat("age", 0.0))
	tbFile.Next()
}
*/
type TableFile struct {
	strFile string
	strFlag string
	tbHead  []string            //头
	listVal []map[string]string //值
	iCurRow int                 //当前读取行
}

func NewTableFile(strFile, strFlag string) *TableFile {
	if !FileExist(strFile) {
		fmt.Println(fmt.Sprintf("file %v not find.", strFile))
		return nil
	}

	file, err := os.Open(strFile)
	if nil != err {
		fmt.Println(err)
		return nil
	}

	defer file.Close()

	tbFile := &TableFile{}
	tbFile.iCurRow = 0
	tbFile.strFile = strFile
	tbFile.strFlag = strFlag

	rd := bufio.NewReader(file)
	for {
		line, _, err := rd.ReadLine()
		if nil != err {
			break
		}

		tbFile.parseIniLine(string(line))
	}

	return tbFile
}

//是否还有值
func (self *TableFile) Eof() bool {
	if self.iCurRow >= len(self.listVal) {
		return true
	}

	return false
}

//转转下一条记录
func (self *TableFile) Next() {
	self.iCurRow++
}

//返回第一条记录
func (self *TableFile) Reset() {
	self.iCurRow = 0
}

//读取值
func (self *TableFile) ReadInt(strKey string, iDef int64) int64 {
	strVal, bOk := self.readVal(strKey)
	if !bOk {
		return iDef
	}

	iVal, err := strconv.ParseInt(strVal, 10, 0)
	if nil != err {
		return iDef
	}

	return iVal
}

func (self *TableFile) ReadFloat(strKey string, fDef float64) float64 {
	strVal, bOk := self.readVal(strKey)
	if !bOk {
		return fDef
	}

	fVal, err := strconv.ParseFloat(strVal, 64)
	if nil != err {
		return fDef
	}

	return fVal
}

func (self *TableFile) ReadString(strKey, strDef string) string {
	strVal, bOk := self.readVal(strKey)
	if !bOk {
		return strDef
	}

	return strVal
}

func (self *TableFile) readVal(strKey string) (strVal string, bOk bool) {
	if self.Eof() {
		fmt.Println("read eof line in file:", self.strFile)
		return "", false
	}

	mapVal := self.listVal[self.iCurRow]
	strVal, ok := mapVal[strKey]
	if !ok {
		fmt.Println(fmt.Sprintf("not find  %v's value in file: %v", strKey, self.strFile))
		return "", false
	}

	return strVal, true
}

//判断是否为注释行 "//"
func (self *TableFile) checkNotes(strLine string) bool {
	if 2 > len(strLine) {
		return false
	}

	cLine := []byte(strLine)
	if '/' == cLine[0] && '/' == cLine[1] {
		return true
	}

	return false
}

func (self *TableFile) parseIniLine(strLine string) {
	strLine = strings.TrimSpace(strLine)
	if 0 == len(strLine) {
		return
	}

	//是否为注释
	if self.checkNotes(strLine) {
		return
	}

	acLine := strings.Split(strLine, self.strFlag)
	//读取头
	if 0 == len(self.tbHead) {
		//检查空值,取得头
		for i := 0; i < len(acLine); i++ {
			if 0 == len(acLine[i]) {
				panic(fmt.Errorf("invalid table head in file: %v", self.strFile))
			}
			//是否有重复的
			for _, val := range self.tbHead {
				if val == acLine[i] {
					panic(fmt.Errorf("repeat head name %v in file: %v", val, self.strFile))
				}
			}

			self.tbHead = append(self.tbHead, acLine[i])
		}

		return
	}

	//读取值
	mapVal := make(map[string]string)
	for key, val := range self.tbHead {
		if key >= len(acLine) {
			mapVal[val] = ""
		} else {
			mapVal[val] = acLine[key]
		}
	}

	self.listVal = append(self.listVal, mapVal)
}
