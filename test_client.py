from mt5linux_updated import MetaTrader5

import json
import os


def test_metatrader5_works_on_linux():
    mt5_client = MetaTrader5(host="0.0.0.0", port=17001)

    with open("config.json", "r") as file:
        config_data_list = json.load(file)

    BASE_DIRECTORY = os.getenv("BASE_DIR")

    for config_data in config_data_list:
        instance_executable_path = os.path.join(BASE_DIRECTORY, config_data["Path"])
        initialized = mt5_client.initialize(
            path=instance_executable_path,
            login=config_data["Login"],
            server=config_data["Server"],
            password=config_data["Password"],
        )
        if initialized:
            print("{0} got connected!".format(config_data["Name"]))
        assert mt5_client.account_info() is not None

        # prepare request
        symbol = "USDJPY"
        symbol_info = mt5_client.symbol_info(symbol)
        assert symbol_info is not None, "symbol_info({}) failed".format(symbol)
        if not symbol_info.visible:
            assert mt5_client.symbol_select(
                symbol, True
            ), "symbol_select({}) failed".format(symbol)

        lot = 1.0
        point = mt5_client.symbol_info(symbol).point
        price = mt5_client.symbol_info_tick(symbol).bid
        deviation = 20
        request = {
            "action": mt5_client.TRADE_ACTION_PENDING,
            "symbol": symbol,
            "volume": lot,
            "type": mt5_client.ORDER_TYPE_BUY_STOP,
            "price": price + 1000 * point,
            "sl": price - 1000 * point,
            "tp": price + 3000 * point,
            "deviation": deviation,
            "magic": 234000,
            "comment": "python test script",
            "type_time": mt5_client.ORDER_TIME_GTC,
            "type_filling": mt5_client.ORDER_FILLING_IOC,
        }

        # send request
        result = mt5_client.order_send(request)
        assert (
            result.retcode == mt5_client.TRADE_RETCODE_DONE
        ), "order_send failed, retcode={}".format(result.retcode)


if __name__ == "__main__":
    test_metatrader5_works_on_linux()
