# TODO: Clean up this section adn add details to readme
# install python3.13 in wine, to be able to install succesfully follow similar steps as to how mt5setup was installed with wine

# all we are doing right now is forwarding request via the package https://pypi.org/project/pymt5linux/ to mt5 running on wine
# should be enough to get mt5 task and python api running on linux

# https://www.mql5.com/en/docs/python_metatrader5

# TODO: Support multiple instances, relocate logs to wine directory

from pymt5linux import MetaTrader5

import json

mt5_client = MetaTrader5()

with open('config.json', 'r') as file:
    config_data = json.load(file)

authorized = mt5_client.initialize(
    path="/home/richard/.mt5/drive_c/Program Files/MetaTrader 5/terminal64.exe",
    login=config_data[0]["Login"],
    server=config_data[0]["Server"],
    password=config_data[0]["Password"],
)

if authorized:
    print("connection was authorized")
    print(mt5_client.account_info())
else:
    print("not authorized")
