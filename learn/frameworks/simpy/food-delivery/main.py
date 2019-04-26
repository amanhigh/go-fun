import logging
import random

import simpy

from core.order_manager import OrderManager
from entities.delivery_boy import DeliveryBoy
from entities.restaurant import Restaurant
from models.order import *


def order_generator(interval):
    i = 1
    while True:
        yield env.timeout(random.randint(interval - 2, interval + 2))
        i += 1
        orderManager.place_order(Order(i, restaurant, Dish(i)))


logging.basicConfig(level=logging.DEBUG)

env = simpy.Environment()
boy = DeliveryBoy(env, 1)
restaurant = Restaurant(env, 1)
orderManager = OrderManager(env, boy)

# Single Order
order = Order(1, restaurant, Dish(1))
orderManager.place_order(order)

# Order Generator
env.process(order_generator(4))

# Simulate
env.run(until=20)
