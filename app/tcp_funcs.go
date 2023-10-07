package app

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/gtcp"
)

func StartTCPServer() {
	err := gtcp.NewServer("127.0.0.1:"+g.Cfg().GetString("tcp.port"), func(conn *gtcp.Conn) {
		defer conn.Close()
		for {
			data, err := conn.Recv(-1)
			if len(data) > 0 {
				fmt.Println("TCP Server Receive: ", string(data))
			}
			if err != nil {
				break
			}
		}
	}).Run()
	if err != nil {
		fmt.Println("TCP Server Error: ", err)
	}
}
