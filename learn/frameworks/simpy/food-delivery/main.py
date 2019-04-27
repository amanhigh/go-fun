import logging

import simpy
import yaml

from db_manager import DeliveryBoyManager
from restaurant_manager import RestaurantManager


def setup(env, config):
    dbManager = DeliveryBoyManager(env, config['delivery'])
    restaurantManager = RestaurantManager(env, config['restaurant'])
    # orderManager = OrderManager(env, dbManager)

    # Single Order
    # id = 1
    # order = Order(id, restaurant, Dish(id))
    # orderManager.place_order(order)

    # Order Generator
    # env.process(orderManager.order_generator(interval=4, id=id))

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
