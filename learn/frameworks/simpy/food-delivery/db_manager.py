import logging

import simpy

from entities.delivery_boy import DeliveryBoy


class DeliveryBoyManager:
    def __init__(self, env, config, xy_generator):
        count = config['hires']
        speed = config['speed']
        self.env = env
        self.freePool = simpy.Store(env, count)
        self.orderServed = 0
        self.xy_generator = xy_generator
        logging.info("Hired %d Delivery boys with Speed %d" % (count, speed))
        for i in range(count):
            x, y = self.xy_generator.next()
            self.freePool.put(DeliveryBoy(self.env, i + 1, self.freePool, x, y, speed))

    def deliverOrder(self, order):
        boy = yield self.freePool.get()
        self.env.process(boy.deliver(order))
        self.orderServed += 1

    def printOrdersServed(self):
        logging.info("Orders Served: %d" % self.orderServed)
