package app

import (
	"bytes"
	"fmt"
	"github.com/foolin/gomap"
	"github.com/foolin/gosupervisor"
	"github.com/gogf/gf/encoding/gini"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gfile"
	"github.com/kolo/xmlrpc"
	"os/exec"
	"strings"
)

// GetProcessInfo 获取进程信息
func GetProcessInfo(AName string) (gomap.Mapx, error) {
	rpcUrl := g.Cfg().GetString("supervisor.url")
	rpc := gosupervisor.New(rpcUrl)
	var list gomap.Mapx = gomap.NewMapx()
	client, err := xmlrpc.NewClient(rpc.Url, nil)
	if err != nil {
		return nil, err
	}
	//ret := make([]interface{}, 0)
	err = client.Call("supervisor.getProcessInfo", AName, &list)
	return list, err
}

// GetAllProcessInfo 获取所有进程信息
func GetAllProcessInfo() ([]g.Map, error) {
	rpcUrl := g.Cfg().GetString("supervisor.url")
	rpc := gosupervisor.New(rpcUrl)
	mlist, er := rpc.GetAllProcessInfo()
	if er != nil {
		return nil, er
	} else {
		var mR []g.Map
		for _, v := range mlist {
			//fmt.Println("Value=", v)
			//fmt.Printf("%v\t\t%v\t\t%v", v.Get("name"), time.Unix(v.Int64("start", 0), 0), v.Get("description"))
			mR = append(mR, v)
		}
		return mR, nil
	}
}

// StopAllProcess 停止所有进程
func StopAllProcess(AWait bool) ([]g.Map, error) {
	rpcUrl := g.Cfg().GetString("supervisor.url")
	rpc := gosupervisor.New(rpcUrl)
	mList, er := rpc.StopAllProcesses(AWait)
	var mR []g.Map
	for _, v := range mList {
		mR = append(mR, v)
	}
	return mR, er

}

// StartAllProcess 启动所有进程
func StartAllProcess(AWait bool) ([]g.Map, error) {
	rpcUrl := g.Cfg().GetString("supervisor.url")
	rpc := gosupervisor.New(rpcUrl)
	mList, er := rpc.StartAllProcess(AWait)
	var mR []g.Map
	for _, v := range mList {
		mR = append(mR, v)
	}
	return mR, er
}

// RestartSup 重启supvivior服务
func RestartSup() (bool, error) {
	rpcUrl := g.Cfg().GetString("supervisor.url")
	rpc := gosupervisor.New(rpcUrl)
	mR, er := rpc.Restart()
	return mR, er
}

// StartProcess 启动进程
func StartProcess(AName string, AWait bool) (bool, error) {
	rpcUrl := g.Cfg().GetString("supervisor.url")
	rpc := gosupervisor.New(rpcUrl)
	mR, er := rpc.StartProcess(AName, AWait)
	return mR, er
}

// StopProcess 停止进程
func StopProcess(AName string, AWait bool) (bool, error) {
	rpcUrl := g.Cfg().GetString("supervisor.url")
	rpc := gosupervisor.New(rpcUrl)
	mR, er := rpc.StopProcess(AName, AWait)
	return mR, er
}

// StartGroup 启动组进程
func StartGroup(AName string, AWait bool) ([]g.Map, error) {
	rpcUrl := g.Cfg().GetString("supervisor.url")
	rpc := gosupervisor.New(rpcUrl)
	mList, er := rpc.StartProcessGroup(AName, AWait)
	var mR []g.Map
	for _, v := range mList {
		mR = append(mR, v)
	}
	return mR, er
}

// StopGroup 启动组进程
func StopGroup(AName string, AWait bool) ([]g.Map, error) {
	rpcUrl := g.Cfg().GetString("supervisor.url")
	rpc := gosupervisor.New(rpcUrl)
	mList, er := rpc.StopProcessGroup(AName, AWait)
	var mR []g.Map
	for _, v := range mList {
		mR = append(mR, v)
	}
	return mR, er
}

// ReloadConfig 重新载入设置
func ReloadConfig() (bool, error) {
	rpcUrl := g.Cfg().GetString("supervisor.url")
	rpc := gosupervisor.New(rpcUrl)
	mR, er := rpc.ReloadConfig()
	return mR, er
}

