package app

import (
	"fmt"
	"github.com/gogf/gf/encoding/gbase64"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gfile"
	"github.com/shirou/gopsutil/cpu"
	_ "github.com/shirou/gopsutil/cpu"
	_ "github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	_ "github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	_ "github.com/shirou/gopsutil/mem"
	_ "github.com/shirou/gopsutil/net"
	_ "github.com/shirou/gopsutil/process"
	"os"
	"strings"
	"time"
)

/*
	单一接口报文处理函数
	报文为JSON字符串，采用RAS加密，收到后解密并分类处理
	报文结构示例：{
		user://简单的用户名
		pwd://简单的密码
		data://RSA加密的报文数据
	}
*/

func SvdProce(r *ghttp.Request) {
	//输出日志
	fmt.Printf("Remote Call from %s at %s\n", r.RemoteAddr, fmt.Sprint(time.Now()))
	//解析提交的报文，报文统一POST到Data字段
	mUser := r.GetString("user")
	mPwd := r.GetString("pwd")
	//校验用户名和密码，用户名和密码写在config.toml文件中
	if CheckAccessUserPwd(mUser, mPwd) == false {
		//校验用户名和密码错误，返回403状态
		r.Response.WriteStatusExit(403)
	} else {
		//校验正确，开始获取data报文并解密
		mData := r.GetString("data")
		//fmt.Println("Data=", mData)
		mDesData, _ := gbase64.Decode([]byte(mData))
		//fmt.Println("unbase64", mDesData, er)
		plainText := RSA_Decrypt(mDesData, g.Cfg().GetString("rsa.private_pem"))
		mJsonStr := string(plainText)  //转换为字符串
		mJson := gjson.New(mJsonStr)   //转换为JSON
		mCmd := mJson.GetString("cmd") //获取命令字
		//mParams := mJson.GetMap("params") //获取参数Map
		switch mCmd {
		case "GetProcessInfo":
			//获取进程信息，参数是进程名
			mR := doGetProcessInfo(mJson)
			_ = r.Response.WriteJson(mR)
		case "GetAllProcessInfo":
			//获取所有进程信息,不需要参数
			mR := doGetAllProcessInfo()
			_ = r.Response.WriteJson(mR)
		case "Restart":
			// 重启supvivior服务
			mR := doRestartSup()
			_ = r.Response.WriteJson(mR)
		case "StartProcess":
			//启动进程，获取参数中的process名
			mR := doStartProcess(mJson)
			_ = r.Response.WriteJson(mR)
		case "StartAllProcess":
			//启动所有进程
			mR := doStartAllProcess(mJson)
			_ = r.Response.WriteJson(mR)
		case "StopProcess":
			//停止进程
			mR := doStopProcess(mJson)
			_ = r.Response.WriteJson(mR)
		case "StopAllProcess":
			//停止所有进程
			mR := doStopAllProcess(mJson)
			_ = r.Response.WriteJson(mR)
		case "StartGroup":
			//启动组进程
			mR := doStartGroup(mJson)
			_ = r.Response.WriteJson(mR)
		case "StopGroup":
			//停止组进程
			mR := doStopGroup(mJson)
			_ = r.Response.WriteJson(mR)
		case "ReloadConfig":
			//重新载入设置
			mR := doReloadConfig(mJson)
			_ = r.Response.WriteJson(mR)
		case "DisplayConf":
			//显示指定进程的设置文件
			mR := doDisplayConf(mJson)
			_ = r.Response.WriteJson(mR)
		case "InsertConf":
			//创建新的设置文件
			mR := doInsertConf(mJson)
			_ = r.Response.WriteJson(mR)
		case "UpdateConf":
			//更新指定进程的设置文件
			mR := doUpdateConf(mJson)
			_ = r.Response.WriteJson(mR)
		case "RemoveConf":
			//移除设置文件
			mR := doRemoveConf(mJson)
			_ = r.Response.WriteJson(mR)
		case "DisplayStdoutLog":
			//显示进程的标准输出日志内容
			mR := doDisplayStdoutLog(mJson)
			_ = r.Response.WriteJson(mR)
		case "DisplayStderrLog":
			//显示进程的标准错误日志内容
			mR := doDisplayStderrLog(mJson)
			_ = r.Response.WriteJson(mR)
		case "SysConfig":
			//获取客户端设置信息
			mR := doDisplaySysConfig()
			_ = r.Response.WriteJson(mR)
		case "SysInfo":
			//获取客户端运行环境信息
			mR := doDisplaySysInfo()
			_ = r.Response.WriteJson(mR)
		case "Reset":
			//重新启动本服务器
			mR := doReset()
			_ = r.Response.WriteJson(mR)
		}
		//将返回的信息加密并Base64后返回
		//r.Response.Write(plainText)
	}
}

