package main

import (
	"github.com/qianyaozu/qgenerate/mssql"
	"github.com/qianyaozu/qconf"
	"fmt"
)

func main(){
	LoadConf()
}

func LoadConf() {
	conf, err := qconf.LoadConfiguration("conf.ini")
	if err != nil {
		panic(err)
	}
	MSSqlConnection  := conf.GetString("MSSqlConnection")
 	if MSSqlConnection!=""{
		mssql.Generator(MSSqlConnection)
	}
	var command=""
	fmt.Print("生成完成，输入任意字符关闭")
	fmt.Scanln(&command)
}

