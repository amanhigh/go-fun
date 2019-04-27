import logging
import random

import simpy

from db_manager import DeliveryBoyManager
from entities.restaurant import Restaurant
from models.order import Order, Dish
from order_manager import OrderManager


def order_generator(interval, id):
    while True:
        yield env.timeout(random.randint(interval - 2, interval + 2))
        id += 1
        orderManager.place_order(Order(id, restaurant, Dish(id)))


logging.basicConfig(level=logging.DEBUG)

env = simpy.Environment()
restaurant = Restaurant(env, id=1, kitchencount=2)
dbManager = DeliveryBoyManager(env, count=2)
orderManager = OrderManager(env, dbManager)

# Single Order
id = 1
order = Order(id, restaurant, Dish(id))
orderManager.place_order(order)

# Order Generator
env.process(order_generator(interval=4, id=id))

# Simulate
env.run(until=20)
