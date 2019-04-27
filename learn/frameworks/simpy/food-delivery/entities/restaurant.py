import logging

import simpy


class Restaurant:
    def __init__(self, env, id,kitchencount):
        self.id = id
        self.env = env
        self.name = "RE-%d" % id
        self.orderStore = simpy.FilterStore(env)
        self.kitchen = simpy.Resource(env, kitchencount)

    def prepare_food(self, order):
        logging.info("%s (O%d): received at %s" % (self.name, order.id, self.env.now))
        with self.kitchen.request() as req:
            yield req

            logging.info("%s (O%d): prep-started at %s" % (self.name, order.id, self.env.now))
            yield self.env.timeout(order.dish.prep_time())
            yield self.orderStore.put(order)

            logging.info("%s (O%d): cooked at %s" % (self.name, order.id, self.env.now))

    def handover_food(self, order):
        yield self.orderStore.get(lambda o: o.id == order.id)
        logging.info("%s (O%d): handed over at %s" % (self.name, order.id, self.env.now))
