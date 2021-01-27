package main

import (
	"cmall/conf"
	"cmall/model"
	"cmall/routes"
)

func main() {
	//从配置文件读入配置
	conf.Init()
	model.ListenOrder()
	//转载路由
	r := routes.NewRouter()
	_ = r.Run(":3000")
}
