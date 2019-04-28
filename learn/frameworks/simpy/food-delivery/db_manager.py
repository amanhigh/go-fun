import logging

import simpy

from entities.delivery_boy import DeliveryBoy


class DeliveryBoyManager:
    def __init__(self, env, config, xy_generator):
        count = config['hires']
        speed = config['speed']
        self.env = env
        self.freePool = simpy.FilterStore(env, count)
        self.orderServed = 0
        self.deliveryTimeTotal = 0
        self.xy_generator = xy_generator
        logging.critical("Hired %d Delivery boys with Speed %d" % (count, speed))
        for i in range(count):
            x, y = self.xy_generator.next()
            self.freePool.put(DeliveryBoy(self.env, i + 1, self, x, y, speed))

    def deliverOrder(self, order):
        boy = yield self.freePool.get()
        self.env.process(boy.deliver(order))

    def reportOrderServed(self, boy, deliveryTime):
        self.orderServed += 1
        self.deliveryTimeTotal += deliveryTime
        self.freePool.put(boy)

    def printSummary(self):
        logging.critical("-------- Simulation Summary ------------")
        logging.critical("Orders Served: %d" % self.orderServed)
        logging.critical("Average Delivery Time: %f" % (self.deliveryTimeTotal / self.orderServed))
