package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//super visor
func SuperVisor(cmd string,blankTime int32) {
	//先查找这个进程是否存在 不设置退出的条件
	//处理掉命令中多余的空格
	newCmd := DeleteExtraSpace(cmd)
	pid := 0
	for{
		for {
			if cmd != newCmd{
				pid = GetServicePid(cmd)
			}else{
				pid = GetServicePid(newCmd)
			}
			if 0 == pid{
				break
			}
			//	fmt.Println(cmd," exist ",pid)
			time.Sleep(time.Duration(time.Second * time.Duration(blankTime)))
		}
		//进程不存在了 执行启动命令 启动不需要等待
		fmt.Println("do cmd ",newCmd)
		err := DoLinuxCmd(newCmd)
		if nil != err{
			fmt.Println("do cmd fail ",newCmd,err)
		}
		//	fmt.Println("start cmd ",cmd,err)
		//	time.Sleep(time.Duration(time.Second * time.Duration(blankTime)))
	}
}
//do cmd
func DoLinuxCmd(cmd string)error{
	args := strings.Split(cmd," ")
	execCmd := exec.Command(args[0],args[1:]...)
	return execCmd.Run()
}
//get app info by ps
func GetServicePsInfo(serviceCmd string) string{
	ps := exec.Command("ps","-aux")
	grep := exec.Command("grep",serviceCmd)
	p,g := io.Pipe()
	defer p.Close()
	defer g.Close()
	ps.Stdout = g
	grep.Stdin = p
	var buffer bytes.Buffer
	grep.Stdout = &buffer
	_ = ps.Start()
	_ = grep.Start()
	io.Copy(os.Stdout,&buffer)
	ps.Wait()
	g.Close()
	grep.Wait()
	return DeleteExtraSpace(buffer.String())
}

//get the process pid by app name or start cmd
func GetServicePid(serviceCmd string)int{
	res := GetServicePsInfo(serviceCmd)
	if "" == res{
		return 0
	}
	resSplit := strings.Split(res," ")
	if len(resSplit) < 2{
		return 0
	}
	val,_ := strconv.Atoi(resSplit[1])
	return val
}

//delete the space where bigger than 2
func DeleteExtraSpace(src string) string {
	//删除字符串中的多余空格，有多个空格时，仅保留一个空格
	newSrc := strings.Replace(src, "  ", " ", -1)      //替换tab为空格
	regStr := "\\s{2,}"                          //两个及两个以上空格的正则表达式
	reg, _ := regexp.Compile(regStr)             //编译正则表达式
	res := make([]byte, len(newSrc))                  //定义字符数组切片
	copy(res, newSrc)                                 //将字符串复制到切片
	spcIndex := reg.FindStringIndex(string(res)) //在字符串中搜索
	for len(spcIndex) > 0 {                     //找到适配项
		res = append(res[:spcIndex[0]+1], res[spcIndex[1]:]...) //删除多余空格
		spcIndex = reg.FindStringIndex(string(res))            //继续在字符串中搜索
	}
	return string(res)
}

//kill process by app name or start cmd
func KillProcess(serviceCmd string){
	err := DoLinuxCmd(fmt.Sprintf("kill -9 %d",GetServicePid(serviceCmd)))
	fmt.Println(err)
}

func main() {
	var serviceCmd string
	flag.StringVar(&serviceCmd,"c","","")
	flag.Parse()
	if serviceCmd == ""{
		panic("empty")
	}
	serviceCmd += " ../winwinlog"
	fmt.Println(GetServicePid(serviceCmd))
	SuperVisor(serviceCmd,3)
}
