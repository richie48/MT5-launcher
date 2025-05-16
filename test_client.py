from pymt5linux import MetaTrader5

import json

mt5_client = MetaTrader5()

with open('config.json', 'r') as file:
    config_data = json.load(file)

# TODO: Support multiple instances, relocate logs to wine directory
authorized = mt5_client.initialize(
    path="/home/richard/.mt5/drive_c/Program Files/MetaTrader 5/terminal64.exe",
    login=config_data[0]["Login"],
    server=config_data[0]["Server"],
    password=config_data[0]["Password"],
)

assert mt5_client.account_info() is not None
