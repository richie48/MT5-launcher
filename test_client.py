from pymt5linux import MetaTrader5

import json
import os


mt5_client = MetaTrader5()

with open("config.json", "r") as file:
    config_data_list = json.load(file)

BASE_DIRECTORY = os.getenv("BASE_DIR")

for config_data in config_data_list:
    instance_executable_path = os.path.join(BASE_DIRECTORY, config_data["Path"])
    authorized = mt5_client.initialize(
        path=instance_executable_path,
        login=config_data["Login"],
        server=config_data["Server"],
        password=config_data["Password"],
    )

    print("{0} got connected!".format(config_data["Name"]))
    assert mt5_client.account_info() is not None