//校验用户名和密码
func CheckAccessUserPwd(AUser string, APwd string) bool {
	mUser := g.Cfg().GetString("rsa.access_user")
	mPwd := g.Cfg().GetString("rsa.access_pwd")
	if mUser == AUser && mPwd == APwd {
		return true
	} else {
		return false
	}
}

//停止所有进程
func doStopAllProcess(AJson *gjson.Json) g.Map {
	mParams := AJson.GetMap("params")
	mWait := fmt.Sprint(mParams["wait"]) == "true"
	mR, er := StopAllProcess(mWait)
	var mEr string
	var mCode int
	if er == nil {
		mEr = ""
		mCode = 200
	} else {
		mEr = er.Error()
		mCode = 501
	}
	return g.Map{
		"code":  mCode,
		"error": mEr,
		"msg":   "StopAllProcess",
		"data":  mR,
	}
}

//停止组进程
func doStopGroup(AJson *gjson.Json) g.Map {
	mParams := AJson.GetMap("params")
	mWait := fmt.Sprint(mParams["wait"]) == "true"
	mGroup := fmt.Sprint(mParams["group"])
	mR, er := StopGroup(mGroup, mWait)
	var mEr string
	var mCode int
	if er == nil {
		mEr = ""
		mCode = 200
	} else {
		mEr = er.Error()
		mCode = 501
	}
	return g.Map{
		"code":  mCode,
		"error": mEr,
		"msg":   "StopGroup",
		"data":  mR,
	}
}

//启动组进程
func doStartGroup(AJson *gjson.Json) g.Map {
	mParams := AJson.GetMap("params")
	mWait := fmt.Sprint(mParams["wait"]) == "true"
	mGroup := fmt.Sprint(mParams["group"])
	mR, er := StartGroup(mGroup, mWait)
	var mEr string
	var mCode int
	if er == nil {
		mEr = ""
		mCode = 200
	} else {
		mEr = er.Error()
		mCode = 501
	}
	return g.Map{
		"code":  mCode,
		"error": mEr,
		"msg":   "StartGroup",
		"data":  mR,
	}
}

//重新载入设置
func doReloadConfig(AJson *gjson.Json) g.Map {
	mR, er := ReloadConfig()
	var mEr string
	var mCode int
	if er == nil {
		mEr = ""
		mCode = 200
	} else {
		mEr = er.Error()
		mCode = 501
	}
	return g.Map{
		"code":  mCode,
		"msg":   "ReloadConfig",
		"data":  mR,
		"error": mEr,
	}
}

//启动所有进程
func doStartAllProcess(AJson *gjson.Json) g.Map {
	mParams := AJson.GetMap("params")
	mWait := fmt.Sprint(mParams["wait"]) == "true"
	mR, er := StartAllProcess(mWait)
	var mEr string
	var mCode int
	if er == nil {
		mEr = ""
		mCode = 200
	} else {
		mEr = er.Error()
		mCode = 501
	}
	return g.Map{
		"code":  mCode,
		"error": mEr,
		"msg":   "StartAllProcess",
		"data":  mR,
	}
}

//停止进程，获取参数中的process名
func doStopProcess(AJson *gjson.Json) g.Map {
	mParams := AJson.GetMap("params")
	mProcessName := fmt.Sprint(mParams["process"])
	mWait := fmt.Sprint(mParams["wait"]) == "true"
	var mEr string
	var mCode int
	mR, er := StopProcess(mProcessName, mWait)
	if er == nil {
		mEr = ""
		mCode = 200
	} else {
		mEr = er.Error()
		mCode = 501
	}
	return g.Map{
		"code":  mCode,
		"error": mEr,
		"msg":   "StopProcess",
		"data":  mR,
	}
}

