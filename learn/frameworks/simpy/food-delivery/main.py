import logging

import simpy

from core.order_manager import OrderManager
from entities.delivery_boy import DeliveryBoy
from entities.restaurant import Restaurant
from models.order import *

logging.basicConfig(level=logging.DEBUG)

env = simpy.Environment()
boy = DeliveryBoy(env, 1)
restaurant = Restaurant(env, 1)
order = Order(1, restaurant, Dish(1))
orderManager = OrderManager(env, boy)

# Driver Interrupts last Car
orderManager.place_order(order)
env.run(until=20)
