import logging

import simpy
import yaml

from db_manager import DeliveryBoyManager
from order_manager import OrderManager
from restaurant_manager import RestaurantManager


def setup(env, config):
    dbManager = DeliveryBoyManager(env, config['delivery'])
    restaurantManager = RestaurantManager(env, config['restaurant'])
    orderManager = OrderManager(env, dbManager, restaurantManager)

    # Start Order Generator
    env.process(orderManager.order_generator(interval=config['sim']['order']['generateInterval'],id=0))

    return dbManager


logging.basicConfig(level=logging.INFO)

with open("config.yaml", 'r') as stream:
    config = yaml.load(stream)

env = simpy.Environment()
dbManager = setup(env, config)

# Simulate
until = config['sim']['until']
logging.info("Running Simulation for %d Time" % until)
env.run(until=until)
# dbManager.printOrdersServed()
