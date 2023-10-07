package app

import (
	"fmt"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/gtcp"
	"strings"
)

func StartTCPServer() {
	err := gtcp.NewServer("127.0.0.1:"+g.Cfg().GetString("tcp.port"), func(conn *gtcp.Conn) {
		defer conn.Close()
		for {
			data, err := conn.Recv(-1)
			if len(data) > 0 {
				//转换成JSON
				mJson := gjson.New(string(data))
				//获取header.eventname
				mEventName := mJson.GetString("header.eventname")
				//获取payload字符串
				mPayload := mJson.GetString("payload")
				//拆分字符串，解析出processname字段和from_state字段
				mPayloadMap := ParsePayload(mPayload)
				fmt.Println(mEventName, mPayloadMap)
				//
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

func ParsePayload(AIn string) g.Map {
	mR := g.Map{}
	//先按照空格拆分
	mT1 := strings.Split(AIn, " ")
	for _, v := range mT1 {
		if strings.TrimSpace(v) != "" {
			//按照冒号拆分
			mT2 := strings.Split(v, ":")
			mKey := mT2[0]
			mValue := mT2[1]
			mR[mKey] = mValue
		}
	}
	return mR
}