// DisplayConf 显示指定进程的设置文件
func DisplayConf(AProcess string) (g.Map, error) {
	//命令示例： sudo grep -r ":pa_api_go" /etc/
	//返回：pingan_api_go.conf:[program:pa_api_go]，从而确定文件
	mS, err := DoGrep(AProcess, g.Cfg().GetString("supervisor.dir"))

	mTmp := strings.Split(mS, ":")
	var mCode int
	var mErr string
	var mFile string
	var mContent string
	if len(mTmp) > 0 {
		mFile = mTmp[0]
		if gfile.Exists(mFile) == true {
			mContent = gfile.GetContents(mFile)
			mCode = 200
		} else {
			mCode = 404
			mContent = ""
		}
	} else {
		mFile = ""
		mCode = 404
	}
	if err == nil {
		mErr = ""
	} else {
		mErr = err.Error()
	}
	return g.Map{
		"code":    mCode,
		"process": AProcess,
		"file":    mFile,
		"error":   mErr,
		"content": mContent,
	}, err
}

// DoGrep 查找包含字符串的文件，返回的是文件名：关键字样式的字符串
func DoGrep(AKey string, ADir string) (string, error) {
	mR, er := DoShell(fmt.Sprintf("grep -r \"%s\" %s", AKey, ADir))
	return mR, er
}

// Go语言调用命令行程序的代码
func DoShell(Avalue string) (string, error) {

	cmd := exec.Command("/bin/sh", "-c", Avalue) //linux
	//cmd := exec.Command("cmd", Avalue)
	out, err := cmd.Output()
	if err == nil {
		return string(out), nil
	} else {
		return "", err
	}
}

// DoGenerateConf 生成设置文件内容
func DoGenerateConf(AParams g.Map) string {
	var mTmp bytes.Buffer
	//获取参数变量
	if _, ok := AParams["process"]; ok {
		mProcess := fmt.Sprint(AParams["process"])
		mTmp.WriteString(fmt.Sprintf("[program:%s]\n", mProcess))
	}
	if _, ok := AParams["command"]; ok {
		mCommand := fmt.Sprint(AParams["command"])
		mTmp.WriteString(fmt.Sprintf("command=%s\n", mCommand))
	}
	if _, ok := AParams["numprocs"]; ok {
		mNumProcs := fmt.Sprint(AParams["numprocs"])
		if strings.TrimSpace(mNumProcs) != "" {
			mTmp.WriteString(fmt.Sprintf("numprocs=%s\n", mNumProcs))
		}
	}
	if _, ok := AParams["numprocs_start"]; ok {
		mNumprocs_start := fmt.Sprint(AParams["numprocs_start"])
		if strings.TrimSpace(mNumprocs_start) != "" {
			mTmp.WriteString(fmt.Sprintf("numprocs_start=%s\n", mNumprocs_start))
		}
	}
	if _, ok := AParams["priority"]; ok {
		mPriority := fmt.Sprint(AParams["priority"])
		if strings.TrimSpace(mPriority) != "" {
			mTmp.WriteString(fmt.Sprintf("priority=%s\n", mPriority))
		}
	}
	if _, ok := AParams["autostart"]; ok {
		mAutoStart := fmt.Sprint(AParams["autostart"])
		if strings.TrimSpace(mAutoStart) != "" {
			mTmp.WriteString(fmt.Sprintf("autostart=%s\n", mAutoStart))
		}
	}
	if _, ok := AParams["autorestart"]; ok {
		mAutoRestart := fmt.Sprint(AParams["autorestart"])
		if strings.TrimSpace(mAutoRestart) != "" {
			mTmp.WriteString(fmt.Sprintf("autorestart=%s\n", mAutoRestart))
		}
	}
	if _, ok := AParams["user"]; ok {
		mUser := fmt.Sprint(AParams["user"])
		if strings.TrimSpace(mUser) != "" {
			mTmp.WriteString(fmt.Sprintf("user=%s\n", mUser))
		}
	}
	return mTmp.String()
}

// DoGenerateConfFilename 生成设置文件名
func DoGenerateConfFilename(AParams g.Map) string {
	//获取参数变量
	mProcess := fmt.Sprint(AParams["process"])
	return strings.TrimSpace(mProcess) + ".conf"
}

//检查supervisor是否安装，只需要检查设置文件中的supervisor.conf是否存在即可
func isSpInstalled() bool {
	return gfile.Exists(g.Cfg().GetString("supervisor.conf"))
}

//检查supervisor设置中是否开启RPC服务
func SpInetOpened() bool {
	mFile := gfile.GetContents(g.Cfg().GetString("supervisor.conf"))
	mTom, _ := gini.Decode([]byte(mFile))
	if _, ok := mTom["inet_http_server"]; ok {
		return true
	} else {
		return false
	}
}

//在supervisor.conf文件中增加RPC设置
func SpOpenRPC() error {
	mFile := gfile.GetContents(g.Cfg().GetString("supervisor.conf"))
	mFile = mFile + "\n[inet_http_server]\nport=127.0.0.1:9001\n"
	return gfile.PutContents(g.Cfg().GetString("supervisor.conf"), mFile)
}
