package main

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	_ "github.com/lib/pq"
	"goJLinkClient/app"
)

func main() {
	//启动时检查supervisor是否开启RPC，如没有则自动修改设置并重启sp
	if app.SpInetOpened() == false {
		er := app.SpOpenRPC()
		if er == nil {
			//重启supervisord
			mOut, _ := app.DoShell("/etc/init.d/supervisor restart") //ubuntu 20.04
			fmt.Println(mOut)
		}
	}
	//启动监听Socket Server，接收Listener发过来的实时消息
	go app.StartTCPServer()
	//启动http服务
	s := g.Server()
	GroupApi := s.Group("/api")
	GroupApi.ALL("/cmd", app.SvdProce) //通用的单一报文接口处理
	GroupApi.ALL("/test0", app.Test0)  //简单的文本信息测试，已废弃
	GroupApi.ALL("/test1", app.Test1)  //测试GetAllProcessInfo
	GroupApi.ALL("/test2", app.Test2)  //测试DisplayConf
	GroupApi.ALL("/test3", app.Test3)  //测试InsertConf
	GroupApi.ALL("/test4", app.Test4)  //测试UpdateConf
	GroupApi.ALL("/test5", app.Test5)  //测试RemoveConf
	GroupApi.ALL("/test6", app.Test6)  //测试GetProcessInfo
	GroupApi.ALL("/test7", app.Test7)  //测试DisplayStdoutLog
	GroupApi.ALL("/test8", app.Test8)  //测试DisplaySysConfig
	GroupApi.ALL("/test9", app.Test9)  //测试DisplaySysInfo

	s.Run()
}
