# mindoc-daemon
[mindoc](https://github.com/mindoc-org/mindoc) 守护进程

[English](./README.md)

## 编译 & 运行
```bash
# 编译
go build -o mindoc-daemon.exe main.go
# 直接运行 mindoc
mindoc-daemon.exe
# mindoc-daemon service 命令
mindoc-daemon.exe -service <命令>
```

## mindoc-daemon service 可用的命令
- `install`
- `uninstall`
- `start`
- `stop`
- `restart`

## 配置
配置文件是 `mindoc-daemon.json`, 必须与 mindoc-daemon 可执行文件同一目录, 配置文件内容如下:
```json
{
    "Name": "MinDoc",
    "DisplayName": "MinDoc",
    "Description": "A document online management program.",
    
    "Dir": "E:\\own\\go\\mindoc",
    "Exec": "E:\\own\\go\\mindoc\\mindoc_windows_amd64.exe",
    "Args": [],
    "Env": [
        "MINDOC_RUN_MODE=dev",
        "MINDOC_HIGHLIGHT_STYLE=github"
    ],
    
    "Stderr": "",
    "Stdout": ""
}
```
主要配置项解释如下:
- `Dir` mindoc 的工作目录
- `Exec` mindoc 可执行文件的完整路径
- `Stderr` 错误日志完整文件路径, 留空则使用标准错误
- `Stdout` 输出日志完整文件路径, 留空则使用标准输出

## 参考
- https://stackoverflow.com/questions/43135919/how-to-run-a-shell-command-in-a-specific-folder
- https://pkg.go.dev/os/exec#Cmd.Run
- https://pkg.go.dev/os#Executable
- https://github.com/kardianos/service/blob/master/example/simple/main.go
- https://github.com/kardianos/service/issues/229