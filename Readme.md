# JLink客户端

## 功能
* 转接supervisor的API，包括启动、停止、获取进程信息等
* 创建Socket Server，接收监听器发回来的消息
* 创建MQTT客户端，以便发布报警消息
* HTTP服务采用单一接口的报文通讯（详见代码）

## 提交数据格式
* 提交方法：POST
* 参数：
  * user ： access_user名，在config.toml中定义
  * pwd  : access_pwd密码，在config.toml中定义
  * data : 经过公钥加密即Base64编码的报文

## 提交报文为JSON字符串，通用格式如下
```
{
    "cmd":"命令",
    "params":{
        参数JSON记录
    }
}
```
## 返回报文为JSON字符串，通用格式如下
```
{
  "code":状态码,
  "msg":"信息文本",
  "data":{
    返回数据
  }
}
```
## 命令列表
* GetProcessInfo
  * 获取指定进程信息
  * 参数：进程名
* GetAllProcessInfo
  * 获取所有进程信息
  * 参数：无
* Restart
  * 重启supervisor服务
  * 参数：无
* StartProcess
  * 启动进程
  * 参数为：process = 进程名, wait=是否延时启动，true=是，其余为否
* StartAllProcess
  * 启动所有进程
  * 参数：无
* StopAllProcess
  * 停止所有进程
* StopProcess
  * 停止进程
* StartGroup
  * 启动组进程
* StopGroup
  * 停止组进程
* ReloadConfig
  * 重新载入设置
* DisplayConf
  * 显示指定进程的设置文件
* UpdateConf
  * 更新指定进程的设置文件
* InsertConf
  * 创建新的设置文件
* RemoveConf
  * 移除设置文件
* DisplayStdoutLog
  * 显示进程的标准输出日志内容
* DisplayStderrLog
  * 显示进程的标准错误日志内容
* SysConfig
  * 获取客户端设置信息
* SysInfo
  * 获取客户端运行环境信息
* Reset
  * 重新启动客户端