//创建新的设置文件
func doInsertConf(AJson *gjson.Json) g.Map {
	mParams := AJson.GetMap("params")
	mContent := DoGenerateConf(mParams)          //生成设置文件内容
	mFilename := DoGenerateConfFilename(mParams) //生成设置文件名
	mPath := gfile.Join(g.Cfg().GetString("supervisor.dir"), mFilename)
	//写入文件
	err := gfile.PutContents(mPath, mContent)
	var mCode int
	var mErr string
	if err == nil {
		mCode = 200
		mErr = ""
	} else {
		mCode = 403
		mErr = err.Error()
	}
	return g.Map{
		"code":  mCode,
		"error": mErr,
		"msg":   "InsertConf",
		"data":  mPath,
	}
}

//移除指定的进程设置文件
func doRemoveConf(AJson *gjson.Json) g.Map {
	//获取进程名参数
	mParams := AJson.GetMap("params")
	mProcess := fmt.Sprint(mParams["process"])
	//查找设置文件名
	mS, _ := DoGrep(mProcess, g.Cfg().GetString("supervisor.dir"))
	mTmp := strings.Split(mS, ":")
	var mConfFile string
	var mCode int
	var mErr string
	var mPath string
	if len(mTmp) > 0 {
		mConfFile = mTmp[0]
		if gfile.Exists(mConfFile) == true {
			//重新生成设置文件内容
			mFilename := DoGenerateConfFilename(mParams) //生成设置文件名
			mPath = gfile.Join(g.Cfg().GetString("supervisor.dir"), mFilename)
			//移除文件
			err := gfile.Remove(mPath)
			if err != nil {
				mCode = 501
				mErr = err.Error()
			} else {
				mCode = 200
				mErr = ""
			}
		} else {
			mConfFile = ""
			mCode = 404
			mErr = "Conf file missing"
		}
	} else {
		mConfFile = ""
		mCode = 404
		mErr = "Conf file missing"
	}
	return g.Map{
		"code":  mCode,
		"error": mErr,
		"msg":   "RemoveConf",
		"data":  mPath,
	}
}

//更新指定进程的设置文件
func doUpdateConf(AJson *gjson.Json) g.Map {
	//获取进程名参数
	mParams := AJson.GetMap("params")
	mProcess := fmt.Sprint(mParams["process"])
	//查找设置文件名
	mS, _ := DoGrep(mProcess, g.Cfg().GetString("supervisor.dir"))
	mTmp := strings.Split(mS, ":")
	var mConfFile string
	var mCode int
	var mErr string
	var mPath string
	if len(mTmp) > 0 {
		mConfFile = mTmp[0]
		if gfile.Exists(mConfFile) == true {
			//重新生成设置文件内容
			mContent := DoGenerateConf(mParams)          //生成设置文件内容
			mFilename := DoGenerateConfFilename(mParams) //生成设置文件名
			mPath = gfile.Join(g.Cfg().GetString("supervisor.dir"), mFilename)
			//写入文件
			err := gfile.PutContents(mPath, mContent)
			if err != nil {
				mCode = 501
				mErr = err.Error()
			} else {
				mCode = 200
				mErr = ""
			}
		} else {
			mConfFile = ""
			mCode = 404
			mErr = "Conf file missing"
		}
	} else {
		mConfFile = ""
		mCode = 404
		mErr = "Conf file missing"
	}
	return g.Map{
		"code":  mCode,
		"error": mErr,
		"msg":   "UpdateConf",
		"data":  mPath,
	}
}

//显示标准输出日志
func doDisplayStderrLog(AJson *gjson.Json) g.Map {
	mParams := AJson.GetMap("params")
	mProcessName := fmt.Sprint(mParams["process"])
	//获取进程信息
	mInfo, er := GetProcessInfo(mProcessName)
	var mEr string
	var mCode int
	var mContent string
	if _, ok := mInfo["stdout_logfile"]; ok {
		mLog := fmt.Sprint(mInfo["stderr_logfile"])
		if gfile.Exists(mLog) {
			mContent = gfile.GetContents(mLog)
			mCode = 200
			mEr = ""
		}
	} else {
		mCode = 501
		mEr = er.Error() + " & PARAMETER MISSING"
	}
	return g.Map{
		"code":  mCode,
		"error": mEr,
		"msg":   "DisplayStdoutLog",
		"data":  mContent,
	}
}

