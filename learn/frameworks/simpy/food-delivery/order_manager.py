import random

from models.order import Order


class OrderManager:
    def __init__(self, env, dbManager):
        self.env = env
        self.dbManager = dbManager

    def place_order(self, order):
        # Tell Restaurant to Prepare Food
        self.env.process(order.restaurant.prepare_food(order))

        # Tell Delivery Boy to Deliver Food
        self.env.process(self.dbManager.deliverOrder(order))

    def order_generator(self, interval, id):
        while True:
            yield self.env.timeout(random.randint(interval - 2, interval + 2))
            id += 1
            self.place_order(Order(id, restaurant, Dish(id)))
