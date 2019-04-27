import logging
import random

from models.order import Order, Customer


class OrderManager:
    def __init__(self, env, dbManager, resManager, xy_generator):
        self.env = env
        self.dbManager = dbManager
        self.resManager = resManager
        self.xy_generator = xy_generator

    def place_order(self, order):
        # Tell Restaurant to Prepare Food
        self.env.process(order.restaurant.prepare_food(order))

        # Tell Delivery Boy to Deliver Food
        self.env.process(self.dbManager.deliverOrder(order))

    def order_generator(self, interval, id):
        while True:
            yield self.env.timeout(random.randint(interval - 2, interval + 2))
            id += 1
            restaurant = self.resManager.get_random_restaurant()
            dish = self.resManager.get_random_dish()
            x, y = self.xy_generator.next()
            customer = Customer(x, y)

            logging.info("NEW_ORDER (O%d): Dish %d Restaurant %d Customer (%d,%d) at %d" % (
                id, dish.id, restaurant.id, x, y, self.env.now))
            self.place_order(Order(id, restaurant, dish, customer))
