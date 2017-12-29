//ini 文件读取
package gokit

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

//ini结构体
type STIni struct {
	mapVal  map[string]map[string]string
	strFile string
}

//判断是否为节点 是则读取值
func checkNode(strLine string) (strNode string, bNode bool) {
	bNode = false
	strNode = ""

	if 2 >= len(strLine) {
		return
	}

	cLine := []byte(strLine)
	if '[' == cLine[0] && ']' == cLine[len(cLine)-1] {
		bNode = true
		strNode = string(cLine[1 : len(cLine)-1])
		return
	}

	return
}

//判断是否为注释行 "//"
func checkNotes(strLine string) bool {
	if 2 > len(strLine) {
		return false
	}

	cLine := []byte(strLine)
	if '/' == cLine[0] && '/' == cLine[1] {
		return true
	}

	return false
}

func parseIniLine(strLine string, strCurNode string, iniFile *STIni) string {
	strTmp := strings.TrimSpace(strLine)
	//空行
	if 0 == len(strTmp) {
		return strCurNode
	}

	//判断是否为注释
	if checkNotes(strTmp) {
		return strCurNode
	}

	//判断是否为节点
	strNode, bNode := checkNode(strTmp)
	if bNode {
		_, exists := iniFile.mapVal[strNode]
		if !exists {
			iniFile.mapVal[strNode] = make(map[string]string)
		}
		return strNode
	}

	//还没有节点
	if 0 == len(strCurNode) {
		return strCurNode
	}

	//读取值
	iPos := strings.Index(strTmp, "=")
	if 0 >= iPos {
		return strCurNode
	}

	rsTmp := []rune(strTmp)
	strKey := string(rsTmp[0:iPos])
	strKey = strings.TrimSpace(strKey)
	strVal := string(rsTmp[iPos+1 : len(strTmp)])
	strVal = strings.TrimSpace(strVal)

	iniFile.mapVal[strCurNode][strKey] = strVal

	return strCurNode
}

func NewIni(strFile string) *STIni {
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

	iniFile := &STIni{}
	iniFile.strFile = strFile
	iniFile.mapVal = make(map[string]map[string]string)

	var strCurNode string
	rd := bufio.NewReader(file)
	for {
		line, _, err := rd.ReadLine()
		if nil != err {
			break
		}

		strCurNode = parseIniLine(string(line), strCurNode, iniFile)
	}

	return iniFile
}

func (self *STIni) readVal(strNode, strKey string) (strVal string, bOk bool) {
	mapNode, exists := self.mapVal[strNode]
	if !exists {
		fmt.Println(fmt.Sprintf("not find node %v in file: %v", strNode, self.strFile))
		return "", false
	}

	strVal, exists = mapNode[strKey]
	if !exists {
		fmt.Println(fmt.Sprintf("not find key %v in file: %v", strKey, self.strFile))
		return "", false
	}

	return strVal, true
}

//读取整形类型值
func (self *STIni) ReadInt(strNode string, strKey string, iDef int64) int64 {
	strVal, bOk := self.readVal(strNode, strKey)
	if !bOk {
		return iDef
	}

	iVal, err := strconv.ParseInt(strVal, 10, 0)
	if nil != err {
		return iDef
	}

	return iVal
}

//读取浮点类型值
func (self *STIni) ReadFloat(strNode string, strKey string, fDef float64) float64 {
	strVal, bOk := self.readVal(strNode, strKey)
	if !bOk {
		return fDef
	}

	fVal, err := strconv.ParseFloat(strVal, 64)
	if nil != err {
		return fDef
	}

	return fVal
}

//读取字符串类型值
func (self *STIni) ReadString(strNode, strKey, strDef string) string {
	strVal, bOk := self.readVal(strNode, strKey)
	if !bOk {
		return strDef
	}

	return strVal
}
