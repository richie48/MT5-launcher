# MT5 Launcher
There is presently no clean way to launch multiple metatrader5 sessions on a single task. Every session has to be run on it's own task making it difficult to run multiple accounts on the same server easily. This task aims to create multiple instances of metatrader5 launched with it's own unique session driven by config

### How to run launcher   
Iinstall go from the official go [website](https://go.dev/doc/install)
```
go run main.go
```   
Ensure to format if making any modification `gofmt -w main.go`
