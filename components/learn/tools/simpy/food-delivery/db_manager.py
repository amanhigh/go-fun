import logging
import sys

import simpy

from entities.delivery_boy import DeliveryBoy


class DeliveryBoyManager:
    def __init__(self, env, config, xy_generator):
        speed = config['speed']
        self.hireCount = config['hires']
        self.algoConfig = config['algo']
        self.env = env
        self.freePool = simpy.FilterStore(env, self.hireCount)
        self.orderServed = 0
        self.deliveryTimeTotal = 0
        self.idleTimeTotal = 0
        self.xy_generator = xy_generator
        logging.critical(
            "Delivery Manager: Hired %d boys with Speed %d, Algo: %s" % (self.hireCount, speed, self.algoConfig["type"]))
        for i in range(self.hireCount):
            x, y = self.xy_generator.next()
            self.freePool.put(DeliveryBoy(self.env, i + 1, self, x, y, speed))

    def deliverOrder(self, order):
        # Get DB Id to be assigned or assign First Free
        boyId = self.getDeliveryBoyId(order)
        if boyId is None:
            boy = yield self.freePool.get()
        else:
            boy = yield self.freePool.get(lambda boy: boy.id == boyId)

        # Note down idle time and assign Order
        self.idleTimeTotal += (self.env.now - boy.lastDeliveryTime)
        self.env.process(boy.deliver(order))

    def reportOrderServed(self, boy, deliveryTime):
        self.orderServed += 1
        self.deliveryTimeTotal += deliveryTime
        self.freePool.put(boy)

    def printSummary(self):
        logging.critical("-------- Simulation Summary ------------")
        logging.critical("Orders Served: %d" % self.orderServed)
        logging.critical("Average Delivery Time: %f" % (self.deliveryTimeTotal / self.orderServed))
        logging.critical("Average Idle Time: %f" % (self.idleTimeTotal / self.hireCount))

    def getDeliveryBoyId(self, order):

        if self.algoConfig["type"] == "LEAST_COST":
            # Find DB with least cost free pool and return
            minCost, minId = sys.maxint, None
            for boy in self.freePool.items:
                cost = self.getCost(order, boy)
                if cost < minCost:
                    minCost = cost
                    minId = boy.id
            return minId
        else:
            return None

    def getCost(self, order, boy):
        weight = self.algoConfig["weight"]
        return boy.getCost(order, weight["restaurant"], weight["idle"])
