import logging

import simpy


class Restaurant:
    def __init__(self, env, id):
        self.id = id
        self.env = env
        self.name = "Food Point %d" % id
        self.orderStore = simpy.FilterStore(env)

    def prepare_food(self, order):
        logging.info("%s: Order: %d received at %s" % (self.name, order.id, self.env.now))
        yield self.env.timeout(order.dish.prep_time())
        yield self.orderStore.put(order)

        logging.info("%s: Order: %d prepared at %s" % (self.name, order.id, self.env.now))

    def handover_food(self, order):
        yield self.orderStore.get(lambda o: o.id == order.id)
        logging.info("%s: Order: %d handed over at %s" % (self.name, order.id, self.env.now))
