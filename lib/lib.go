package lib

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// 只获取当前目录的
func GetAllFile(pwd string) ([]string, error) {
	//pwd, _ := os.Getwd()
	//fmt.Println(pwd)
	//获取文件或目录相关信息
	result := make([]string, 0)
	fileInfoList, err := ioutil.ReadDir(pwd)
	if err != nil {
		return []string{}, err
	}
	for i := range fileInfoList {
		//fmt.Println(fileInfoList[i].Name()) //打印当前文件或目录下的文件或目录名
		fileName := fileInfoList[i].Name()
		if !strings.HasPrefix(fileName, ".") && strings.HasSuffix(fileName, ".mod") {
			fmt.Println(path.Join(pwd, fileName))
			result = append(result, path.Join(pwd, fileName))
		}
	}
	return result, err
}

// 获取所有的
func RecursionGetAllFile(pwd string) ([]string, error) {
	result := make([]string, 0)
	err := filepath.Walk(pwd, func(dir string, info os.FileInfo, err error) error {
		fileName := info.Name()
		if !strings.HasPrefix(fileName, ".") && strings.HasSuffix(fileName, ".log") {
			result = append(result, dir)
		}
		return nil
	})
	return result, err
}

//get server ip

func GetServerIp() string {
	ip := ""
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ip
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ip = ipNet.IP.String()
				break
			}
		}
	}
	return ip
}
