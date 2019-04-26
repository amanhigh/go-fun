import logging


class Restaurant:
    def __init__(self, env, id):
        self.id = id
        self.env = env
        self.name = "Food Point %d" % id

    def prepare_food(self, order):
        logging.info("%s: Order %d received at %s" % (self.name, order.id, self.env.now))
        yield self.env.timeout(order.dish.prep_time())

        logging.info("%s: Order %d prepared at %s" % (self.name, order.id, self.env.now))
