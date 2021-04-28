package lib

import (
	"bytes"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"strings"
)

var (
	ConfEnvPath string //配置文件夹
	ConfEnv     string //配置环境名 比如：dev prod test
)

// 解析配置文件目录
//
// 配置文件必须放到一个文件夹中
// 如：config=conf/dev/base.json 	ConfEnvPath=conf/dev	ConfEnv=dev
// 如：config=conf/base.json		ConfEnvPath=conf		ConfEnv=conf
func ParseConfPath(config string) error {
	path := strings.Split(config, "/")
	prefix := strings.Join(path[:len(path)-1], "/")
	ConfEnvPath = prefix
	ConfEnv = path[len(path)-2]

	return nil
}

func GetConfEnv() string {
	return ConfEnv
}

func GetConfPath(filename string) string {
	return ConfEnvPath + "/" + ".toml" //解析后缀名为 .toml的文件
}

func GetConfFilePath(fileName string) string {
	return ConfEnvPath + "/" + fileName
}

//本地文件解析
func ParseLocalConfig(fileName string, st interface{}) error {
	path := GetConfFilePath(fileName)

	err := ParseConfig(path, st)
	if err != nil {
		return err
	}

	return nil
}

func ParseConfig(path string, conf interface{}) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("Open config %v fail, %v", path, err)
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("Read config fail, %v", err)
	}

	v := viper.New()
	v.SetConfigType("toml") //解析文件类型
	_ = v.ReadConfig(bytes.NewBuffer(data))

	if err := v.Unmarshal(conf); err != nil {
		return fmt.Errorf("Parse config fail, config:%v, err:%v", string(data), err)
	}

	return nil
}
