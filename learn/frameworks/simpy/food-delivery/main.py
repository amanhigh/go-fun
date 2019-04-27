import logging
import random

import simpy
import yaml

from db_manager import DeliveryBoyManager
from order_manager import OrderManager
from restaurant_manager import RestaurantManager


def coordinate_generator(grid):
    while True:
        x = random.randint(0, grid['x'])
        y = random.randint(0, grid['x'])
        yield x, y


def setup(env, config):
    xy_generator = coordinate_generator(config['sim']['grid'])
    dbManager = DeliveryBoyManager(env, config['delivery'], xy_generator)
    restaurantManager = RestaurantManager(env, config['restaurant'], xy_generator)
    orderManager = OrderManager(env, dbManager, restaurantManager, xy_generator)

    # Start Order Generator
    env.process(orderManager.order_generator(interval=config['sim']['order']['generateInterval'], id=0))
    logging.info("-------- Starting Order Generation ----------")

    return dbManager


# Load Config and Setup
with open("config.yaml", 'r') as stream:
    config = yaml.load(stream)
logging.basicConfig(level=logging.INFO)

env = simpy.Environment()
dbManager = setup(env, config)

# Simulate
until = config['sim']['until']
logging.info("Running Simulation for %d Time" % until)
env.run(until=until)

dbManager.printSummary()
