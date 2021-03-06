import logging

import simpy


class Restaurant:
    def __init__(self, env, id, kitchencount, x, y):
        self.id = id
        self.env = env
        self.name = "RE-%d" % id
        self.x = x
        self.y = y
        self.orderStore = simpy.FilterStore(env)
        self.kitchen = simpy.Resource(env, kitchencount)
        logging.debug("%s: Setup at X:%d,Y:%d" % (self.name, x, y))

    def prepare_food(self, order):
        logging.debug("%s (O%d): RECEIVED at %s" % (self.name, order.id, self.env.now))
        with self.kitchen.request() as req:
            yield req

            logging.info("%s (O%d): PREP_STARTED at %s" % (self.name, order.id, self.env.now))
            yield self.env.timeout(order.dish.prep_time())
            yield self.orderStore.put(order)

            logging.info("%s (O%d): #COOKED# at %s" % (self.name, order.id, self.env.now))

    def handover_food(self, order):
        yield self.orderStore.get(lambda o: o.id == order.id)
        logging.debug("%s (O%d): HANDED_OVER at %s" % (self.name, order.id, self.env.now))
