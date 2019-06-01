package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

var configPath = "config.txt"
var cfstr, _ = ReadConfig(configPath)
var arrstr = strings.Split(cfstr, "\n") //config src arr

var fromFolder = strings.Replace(arrstr[1], "\\r", "", -1)
var outFolder = strings.Replace(arrstr[3], "\\r", "", -1)
var copyList = strings.Replace(arrstr[5], "\\r", "", -1)
var regString = strings.Replace(arrstr[7], "\\r", "", -1)

var restr = strings.Split(strings.Replace(arrstr[9], "\\r", "", -1), ",")
var outstr = strings.Split(strings.Replace(arrstr[11], "\\r", "", -1), ",")

func main() {
	defer func() {
		if eEnd := recover(); eEnd != nil {
			fmt.Println(eEnd)
			Log("--" + eEnd.(string))
		}
	}()
	if copyList == "" {
		e1 := copyDirFile(fromFolder, outFolder, restr, outstr)
		if e1 != nil {
			Log("error" + e1.Error())
		}
		Log("******log********")
	} else {
		e2 := copyFile(fromFolder, outFolder, copyList, restr, outstr)
		if e2 != nil {
			Log("error:" + e2.Error())
		}
		Log("******log********")
	}
}

//ReadConfig 读取配置
func ReadConfig(ru string) (string, error) {
	data, err := ioutil.ReadFile(ru)
	if err != nil {
		return "", err
	}
	aString := string(data)

	return aString, err
}

//---------------end config------------------------
func doReString(data []byte, restr []string, outstr []string) []byte {

	dataS := string(data)
	for i := 0; i < len(restr); i++ {
		dataS = strings.Replace(dataS, restr[i], outstr[i], -1)
		Log("-------文件的字符串{" + restr[i] + "}被修改为{" + outstr[i] + "}")
	}
	return []byte(dataS)
}

//复制某文件夹下所有的文件至新文件夹下
func copyDirFile(fromPath string, toPath string, restr []string, outstr []string) error {
	infor, e := ioutil.ReadDir(fromPath)
	if e != nil {
		return e
	}
	Log("--开始读取(" + fromPath + ")下的文件")
	for _, info := range infor {
		phNext := path.Join(fromPath, info.Name())
		thNext := path.Join(toPath, info.Name())
		if info.IsDir() { //is dir
			//检查toPath中是否存在这个dir thNext，不存在则创建
			if !Exists(thNext) {
				e1 := os.MkdirAll(thNext, 0666)
				if e1 != nil {
					return e1
				}
				Log("--" + thNext + "路径不存在，创建该路径文件")
			}
			//复制文件夹下的内容至thNext
			e2 := copyDirFile(phNext, thNext, restr, outstr)
			if e2 != nil {
				return e2
			}
			Log("--已复制文件夹 " + phNext + " 的内容至 " + thNext)
		} else { //is file
			rf, e3 := os.Open(phNext)
			defer rf.Close()
			if e3 != nil {
				return e3
			}
			Log("--开始读取文件(" + fromPath + ")--")

			data, e4 := ioutil.ReadAll(rf)
			if e4 != nil {
				return e4
			}
			if len(restr) != 0 && len(outstr) != 0 {
				data = doReString(data, restr, outstr)
			}
			e5 := ioutil.WriteFile(thNext, data, 0666)
			if e5 != nil {
				return e5
			}
			Log("--已复制文件 " + phNext + " 至 " + thNext)

		}
	}

	return nil
}

//复制文件至新文件夹
func copyFile(fromFolder string, outFolder string, copyList string, restr []string, outstr []string) error {
	//do read file and than copy and restr
	fileNameMsg, e1 := ReadConfig(copyList)
	if e1 != nil {
		return e1
	}
	Log("--开始读取(" + copyList + ")下的内容")
	fileNameList := strings.Split(fileNameMsg, "\n")

	alen := len(fileNameList)
	for i := 0; i < alen; i++ {

		reg := regexp.MustCompile(regString)
		fileLongName := reg.FindString(fileNameList[i])
		Log("------正则匹配的字符串为-" + fileLongName + "------")
		// strings.Replace(arrstr[5], "\\r", "", -1)

		// fromPath := fromFolder + fileLongName
		fromPath := path.Join(fromFolder, fileLongName)
		outPath := path.Join(outFolder, fileLongName)
		folder, _ := filepath.Split(outPath)
		if !Exists(folder) {
			ex := os.MkdirAll(folder, 0666)
			if ex != nil {
				return ex
			}
			Log("--" + folder + "路径不存在，创建该路径文件")

		}
		Log("--开始读取文件(" + fromPath + ")--")
		data, e2 := ioutil.ReadFile(fromPath)
		if e2 != nil {
			return e2
		}
		if len(restr) != 0 && len(outstr) != 0 {
			data = doReString(data, restr, outstr)
		}
		e3 := ioutil.WriteFile(outPath, data, 0666)
		if e3 != nil {
			return e3
		}
		Log("--已复制文件 " + fromPath + " 至 " + outPath)

	}
	return nil
}

//Exists 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

//Log 打印log至config
func Log(s string) {
	logf := configPath
	of, e := os.OpenFile(logf, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	defer of.Close()
	if e != nil {
		fmt.Println("********************************")
		fmt.Println("日志打印失败")
		fmt.Println("失败原因：" + e.Error())
		fmt.Println("失败内容：" + s)
		fmt.Println("********************************")
	}
	log.SetOutput(of)
	log.SetFlags(log.Ldate | log.Ltime)
	log.Println(s)
}
