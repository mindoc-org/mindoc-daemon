# mindoc-daemon
daemon for [mindoc](https://github.com/mindoc-org/mindoc)

[简体中文](./README_zh-CN.md)

## build & run
```bash
# build
go build -o mindoc-daemon.exe main.go
# direct run mindoc
mindoc-daemon.exe
# mindoc-daemon service command
mindoc-daemon.exe -service <command>
```

## available commands for mindoc-daemon service
- `install`
- `uninstall`
- `start`
- `stop`
- `restart`

## configuration
the configuration file is `mindoc-daemon.json`, which must be in the same folder as the mindoc-daemon executable file, and the content is as follows:
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
the full configuration file example is `mindoc-daemon-full.json.example`.

explanation of main configuration items:
- `Dir` working directory of mindoc, please change it to your self.
- `Exec` the full path of the mindoc executable file, please change it to your self.
- `Stderr` the full path of error log, or standard error if empty, please change it to your self if you want use file to store log.
- `Stdout` the full path of output log, or standard output if empty, please change it to your self if you want use file to store log.

## reference
- https://stackoverflow.com/questions/43135919/how-to-run-a-shell-command-in-a-specific-folder
- https://pkg.go.dev/os/exec#Cmd.Run
- https://pkg.go.dev/os#Executable
- https://github.com/kardianos/service/blob/master/example/simple/main.go
- https://github.com/kardianos/service/issues/229