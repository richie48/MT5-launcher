# MT5 Launcher
There is presently no clean way to launch multiple metatrader5 sessions on a single task. Every session has to be run on it's own task making it difficult to run multiple accounts on the same server easily. This task aims to create multiple instances of metatrader5 launched with it's own unique session driven by config. `config-example.json` as a template of the config.json, this is to help provide a plug and use format for anyone looking to use this task for the same purpose.   

### How to setup   
Download MT5 following the guidelines for your operating system [here](https://www.metatrader5.com/en/download). If on a linux machine what this does is attempt to use wine to run mt5setup.exe. MT5 is native to windows and therefore cannot be run on a different operating system directly. Once installations are done we expect to find terminal64.exe at `/home/<your_user>/.mt5/drive_c/Program Files/MetaTrader 5`. We are now ready to run our launcher.

### How to run launcher   
Install go from the official go [website](https://go.dev/doc/install)   

This script is not packaged as it only relies on go official packages. Run with the provided code below. By default it will attempt to both create instance directories and launch scripts. To create task but not launch add `--no-launch`. It's more recommended to initialize metaTraders5 via a library officially provided, so it's best to create instances and then give metaTrader5 official packages the path to your task to run themselves.
```
BASE_DIR=/home/<your_user>/.mt5/drive_c/"Program Files"/"MetaTrader 5" go run main.go --no-launch
```   
Ensure to format if making any modification `gofmt -w main.go`   
Logging only top level issues, everything else is printed to the terminal. This is a stop on first failure task, it adviseable to cleanup folders if the script is terminated in the middle of running as the script may have made partial updates to folders. 
