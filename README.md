# MT5 Launcher
There is presently no clean way to launch multiple metatrader5 sessions on a single task. Every session has to be run on it's own task making it difficult to run multiple accounts on the same task simultaneously. This task aims to create multiple instances of metatrader5 launched with it's own unique session driven by config. `config-example.json` as a template of the config.json, this is to help provide a plug and use format for anyone looking to use this task for the same purpose.   

### How to setup   
Download MT5 following the guidelines for your operating system [here](https://www.metatrader5.com/en/download). If on a linux machine what this does is attempt to use wine to run mt5setup.exe. MT5 is native to windows and therefore cannot be run on a different operating system directly. Once installations are done we expect to find terminal64.exe at `/home/<your_user>/.mt5/drive_c/Program Files/MetaTrader 5`. We are now ready to run our launcher.

### How to run launcher   
Install go from the official go [website](https://go.dev/doc/install)   

This script is not packaged as it only relies on go official packages. Run with the provided code below. By default it will attempt to both create instance directories and run the launch scripts to start all the instances on metaTrader5. To create task but not launch add `--no-launch`. It's recommended to initialize metaTrader5 via a library officially provided as there are reportedly issues with logging in at startup using the config.ini. It's best to create instances and then give a metaTrader5 official library package the path to your task to run themselves ([`For examples`](https://www.mql5.com/en/docs/python_metatrader5)).

The BASE_DIR is on the path to the original program, the full path to the original program is gotten by combining with SRC_DIR. The BASE_DIR is where we want to put our new instance also.    

```
BASE_DIR=/home/<your_user>/.mt5/drive_c/"Program Files" SRC_DIR="MetaTrader 5" go run main.go --no-launch
```   

### How to run test   
To test we would need to install [pymt5linux](https://pypi.org/project/pymt5linux/). Make sure to follow the steps outlined to get it running. The goal of this package to provide a way to run metaTrader5 on linux. It does it by forwarding request to a different port where we have our windows emulator running metaTrader5. That way metaTrader5 is able to carry out our request as if we were on a windows machine. Install pymt5linux on windows emulator but use [mt5linux_updated](https://pypi.org/project/mt5linux-updated/) on linux. pymt5linux needs to run on both linux and windows with python3.13 but version 3.13 is still an experimental version, therefore it's better to run on python3.12 where i can as this is a stable version.  Now that we have everything setup we can run our test!
```
BASE_DIR=/home/<your_user>/.mt5/drive_c/"Program Files" python3.12 test_client.py
```   
   
Ensure to format if making any modification `gofmt -w main.go`   
Logging only top level issues, everything else is printed to the terminal. This is a stop on first failure task, it adviseable to cleanup folders if the script is terminated in the middle of running as the script may have made partial updates to folders. 
