# MT5 Launcher
There is presently no clean way to launch multiple metatrader5 sessions on a single task. Every session has to be run on it's own task making it difficult to run multiple accounts on the same server easily. This task aims to create multiple instances of metatrader5 launched with it's own unique session driven by config. `config-example.json` as a template of the config.json, this is to help provide a plug and use format for anyone looking to use this task for the same purpose.

### How to run launcher   
Install go from the official go [website](https://go.dev/doc/install)   

This script is not packages as it only relies on go official packages
```
go run main.go
```   
Ensure to format if making any modification `gofmt -w main.go`