//显示标准输出日志
func doDisplayStdoutLog(AJson *gjson.Json) g.Map {
	mParams := AJson.GetMap("params")
	mProcessName := fmt.Sprint(mParams["process"])
	//获取进程信息
	mInfo, er := GetProcessInfo(mProcessName)
	var mEr string
	var mCode int
	var mContent string
	if _, ok := mInfo["stdout_logfile"]; ok {
		mLog := fmt.Sprint(mInfo["stdout_logfile"])
		if gfile.Exists(mLog) {
			mContent = gfile.GetContents(mLog)
			mCode = 200
			mEr = ""
		}
	} else {
		mCode = 501
		mEr = er.Error() + " & PARAMETER MISSING"
	}
	return g.Map{
		"code":  mCode,
		"error": mEr,
		"msg":   "DisplayStdoutLog",
		"data":  mContent,
	}

}

//显示指定进程的Conf文件信息
func doDisplayConf(AJson *gjson.Json) g.Map {
	mParams := AJson.GetMap("params")
	mProcessName := fmt.Sprint(mParams["process"])
	var mEr string
	var mCode int
	mR, er := DisplayConf(mProcessName)
	if er == nil {
		mEr = ""
		mCode = 200
	} else {
		mEr = er.Error()
		mCode = 501
	}
	return g.Map{
		"code":  mCode,
		"error": mEr,
		"msg":   "DisplayConf",
		"data":  mR,
	}
}

//启动进程，获取参数中的process名
func doStartProcess(AJson *gjson.Json) g.Map {
	mParams := AJson.GetMap("params")
	mProcessName := fmt.Sprint(mParams["process"])
	mWait := fmt.Sprint(mParams["wait"]) == "true"
	var mEr string
	var mCode int
	mR, er := StartProcess(mProcessName, mWait)
	if er == nil {
		mEr = ""
		mCode = 200
	} else {
		mEr = er.Error()
		mCode = 501
	}
	return g.Map{
		"code":  mCode,
		"error": mEr,
		"msg":   "StartProcess",
		"data":  mR,
	}
}

// 重启supvivior服务
func doRestartSup() g.Map {
	mR, er := RestartSup()
	var mEr string
	var mCode int
	if er == nil {
		mEr = ""
		mCode = 200
	} else {
		mEr = er.Error()
		mCode = 501
	}
	return g.Map{
		"code":  mCode,
		"msg":   "RestartSup",
		"data":  mR,
		"error": mEr,
	}
}

//执行GetAllProcessInfo，返回值是报文Map
func doGetAllProcessInfo() g.Map {
	mR, er := GetAllProcessInfo()
	var mEr string
	var mCode int
	if er == nil {
		mEr = ""
		mCode = 200
	} else {
		mEr = er.Error()
		mCode = 501
	}
	return g.Map{
		"code":  mCode,
		"msg":   "GetAllProcessInfo",
		"data":  mR,
		"error": mEr,
	}
}

//获取进程信息，doGetProcessInfo
func doGetProcessInfo(AJson *gjson.Json) g.Map {
	mParams := AJson.GetMap("params")
	mProcessName := fmt.Sprint(mParams["process"])
	var mEr string
	var mCode int
	mR, er := GetProcessInfo(mProcessName)
	if er == nil {
		mEr = ""
		mCode = 200
	} else {
		mEr = er.Error()
		mCode = 501
	}
	return g.Map{
		"code":  mCode,
		"error": mEr,
		"msg":   "GetProcessInfo",
		"data":  mR,
	}
}

//显示客户端设置信息
func doDisplaySysConfig() g.Map {
	return g.Map{
		"code":  200,
		"error": "",
		"msg":   "SysConfig",
		"data":  g.Cfg().Map(),
	}
}

//重新启动本服务器
func doReset() g.Map {
	var mCode int
	var mErr string
	err := ghttp.RestartAllServer()
	if err != nil {
		mCode = 403
		mErr = err.Error()
	} else {
		mCode = 200
		mErr = ""
	}
	return g.Map{
		"code":  mCode,
		"error": mErr,
		"msg":   "Reset",
		"data":  "",
	}
}

//显示系统信息
func doDisplaySysInfo() g.Map {
	return g.Map{
		"code":  200,
		"error": "",
		"msg":   "SysInfo",
		"data":  CollectOSInfo(),
	}
}

//搜集系统信息
func CollectOSInfo() g.Map {
	mCurrentDir, _ := os.Getwd()
	mHostInfo, _ := host.Info()
	mCpuInfo, _ := cpu.Info()
	mMemInfo, _ := mem.VirtualMemory()
	return g.Map{
		"CurrentDir": mCurrentDir, //程序所在目录
		"host":       mHostInfo,   //主机信息
		"cpu":        mCpuInfo,    //Cpu信息
		"mem":        mMemInfo,    //内存信息
	}
}
