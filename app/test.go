package app

import (
	"fmt"
	"github.com/gogf/gf/encoding/gbase64"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

//测试
func Test0(r *ghttp.Request) {
	mUser := "test"                                             //用户名
	mPwd := "test"                                              //密码
	mData := "Hello, this is a test"                            //数据明文
	mEncrypt := RSA_Encrypt([]byte(mData), "./cert/public.pem") //对数据进行公钥加密
	//向测试API提交数据
	mR, er := g.Client().Post("http://127.0.0.1:10100/api/cmd", g.Map{
		"user": mUser,
		"pwd":  mPwd,
		"data": string(gbase64.Encode(mEncrypt)), //BASE64编码字符串
	})
	//显示API返回信息
	mA := mR.ReadAllString()
	fmt.Println(er)
	r.Response.Write(mA)
}

func Test1(r *ghttp.Request) {
	mUser := "test"                               //用户名
	mPwd := "test"                                //密码
	mDataMap := g.Map{"cmd": "GetAllProcessInfo"} //数据明文
	mDataJson := gjson.New(mDataMap)
	mData := mDataJson.Export()
	mEncrypt := RSA_Encrypt([]byte(mData), "./cert/public.pem") //对数据进行公钥加密
	//向测试API提交数据
	mR, er := g.Client().Post("http://127.0.0.1:10100/api/cmd", g.Map{
		"user": mUser,
		"pwd":  mPwd,
		"data": string(gbase64.Encode(mEncrypt)), //BASE64编码字符串
	})
	//显示API返回信息
	mA := mR.ReadAllString()
	fmt.Println("连接客户端错误=", er)
	r.Response.Write(mA)
}

func Test2(r *ghttp.Request) {
	mUser := "test" //用户名
	mPwd := "test"  //密码
	mDataMap := g.Map{
		"cmd": "DisplayConf",
		"params": g.Map{
			"process": "pa_api_go",
		},
	} //数据明文
	mDataJson := gjson.New(mDataMap)
	mData := mDataJson.Export()
	mEncrypt := RSA_Encrypt([]byte(mData), "./cert/public.pem") //对数据进行公钥加密
	//向测试API提交数据
	mR, er := g.Client().Post("http://127.0.0.1:10100/api/cmd", g.Map{
		"user": mUser,
		"pwd":  mPwd,
		"data": string(gbase64.Encode(mEncrypt)), //BASE64编码字符串
	})
	//显示API返回信息
	mA := mR.ReadAllString()
	fmt.Println("连接客户端错误=", er)
	r.Response.Write(mA)
}

func Test3(r *ghttp.Request) {
	mUser := "test" //用户名
	mPwd := "test"  //密码
	mDataMap := g.Map{
		"cmd": "InsertConf",
		"params": g.Map{
			"process":     "test_process",
			"autostart":   "true",
			"autorestart": "true",
			"command":     "pwd",
			"user":        "user",
		},
	} //数据明文
	mDataJson := gjson.New(mDataMap)
	mData := mDataJson.Export()
	mEncrypt := RSA_Encrypt([]byte(mData), "./cert/public.pem") //对数据进行公钥加密
	//向测试API提交数据
	mR, er := g.Client().Post("http://127.0.0.1:10100/api/cmd", g.Map{
		"user": mUser,
		"pwd":  mPwd,
		"data": string(gbase64.Encode(mEncrypt)), //BASE64编码字符串
	})
	//显示API返回信息
	mA := mR.ReadAllString()
	fmt.Println("连接客户端错误=", er)
	r.Response.Write(mA)
}

func Test4(r *ghttp.Request) {
	mUser := "test" //用户名
	mPwd := "test"  //密码
	mDataMap := g.Map{
		"cmd": "UpdateConf",
		"params": g.Map{
			"process":     "test_process",
			"autostart":   "false",
			"autorestart": "false",
			"command":     "ls",
			"user":        "",
		},
	} //数据明文
	mDataJson := gjson.New(mDataMap)
	mData := mDataJson.Export()
	mEncrypt := RSA_Encrypt([]byte(mData), "./cert/public.pem") //对数据进行公钥加密
	//向测试API提交数据
	mR, er := g.Client().Post("http://127.0.0.1:10100/api/cmd", g.Map{
		"user": mUser,
		"pwd":  mPwd,
		"data": string(gbase64.Encode(mEncrypt)), //BASE64编码字符串
	})
	//显示API返回信息
	mA := mR.ReadAllString()
	fmt.Println("连接客户端错误=", er)
	r.Response.Write(mA)
}

func Test5(r *ghttp.Request) {
	mUser := "test" //用户名
	mPwd := "test"  //密码
	mDataMap := g.Map{
		"cmd": "RemoveConf",
		"params": g.Map{
			"process":     "test_process",
			"autostart":   "false",
			"autorestart": "false",
			"command":     "ls",
			"user":        "",
		},
	} //数据明文
	mDataJson := gjson.New(mDataMap)
	mData := mDataJson.Export()
	mEncrypt := RSA_Encrypt([]byte(mData), "./cert/public.pem") //对数据进行公钥加密
	//向测试API提交数据
	mR, er := g.Client().Post("http://127.0.0.1:10100/api/cmd", g.Map{
		"user": mUser,
		"pwd":  mPwd,
		"data": string(gbase64.Encode(mEncrypt)), //BASE64编码字符串
	})
	//显示API返回信息
	mA := mR.ReadAllString()
	fmt.Println("连接客户端错误=", er)
	r.Response.Write(mA)
}

/*
func Test6(r *ghttp.Request) {
	rpcUrl := g.Cfg().GetString("supervisor.url")
	rpc := gosupervisor.New(rpcUrl)
	a, er := GetProcessInfo(rpc, "pa_api_go")
	fmt.Println(a, er)
}
*/

func Test6(r *ghttp.Request) {
	mUser := "test" //用户名
	mPwd := "test"  //密码
	mDataMap := g.Map{
		"cmd": "GetProcessInfo",
		"params": g.Map{
			"process": "pa_api_go",
		},
	} //数据明文
	mDataJson := gjson.New(mDataMap)
	mData := mDataJson.Export()
	mEncrypt := RSA_Encrypt([]byte(mData), "./cert/public.pem") //对数据进行公钥加密
	//向测试API提交数据
	mR, er := g.Client().Post("http://127.0.0.1:10100/api/cmd", g.Map{
		"user": mUser,
		"pwd":  mPwd,
		"data": string(gbase64.Encode(mEncrypt)), //BASE64编码字符串
	})
	//显示API返回信息
	mA := mR.ReadAllString()
	fmt.Println("连接客户端错误=", er)
	r.Response.Write(mA)
}

func Test7(r *ghttp.Request) {
	mUser := "test" //用户名
	mPwd := "test"  //密码
	mDataMap := g.Map{
		"cmd": "DisplayStdoutLog",
		"params": g.Map{
			"process": "pa_api_go",
		},
	} //数据明文
	mDataJson := gjson.New(mDataMap)
	mData := mDataJson.Export()
	mEncrypt := RSA_Encrypt([]byte(mData), "./cert/public.pem") //对数据进行公钥加密
	//向测试API提交数据
	mR, er := g.Client().Post("http://127.0.0.1:10100/api/cmd", g.Map{
		"user": mUser,
		"pwd":  mPwd,
		"data": string(gbase64.Encode(mEncrypt)), //BASE64编码字符串
	})
	//显示API返回信息
	mA := mR.ReadAllString()
	fmt.Println("连接客户端错误=", er)
	r.Response.Write(mA)
}

func Test8(r *ghttp.Request) {
	mUser := "test" //用户名
	mPwd := "test"  //密码
	mDataMap := g.Map{
		"cmd":    "SysConfig",
		"params": g.Map{},
	} //数据明文
	mDataJson := gjson.New(mDataMap)
	mData := mDataJson.Export()
	mEncrypt := RSA_Encrypt([]byte(mData), "./cert/public.pem") //对数据进行公钥加密
	//向测试API提交数据
	mR, er := g.Client().Post("http://127.0.0.1:10100/api/cmd", g.Map{
		"user": mUser,
		"pwd":  mPwd,
		"data": string(gbase64.Encode(mEncrypt)), //BASE64编码字符串
	})
	//显示API返回信息
	mA := mR.ReadAllString()
	fmt.Println("连接客户端错误=", er)
	r.Response.Write(mA)
}

func Test9(r *ghttp.Request) {
	mUser := "test" //用户名
	mPwd := "test"  //密码
	mDataMap := g.Map{
		"cmd":    "SysInfo",
		"params": g.Map{},
	} //数据明文
	mDataJson := gjson.New(mDataMap)
	mData := mDataJson.Export()
	mEncrypt := RSA_Encrypt([]byte(mData), "./cert/public.pem") //对数据进行公钥加密
	//向测试API提交数据
	mR, er := g.Client().Post("http://127.0.0.1:10100/api/cmd", g.Map{
		"user": mUser,
		"pwd":  mPwd,
		"data": string(gbase64.Encode(mEncrypt)), //BASE64编码字符串
	})
	//显示API返回信息
	mA := mR.ReadAllString()
	fmt.Println("连接客户端错误=", er)
	r.Response.Write(mA)
}
