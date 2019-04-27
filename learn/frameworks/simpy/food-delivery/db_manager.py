import logging

import simpy

from entities.delivery_boy import DeliveryBoy


class DeliveryBoyManager:
    def __init__(self, env, config):
        count = config['hires']
        self.env = env
        self.freePool = simpy.Store(env, count)
        self.orderServed = 0
        logging.info("Hired %d Delivery boys" % count)
        for i in range(count):
            self.freePool.put(DeliveryBoy(self.env, i, self.freePool))

    def deliverOrder(self, order):
        boy = yield self.freePool.get()
        self.env.process(boy.deliver(order))
        self.orderServed += 1

    def printOrdersServed(self):
        logging.info("Orders Served: %d" % self.orderServed)
